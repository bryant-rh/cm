package hx

import (
	"os"

	"github.com/bryant-rh/cm/pkg/hxctx"
	"github.com/spf13/cobra"
)

func NewCmdHxConfig(ctx *hxctx.Context) *cobra.Command {
	var cmdConfig = &cobra.Command{
		Use:   "config",
		Short: "show helmx config",
		Run: func(cmd *cobra.Command, args []string) {
			ctx.Range(func(w *hxctx.Workspace) {
				w.WriteTo(os.Stdout)
			})
		},
	}
	return cmdConfig
}
