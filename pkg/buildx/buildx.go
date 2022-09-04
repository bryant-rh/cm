package buildx

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bryant-rh/cm/pkg/hxctx"
	"github.com/bryant-rh/cm/pkg/util"
	"github.com/go-courier/dockerfileyml"
	"github.com/go-courier/helmx/spec"
	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
)

type BuildxFlags struct {
	NoOmit      bool
	WithBuilder bool
	Push        bool
	Load        bool
	LocalCache  string
	Platform    string
}

func RunDockerBuildxSetup(buildkitName string) {
	if v := os.Getenv("MULTI_ARCH_BUILDER"); v == "1" || v == "true" {
		util.StdRun("docker", "buildx", "create", "--use", "--name", buildkitName, "--node", buildkitName+"-amd64", "--driver=kubernetes", "--driver-opt=namespace=gitlab")
	} else {
		util.StdRun("docker", "buildx", "create", "--use", "--name", buildkitName, "--platform=linux/amd64", "--node", buildkitName+"-amd64", "--driver=kubernetes", "--driver-opt=namespace=gitlab")
		util.StdRun("docker", "buildx", "create", "--append", "--name", buildkitName, "--platform=linux/arm64", "--node", buildkitName+"-arm64", "--driver=kubernetes", "--driver-opt=namespace=gitlab")
	}
	util.StdRun("docker", "buildx", "inspect", buildkitName)
}

func RunDockerBuildx(workspace *hxctx.Workspace, flags BuildxFlags) {
	dockerfilename := workspace.Path("./Dockerfile")

	projectEnvVars := map[string]string{}

	projectEnvVars["PROJECT_VERSION"] = workspace.Project.Version.String()
	projectEnvVars["PROJECT_NAME"] = workspace.Project.Name
	projectEnvVars["PROJECT_GROUP"] = workspace.Project.Group
	projectEnvVars["PROJECT_FEATURE"] = workspace.Project.Feature

	dockerfile := dockerfileyml.Dockerfile{}

	hasDockerfileYaml := loadDockerfileYAMLs(&dockerfile, workspace.EnvVarSet.Defaults,
		workspace.Path("dockerfile.default.yml"),
		workspace.Path("dockerfile.yml"),
	)

	if dockerfile.Image == "" {
		dockerfile.Image = workspace.Project.Image
	}

	if hasDockerfileYaml {
		dockerfilename = workspace.Path("./out/Dockerfile")

		df := bytes.NewBuffer(nil)

		if strings.HasPrefix(dockerfile.From, "env-") {
			dockerfile.From = (&spec.Image{}).ResolveImagePullSecret().PrefixTag("~infrav2/env-" + dockerfile.From[4:])

			if flags.WithBuilder {
				dockerfile.Stages = map[string]*dockerfileyml.Stage{
					"builder": {
						From:       dockerfile.From + ":onbuild",
						WorkingDir: "/go/src",
						Add: dockerfileyml.Values{
							"./": "./",
						},
						Run: dockerfileyml.Scripts(
							fmt.Sprintf("--mount=type=cache,sharing=locked,id=gomod,target=/go/pkg/mod WORKSPACE=%s make build", workspace.Name),
						),
					},
				}

				copies := map[string]string{}

				for k, v := range dockerfile.Add {
					r, _ := filepath.Rel(workspace.ProjectRoot, workspace.WorkspaceRoot)
					copies["builder:"+path.Join(r, k)] = v
				}

				dockerfile.Copy = copies
				dockerfile.Add = nil
			}

			dockerfile.From = dockerfile.From + ":runtime"
		}

		injectProjectValues := func(state *dockerfileyml.Stage) {
			if state.Env == nil {
				state.Env = dockerfileyml.Values{}
			}

			for k, v := range projectEnvVars {
				state.Env[k] = v
			}
		}

		injectProjectValues(&dockerfile.Stage)

		for _, s := range dockerfile.Stages {
			injectProjectValues(s)
		}

		if err := dockerfileyml.WriteToDockerfile(df, dockerfile); err != nil {
			panic(err)
		}

		_ = os.MkdirAll(path.Dir(dockerfilename), os.ModePerm)
		if err := ioutil.WriteFile(dockerfilename, df.Bytes(), os.ModePerm); err != nil {
			panic(err)
		}
	}

	if flags.NoOmit {
		return
	}

	dockerfileData, _ := ioutil.ReadFile(dockerfilename)

	args := []string{"docker", "buildx", "build", "--progress", "plain"}

	for _, buildArg := range parseBuildArgs(dockerfileData) {
		buildArgValue := ""

		if projectEnvVar := projectEnvVars[buildArg]; projectEnvVar != "" {
			buildArgValue = projectEnvVar
		}

		if workspace.EnvVarSet.Defaults != nil {
			if defaultValue := workspace.EnvVarSet.Defaults[buildArg]; defaultValue != "" {
				buildArgValue = defaultValue
			}
		}

		if envValue := os.Getenv(buildArg); envValue != "" {
			buildArgValue = envValue
		}

		if buildArgValue != "" {
			args = append(args, "--build-arg", fmt.Sprintf("%s=%s", buildArg, buildArgValue))
		}
	}

	if flags.Push {
		args = append(args, "--push")
	}

	if flags.Load {
		args = append(args, "--load")
	}

	if flags.Platform != "" {
		args = append(args, "--platform", flags.Platform)
	}

	if flags.LocalCache != "" {
		args = append(args, "--cache-from", "type=local,src="+flags.LocalCache)
		args = append(args, "--cache-to", "type=local,dest="+flags.LocalCache)
	}

	if v := os.Getenv("CI_PROJECT_URL"); v != "" {
		args = append(args, "--label", "org.opencontainers.image.source="+v)
	}

	if v := os.Getenv("CI_COMMIT_SHA"); v != "" {
		args = append(args, "--label", "org.opencontainers.image.revision="+v)
	}

	if v := os.Getenv("CI_COMMIT_AUTHOR"); v != "" {
		commitUserEmail := getEmail(v)

		if gitlabUserEmail := os.Getenv("GITLAB_USER_EMAIL"); gitlabUserEmail != "" && gitlabUserEmail != commitUserEmail {
			commitUserEmail = gitlabUserEmail
		}

		if strings.HasSuffix(commitUserEmail, "@sensorsdata.cn") {
			args = append(args, "--label", "org.opencontainers.image.maintainer="+commitUserEmail)
		} else {
			panic(fmt.Errorf("unsupported commitUserEmail domain %s, must be `@sensorsdata.cn` ", commitUserEmail))
		}
	}

	args = append(args, "--label", "org.opencontainers.image.version="+workspace.Project.Version.String())
	args = append(args, "--label", "org.opencontainers.image.title="+workspace.Project.Name)
	args = append(args, "--label", "org.opencontainers.image.created="+time.Now().Format(time.RFC3339))

	if data, err := ioutil.ReadFile(dockerfilename); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(aurora.Yellow("nothing to buildx"))
			return
		}
	} else {
		util.PrintifyDockerfile(os.Stdout, data)
	}

	util.StdRun(append(args, "--file", dockerfilename, "--tag", dockerfile.Image, ".")...)
}

