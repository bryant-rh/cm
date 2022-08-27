package helmx

import (
	"github.com/bryant-rh/cm/cmd/client/global"
	"fmt"

	//"git.querycap.com/tools/confclient/v2"
	//"hx/client_spec"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/go-courier/helmx/spec"
	"gopkg.in/yaml.v2"
)

func LoadHelmxSpecFile(s *spec.Spec, filename string) {
	fmt.Fprintf(os.Stdout, "try to load helmx from %s ...\n", filename)

	specFileContent, err := ioutil.ReadFile(filename)
	if err == nil {
		specFileContent = []byte(ResolveEnvVars(s.Envs, string(specFileContent)))
		err := yaml.Unmarshal(specFileContent, s)
		if err != nil {
			panic(err)
		}
	}
}

func LoadModFile(filename string) global.Refers {
	fmt.Fprintf(os.Stdout, "try to load mod from %s ...\n", filename)
	content, err := ioutil.ReadFile(filename)
	reg := regexp.MustCompile(`[a-zA-Z0-9]+/[a-z]+`)
	var refers global.Refers
	if err == nil {
		for _, line := range strings.Split(string(content), "\n") {
			if !strings.Contains(line, "module") {
				str := strings.Replace(line, "\t", "", -1)
				has := reg.MatchString(str)
				if has {
					lineArray := strings.Split(str, " ")
					if len(lineArray) >= 2 {
						refers = append(refers, global.Refer{
							ReferName:    lineArray[0],
							ReferVersion: lineArray[1],
						})

					}
				}
			}

		}
	}
	return refers
}

const (
	GroupId = iota
	ArtifactId
	PackageType
	Version
	LifeCycle
)

var gitCachePath = "/go/pkg/mod/gitlab-cache/.m2/repository"

func LoadPomFile(filename string) global.Refers {
	group := os.Getenv(EnvVarKeyProjectNamespace)
	projectName := os.Getenv(EnvVarKeyProjectName)

	fileName := fmt.Sprintf("%s/%s-%s-%s", gitCachePath, group, projectName, filename)

	fmt.Fprintf(os.Stdout, "try to load mod from %s ...\n", fileName)

	content, err := ioutil.ReadFile(fileName)
	reg := regexp.MustCompile(`[a-zA-Z0-9]+`)
	var refers global.Refers
	if err == nil {
		start := false
		for _, line := range strings.Split(string(content), "\n") {
			if start {
				strs := strings.Split(line, "[INFO]")
				refStr := ""
				if len(strs) == 2 {
					if strs[0] == " " && strs[1] == " " {
						break
					}
					has := reg.MatchString(strs[1])
					if has {
						refStr = strings.Replace(strs[1], " ", "", -1)
						if !strings.Contains(refStr, "--") && refStr != "" {
							versions := strings.Split(refStr, ":")
							if len(versions) == 5 {
								refers = append(refers, global.Refer{
									ReferName:    fmt.Sprintf("%s:%s.%s", versions[GroupId], versions[ArtifactId], versions[PackageType]),
									ReferVersion: fmt.Sprintf("%s:%s", versions[Version], versions[LifeCycle]),
								})
							}
						}
					}
				}

			}

			if strings.Contains(line, "The following files have been resolved") {
				start = true
			}
		}
	}
	return refers
}

// // 推送配置到golive
// func PushSpecToGolive(ctx *Context) {

// 	cli := &confclient.Client{
// 		Host: "golive.rockontrol.com",
// 	}
// 	cli.Init()
// 	cli.SetDefaults()

// 	ignoreEnvKeys := []string{
// 		"IMAGE_PULL_SECRET",
// 		"PROJECT_DESCRIPTION",
// 		"PROJECT_GROUP",
// 		"PROJECT_NAME",
// 		"PROJECT_VERSION",
// 	}

// 	for _, key := range ignoreEnvKeys {
// 		if _, ok := ctx.Envs[key]; ok {
// 			delete(ctx.Envs, key)
// 		}
// 	}

// 	req := client_spec.AddProjectSpec{}
// 	req.GroupName = ctx.Project.Group
// 	req.ProjectName = ctx.Project.Name
// 	req.Body.ProjectName = ctx.Project.Name
// 	req.Body.GroupName = ctx.Project.Group
// 	req.Body.Image = ctx.Project.Version.String()
// 	req.Body.Configuration.Spec = ctx.Spec
// 	req.Body.Configuration.Templates = ctx.Templates
// 	req.Body.RefName = os.Getenv(EnvVarKeyCommitRefName)
// 	req.Body.Refers = ctx.ReferFiles
// 	req.Body.Desc = ctx.CommitShaMessage
// 	req.Body.CommitSha = ctx.CommitSha

// 	fmt.Println("CommitSha:", req.Body.CommitSha)
// 	fmt.Println("CommitMessage:", req.Body.Desc)
// 	fmt.Println("Body Refers:")
// 	for _, refer := range req.Body.Refers {
// 		fmt.Println(fmt.Sprintf("ReferName: %s, ReferVersion:%s", refer.ReferName, refer.ReferVersion))
// 	}

// 	upstream := make([]string, 0)
// 	for _, up := range req.Body.Configuration.Spec.Upstreams {
// 		if up != "" {
// 			upstream = append(upstream, up)
// 		}
// 	}
// 	req.Body.Configuration.Spec.Upstreams = upstream
// 	_, err := req.Invoke(cli)
// 	if err != nil {
// 		panic(err)
// 	}
// }
