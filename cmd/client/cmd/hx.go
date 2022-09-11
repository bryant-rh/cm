package cmd

import (
	"os"

	"github.com/bryant-rh/cm/cmd/client/cmd/hx"
	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/bryant-rh/cm/pkg/hxctx"
	"github.com/spf13/cobra"
)

var ctx = &hxctx.Context{}

var aliases = map[string]string{
	"CI_SERVER_HOST": "GITLAB_HOST",
	"CI_JOB_TOKEN":   "GITLAB_CI_TOKEN",
}

func NewCmdHx() *cobra.Command {
	hxCmd := &cobra.Command{
		Use: "hx",
		//Version: version.Version,
		Short: "用于初始化项目，部署服务",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if global.ImagePullSecret == "" {
				global.ImagePullSecret = "sensorsdata-registry://harbor.example.com/"
			}
			ctx.Init()

		},
	}

	for k, alias := range aliases {
		if v, ok := os.LookupEnv(k); ok {
			_ = os.Setenv(alias, v)
		}
	}

	if ref := os.Getenv("CI_COMMIT_SHA"); ref != "" && len(ref) > 7 {
		_ = os.Setenv("COMMIT_SHA", ref[0:7])
	}

	flags := hxCmd.PersistentFlags()

	flags.StringVarP(&ctx.Workspace, "workspace", "W", "", "workspace")
	flags.StringVarP(&ctx.DeployEnv, "env", "", os.Getenv("CI_ENVIRONMENT_NAME"), "deploy env")
	flags.StringVarP(&ctx.Feature, "feature", "", os.Getenv("PROJECT_FEATURE"), "project feature")

	workingDir := os.Getenv("WORKING_DIR")

	if workingDir != "" {
		if err := os.Chdir(workingDir); err != nil {
			panic(err)
		}
	}

	hxCmd.AddCommand(hx.NewCmdHxBuildXCISetup())
	hxCmd.AddCommand(hx.NewCmdHxBuildX(ctx))
	hxCmd.AddCommand(hx.NewCmdHxInit(ctx))
	hxCmd.AddCommand(hx.NewCmdHxConfig(ctx))
	hxCmd.AddCommand(hx.NewCmdHxApply(ctx))

	return hxCmd
}
