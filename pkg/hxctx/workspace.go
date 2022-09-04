package hxctx

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/bryant-rh/cm/pkg/temp"
	"github.com/bryant-rh/cm/pkg/util"
	"github.com/go-courier/helmx/spec"
	"github.com/go-courier/semver"

	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewWorkspace(workspaceRoot string, projectRoot string, setDefault func(w *Workspace)) *Workspace {
	w := Workspace{
		ProjectRoot:   projectRoot,
		WorkspaceRoot: workspaceRoot,
		IsMain:        filepath.Base(workspaceRoot) == filepath.Base(projectRoot),
	}

	w.loadSpecs()
	setDefault(&w)
	w.setProjectDefaults()
	w.loadEnvs()

	// put project values as env var for deployment
	MustMarshalToEnvs(w.EnvVarSet.Values, w.Project, "PROJECT_")

	return &w
}

type Refer struct {
	ReferName    string `json:"referName"`
	ReferVersion string `json:"referVersion"`
}

type Refers []Refer

type Workspace struct {
	IsMain          bool   `json:"-"`
	ProjectRoot     string `json:"-"`
	WorkspaceRoot   string `json:"-"`
	Project         `json:"project"`
	EnvVarSet       EnvVarSet        `json:"envVars,omitempty"`
	DeployEnv       string           `json:"deployEnv"`
	DeployResources []DeployResource `json:"deployResources,omitempty"`

	spec.Spec
	DefaultEnvs      map[string]string
	Environment      string
	Version          string
	Templates        map[string][]byte
	ReferFiles       Refers
	CommitShaMessage string
	CommitSha        string
}

type EnvVarSet struct {
	Values    map[string]string `json:"values,omitempty"`
	Defaults  map[string]string `json:"defaults,omitempty"`
	FromFiles []string          `json:"fromFiles,omitempty"`
}

type DeployResource struct {
	runtime.Object `json:"object"`
	FromFiles      []string `json:"fromFiles"`
}

func (w *Workspace) WriteTo(writer io.Writer) (int64, error) {
	err := util.PrintifyJSON(writer, w)
	if err != nil {
		return -1, err
	}
	return 0, nil
}

func (w *Workspace) Path(f string) string {
	return filepath.Join(w.WorkspaceRoot, f)
}

func (w *Workspace) loadEnvs() {
	w.EnvVarSet = EnvVarSet{
		Values:   map[string]string{},
		Defaults: map[string]string{},
	}

	feature := w.Project.Feature
	env := w.DeployEnv

	configFiles := []string{
		w.Path("./config/default.yml"),
		w.Path("./config/master.yml"),
	}

	if feature != "" {
		configFiles = append(configFiles, w.Path("./config/"+feature+".yml"))
	}

	// default demo
	if env == "" {
		env = "demo"
	}

	if env != "" {
		configFiles = append(configFiles, w.Path("./config/"+env+".yml"))
		if feature != "" {
			configFiles = append(configFiles, w.Path("./config/"+env+"-"+feature+".yml"))
		}
	}

	for _, filename := range configFiles {
		specFileContent, err := ioutil.ReadFile(filename)
		if err == nil {
			w.EnvVarSet.FromFiles = append(w.EnvVarSet.FromFiles, filename)

			if strings.HasSuffix(filename, "default.yml") {
				if err := yaml.Unmarshal(specFileContent, &w.EnvVarSet.Defaults); err != nil {
					panic(err)
				}
			}

			if err := yaml.Unmarshal(specFileContent, &w.EnvVarSet.Values); err != nil {
				panic(err)
			}
		}
	}
}

