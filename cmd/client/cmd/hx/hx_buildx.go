package hx

import (
	"github.com/bryant-rh/cm/pkg/buildx"
	"github.com/bryant-rh/cm/pkg/hxctx"
	"github.com/spf13/cobra"
)

var (
	buildXFlags buildx.BuildxFlags
)

func NewCmdHxBuildX(ctx *hxctx.Context) *cobra.Command {
	cmdBuildX := &cobra.Command{
		Use:   "buildx",
		Short: "Ship by dockerfile.yml or Dockerfile",
		Run: func(cmd *cobra.Command, args []string) {
			ctx.Range(func(w *hxctx.Workspace) {
				buildx.RunDockerBuildx(w, buildXFlags)
			})
		},
	}

	cmdBuildX.Flags().BoolVarP(&buildXFlags.NoOmit, "no-omit", "", false, "generate Dockerfile only")
	cmdBuildX.Flags().BoolVarP(&buildXFlags.Load, "load", "", false, "load after build")
	cmdBuildX.Flags().BoolVarP(&buildXFlags.Push, "push", "", false, "push after build")
	cmdBuildX.Flags().BoolVarP(&buildXFlags.WithBuilder, "with-builder", "", false, "added multi stage build for crossing build")
	cmdBuildX.Flags().StringVarP(&buildXFlags.LocalCache, "local-cache", "", "", "buildx cache with local absolute path")
	cmdBuildX.Flags().StringVarP(&buildXFlags.Platform, "platform", "", "linux/amd64,linux/arm64", "platform")

	return cmdBuildX
}
