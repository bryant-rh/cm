package hx

import (
	"github.com/bryant-rh/cm/pkg/hxctx"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func NewCmdHxInit(ctx *hxctx.Context) *cobra.Command {
	// certCmd represents the cert command
	cmd := &cobra.Command{
		Use:   "init",
		Short: "project initial",
		Run: func(cmd *cobra.Command, args []string) {
			klog.V(4).Infoln("Hx Init")
			ctx.Range(func(w *hxctx.Workspace) {
				w.InitProject()
			})
		},
	}

	return cmd
}