// todo remove in future
func (w *Workspace) loadSpecs() {

	if w.EnvVarSet.Values == nil {
		w.EnvVarSet.Values = map[string]string{
			spec.EnvKeyImagePullSecret: global.ImagePullSecret,
		}
	}
	SetEnviron(w.EnvVarSet.Values)

	// load config/default.yml for overwrites of dynamic project
	defaultConfig, err := ioutil.ReadFile(w.Path("./config/default.yml"))
	if err == nil {
		if err := yaml.Unmarshal(defaultConfig, &w.EnvVarSet.Values); err != nil {
			panic(err)
		}
	}

	feature := w.Project.Feature
	env := w.DeployEnv

	specFiles := []string{
		w.Path("./helmx.project.yml"),
		w.Path("./helmx.default.yml"),
		w.Path("./helmx.yml"),
	}

	if feature != "" {
		specFiles = append(specFiles, w.Path("./helmx."+feature+".yml"))
	}

	if env != "" {
		specFiles = append(specFiles, w.Path("./helmx."+env+".yml"))
		if feature != "" {
			specFiles = append(specFiles, w.Path("./helmx."+env+"-"+feature+".yml"))
		}
	}

	for _, filename := range specFiles {
		specFileContent, err := ioutil.ReadFile(filename)
		if err == nil {
			if err := yaml.Unmarshal(ResolveEnvVars(w.EnvVarSet.Values, specFileContent), &w.Spec); err != nil {
				panic(err)
			}
		}
	}

	if w.Spec.Project != nil {
		w.Project.Group = w.Spec.Project.Group
		w.Project.Name = w.Spec.Project.Name
		w.Project.Feature = w.Spec.Project.Feature
		w.Project.Description = w.Spec.Project.Description

		if v, err := semver.ParseVersion(w.Spec.Project.Version.String()); err == nil {
			w.Project.Version = *v
		}
	}

	// if s.Service != nil {
	// 	q := v1alpha1.QService{}

	// 	// helmx 中的结构与 QService 中的结构不一样
	// 	buf := bytes.NewBuffer(nil)
	// 	_ = json.NewEncoder(buf).Encode(s.Service)
	// 	_ = json.NewDecoder(buf).Decode(&q.Spec)

	// 	// 将其他字段(volume, toleration)赋值给 qservice
	// 	buf = bytes.NewBuffer(nil)
	// 	_ = json.NewEncoder(buf).Encode(s)
	// 	_ = json.NewDecoder(buf).Decode(&q.Spec)

	// 	q.Spec.Image = "${{ PROJECT_IMAGE }}"

	// 	if conf, err := ioutil.ReadFile(w.Path("config/default.yml")); err == nil {
	// 		m := map[string]string{}
	// 		_ = yaml.Unmarshal(conf, &m)

	// 		if q.Spec.Envs == nil {
	// 			q.Spec.Envs = map[string]string{}
	// 		}

	// 		for k := range m {
	// 			if strings.HasPrefix(k, "PROJECT_") {
	// 				continue
	// 			}

	// 			q.Spec.Envs[k] = fmt.Sprintf(`${{ %s }}`, k)
	// 		}
	// 	}

	// 	q.SetGroupVersionKind(v1alpha1.SchemeGroupVersion.WithKind("QService"))

	// 	data, _ := yaml2.Marshal(q)
	// 	_ = os.MkdirAll(w.Path("deploy"), os.ModePerm)
	// 	_ = ioutil.WriteFile(w.Path("deploy/qservice.yml"), data, os.ModePerm)
	// }
}

func (w *Workspace) setProjectDefaults() {
	if ver, err := ioutil.ReadFile(filepath.Join(w.ProjectRoot, ".version")); err == nil {
		v, err := semver.ParseVersion(string(bytes.TrimSpace(ver)))
		if err != nil {
			panic(err)
		}
		w.Project.Version = *v
	}

	if sha := os.Getenv("COMMIT_SHA"); sha != "" {
		v, err := w.Project.Version.WithPrerelease(sha)
		if err == nil {
			w.Project.Version = *v
		}
	}

	if w.Project.Group == "" {
		projectPath := os.Getenv("CI_PROJECT_PATH")

		if projectPath == "" {
			// if no CI_PROJECT_PATH, use project root
			parts := strings.Split(path.Dir(w.ProjectRoot), "/")

			// git.xxx/a/b
			for i, part := range parts {
				if strings.HasPrefix(part, "git") && i < len(parts)-1 {
					projectPath = strings.Join(parts[i+1:], "/")
				}
			}
		}

		w.Project.Group = strings.Split(projectPath, "/")[0]
	}

	if w.Project.Name == "" {
		if w.IsMain {
			w.Project.Name = filepath.Base(w.WorkspaceRoot)
		} else {
			w.Project.Name = fmt.Sprintf("%s-%s", filepath.Base(w.ProjectRoot), filepath.Base(w.WorkspaceRoot))
		}
	}

	w.Project.Image = (&spec.Image{}).ResolveImagePullSecret().PrefixTag(w.Project.DefaultImageTag())
}

// func (w *Workspace) loadDeployments() {
// 	deployRootDir := w.Path("deploy")

// 	files, err := ioutil.ReadDir(deployRootDir)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			return
// 		}
// 		panic(err)
// 	}

// 	kinds := map[string]bool{}

// 	for i := range files {
// 		f := files[i]

// 		if f.IsDir() {
// 			continue
// 		}

// 		kinds[strings.Split(f.Name(), ".")[0]] = true
// 	}

// 	kindList := make([]string, 0)

// 	for k := range kinds {
// 		kindList = append(kindList, k)
// 	}

// 	sort.Strings(kindList)

// 	for i := range kindList {
// 		kind := kindList[i]

// 		o, err := readDeployResource(w, kind)
// 		if err != nil {
// 			panic(err)
// 		}

// 		if o == nil {
// 			continue
// 		}

// 		w.DeployResources = append(w.DeployResources, *o)
// 	}
// }

// func readDeployResource(w *Workspace, kind string) (*DeployResource, error) {
// 	feature := w.Project.Feature
// 	env := w.DeployEnv