func loadDockerfileYAMLs(d *dockerfileyml.Dockerfile, envs map[string]string, dockerfileYAMLs ...string) bool {
	hasDockerfileYaml := false

	for _, dockerfileYAML := range dockerfileYAMLs {
		specFileContent, err := ioutil.ReadFile(dockerfileYAML)
		if err == nil {
			hasDockerfileYaml = true

			if err := yaml.Unmarshal(hxctx.ResolveEnvVars(envs, specFileContent), d); err != nil {
				panic(err)
			}
		}
	}

	return hasDockerfileYaml
}

var reBuildArg = regexp.MustCompile("ARG ([^=\n]+)")

func parseBuildArgs(data []byte) []string {
	args := map[string]bool{}

	matched := reBuildArg.FindAllSubmatch(data, -1)

	for i := range matched {
		arg := string(matched[i][1])

		if isPlatformArg(arg) {
			continue
		}

		args[arg] = true
	}

	finals := make([]string, 0)

	for arg := range args {
		finals = append(finals, arg)
	}

	return finals
}

var platformArgs = []string{
	"TARGETPLATFORM",
	"TARGETOS",
	"TARGETARCH",
	"TARGETVARIANT",

	"BUILDPLATFORM",
	"BUILDOS",
	"BUILDARCH",
	"BUILDVARIANT",
}

func isPlatformArg(arg string) bool {
	for _, platformArg := range platformArgs {
		if platformArg == arg {
			return true
		}
	}
	return false
}

var reEmail = regexp.MustCompile(`.+<(.+@[^>]+)>`)

func getEmail(v string) string {
	data := reEmail.FindStringSubmatch(v)
	if len(data) == 0 {
		return ""
	}
	return data[1]
}
