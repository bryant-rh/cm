package cmd

import (
	"os"

	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/bryant-rh/cm/pkg/helmx"

	hxc "github.com/bryant-rh/cm/cmd/client/cmd/hx"

	"github.com/go-courier/helmx/spec"
	"github.com/spf13/cobra"
)

func NewCmdHx() *cobra.Command {
	hxCmd := &cobra.Command{
		Use: "hx",
		Short: "用于初始化项目，部署服务",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if global.ImagePullSecret == "" {
				global.ImagePullSecret = "sensorsdata-registry://harbor.example.com/rk-"
			}

			global.Ctx.DefaultEnvs = map[string]string{}
			global.Ctx.Envs = map[string]string{
				spec.EnvKeyImagePullSecret: global.ImagePullSecret,
			}
			commitRefName := os.Getenv(helmx.EnvVarKeyCommitRefName)

			if commitRefName != "" {
				feat, env := helmx.GetFeatAndEnv(commitRefName)
				if feat != "" {
					global.Ctx.Project.Feature = feat
				}
				if env != "" {
					global.Ctx.Environment = env
				}
			}

			{
				helmx.LoadEnvConfigFromFiles(global.Ctx.DefaultEnvs, "default", "")

				for k, v := range global.Ctx.DefaultEnvs {
					global.Ctx.Envs[k] = v
				}

				helmx.LoadEnvConfig(global.Ctx.Envs, global.Ctx.Environment, global.Ctx.Project.Feature)
			}

			{
				specFiles := []string{
					global.HelmxRootFile,
					"./helmx.default.yml",
					"./helmx.yml",
				}

				if global.Ctx.Project.Feature != "" {
					specFiles = append(specFiles, "./helmx."+global.Ctx.Project.Feature+".yml")
				}

				if global.Ctx.Environment != "" {
					specFiles = append(specFiles, "./helmx."+global.Ctx.Environment+".yml")
					if global.Ctx.Project.Feature != "" {
						specFiles = append(specFiles, "./helmx."+global.Ctx.Environment+"-"+global.Ctx.Project.Feature+".yml")
					}
				}

				for _, file := range specFiles {
					helmx.LoadHelmxSpecFile(&global.Ctx.Spec, file)
				}

				// if !livenessCheckSkip {
				// 	if ctx.Spec.Service.Pod.LivenessProbe == nil || ctx.Spec.Service.Pod.ReadinessProbe == nil {
				// 		err := errors.New(`
				// 		################################################################
				// 		livenessProbe 或 readinessProbe 必须存在。
				// 		更多 CICD 流程，
				// 			参考文档:
				// 		################################################################
				// 		`)
				// 		log.Fatal(err)
				// 	}
				// }

			}

			{
				referFile := []string{
					"./go.mod",
				}
				for _, file := range referFile {
					refer := helmx.LoadModFile(file)
					if refer != nil {
						global.Ctx.ReferFiles = refer
					}
				}
			}

			{
				pomFile := []string{
					"pom.txt",
				}
				for _, file := range pomFile {
					refer := helmx.LoadPomFile(file)
					if refer != nil {
						global.Ctx.ReferFiles = append(global.Ctx.ReferFiles, refer...)
					}
				}
			}
			{
				suffix := helmx.GetShortSha()
				if suffix != "" {
					global.Ctx.Project.Version.Suffix = suffix
				}

				// get commit message
				global.Ctx.CommitShaMessage = os.Getenv(helmx.EnvVarKeyCommitMessage)
				global.Ctx.CommitSha = os.Getenv(helmx.EnvVarKeyCommitSha)

				// overwrite by --version
				if global.Ctx.Version != "" {
					v, err := spec.ParseVersion(global.Ctx.Version)
					if err != nil {
						panic(err)
					}
					global.Ctx.Project.Version = *v
				}

				helmx.MustMarshalToEnvs(global.Ctx.Envs, global.Ctx.Project, "PROJECT_")
			}
			helmx.SetEnviron(global.Ctx.Envs)
		},
	}

	global.Ctx.Project = &spec.Project{}

	flags := hxCmd.PersistentFlags()
	flags.BoolVarP(&global.LivenessCheckSkip, "liveness-check-skip", "", false, "skip liveness/readiness check")

	flags.StringVarP(&global.HelmxRootFile, "hxroot", "", "./helmx.project.yml", "project config file")
	flags.StringVarP(&global.Ctx.Project.Feature, "feature", "", os.Getenv("PROJECT_FEATURE"), "project feature")
	flags.StringVarP(&global.Ctx.Project.Description, "description", "", os.Getenv("PROJECT_DESCRIPTION"), "project description")
	flags.StringVarP(&global.Ctx.Version, "version", "", "", "project version")
	flags.StringVarP(&global.Ctx.Environment, "env", "", os.Getenv("CI_ENVIRONMENT_NAME"), "deploy env")

	workingDir := os.Getenv("WORKING_DIR")

	if workingDir != "" {
		if err := os.Chdir(workingDir); err != nil {
			panic(err)
		}
	}

	hxCmd.AddCommand(hxc.NewCmdHxInit())
	hxCmd.AddCommand(hxc.NewCmdHxApply())
	hxCmd.AddCommand(hxc.NewCmdHxCreateProxy())
	hxCmd.AddCommand(hxc.NewCmdHxCreateKubeConfig())
	return hxCmd
}