// 	files := []string{
// 		w.Path(fmt.Sprintf("deploy/%s.yml", kind)),
// 		w.Path(fmt.Sprintf("deploy/%s.overwrites.yml", kind)),
// 	}

// 	if feature != "" {
// 		files = append(files, w.Path(fmt.Sprintf("deploy/%s.overwrites.%s.yml", kind, feature)))
// 	}

// 	if env != "" {
// 		files = append(files, w.Path(fmt.Sprintf("deploy/%s.overwrites.%s.yml", kind, env)))
// 		if feature != "" {
// 			files = append(files, w.Path(fmt.Sprintf("deploy/%s.overwrites.%s-%s.yml", kind, env, feature)))
// 		}
// 	}

// 	return loadDeployResource(w, files...)
// }

// func loadDeployResource(w *Workspace, files ...string) (*DeployResource, error) {
// 	dr := DeployResource{}

// 	for _, filename := range files {
// 		data, err := ioutil.ReadFile(filename)

// 		if err == nil {
// 			dr.FromFiles = append(dr.FromFiles, filename)

// 			data = MustNewTemplate(data).Execute(w.EnvVarSet.Values)

// 			if dr.Object == nil {
// 				o, err := kubeutil.UnmarshalToObject(data)
// 				if err != nil {
// 					return nil, err
// 				}
// 				dr.Object = o
// 			} else {
// 				if err := yaml2.Unmarshal(data, dr.Object); err != nil {
// 					return nil, err
// 				}
// 			}
// 		}
// 	}

// 	if dr.Object == nil {
// 		return nil, nil
// 	}

// 	return &dr, nil
// }

type Project struct {
	Name        string         `env:"name" json:"name"`
	Feature     string         `env:"feature"  json:"feature,omitempty"`
	Version     semver.Version `env:"version" json:"version"`
	Group       string         `env:"group"  json:"group,omitempty"`
	Description string         `env:"description" yaml:"description,omitempty" json:"description,omitempty"`
	Image       string         `env:"image" json:"image"`
}

func (p Project) FullName() string {
	if p.Feature != "" {
		return p.Name + "--" + p.Feature
	}
	return p.Name
}

func (p Project) DefaultImageTag() string {
	return "~" + p.Group + "/" + p.Name + ":" + strings.Replace(p.Version.String(), "+", "--", -1)
}

func (w *Workspace) InitProject() {

	answers := struct {
		Name        string
		Group       string
		Description string
		Version     string
	}{}

	questings := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "项目名称",
				Default: w.Project.Name,
			},
			Validate: survey.Required,
		},
		{
			Name: "description",
			Prompt: &survey.Input{
				Message: "项目描述",
				Default: w.Project.Description,
			},
			Validate: survey.Required,
		},
		{
			Name: "group",
			Prompt: &survey.Input{
				Message: "项目所属部门",
				Default: w.Project.Group,
			},
			Validate: survey.Required,
		},
		{
			Name: "version",
			Prompt: &survey.Input{
				Message: "项目版本号 (x.x.x)",
				Default: w.Project.Version.String(),
			},
			Validate: survey.Required,
		},
	}

	err := survey.Ask(questings, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	w.Spec.Project.Name = answers.Name
	w.Spec.Project.Group = answers.Group
	w.Spec.Project.Description = answers.Description

	if answers.Version != "" {
		if v, err := semver.ParseVersion(answers.Version); err == nil {
			w.Project.Version = *v
		}

	}

	MustWriteYAML(w.Path("./helmx.project.yml"), struct {
		Project *spec.Project
	}{
		Project: w.Spec.Project,
	})
	//helmx.MustWriteYAML(w.Path("./helmx.project.yml"), w.Project)

	//生成dockerfile.default
	file_type := ""
	switch {
	case util.Exists("./go.mod"):
		file_type = "go"
	case util.Exists("./pom.xml"):
		file_type = "java"
	case util.Exists("./requirements.txt"):
		file_type = "python"
	}
	fmt.Println(util.GreenColor("生成helmx.yml"))
	_ = util.WriteToFile(w.Path("./helmx.yml"), temp.Helmxfile(), func(v interface{}) ([]byte, error) {
		return v.([]byte), nil
	})

	fmt.Println(util.GreenColor("生成dockerfile.default"))
	_ = util.WriteToFile(w.Path("./Dockerfile.default"), temp.Dockerfile(file_type, filepath.Base(w.WorkspaceRoot)), func(v interface{}) ([]byte, error) {
		return v.([]byte), nil
	})

	if file_type == "go" {
		fmt.Println(util.GreenColor("生成Makefile.default"))
		_ = util.WriteToFile("./Makefile.default", temp.Makefile(), func(v interface{}) ([]byte, error) {
			return v.([]byte), nil
		})

	}

	fmt.Println(util.GreenColor("项目初始化完成"))

}
