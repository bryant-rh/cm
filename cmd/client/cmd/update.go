package cmd

import (
	"github.com/bryant-rh/cm/cmd/client/cmd/update"
	"github.com/bryant-rh/cm/cmd/client/global"
	"os"

	"github.com/spf13/cobra"
)

func NewCmdUpdate() *cobra.Command {
	// certCmd represents the cert command
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "更新一个或更多 resources",
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("cert called")
			_ = cmd.Help()
		},
	}
	updateCmd.PersistentFlags().BoolVar(&global.EnableDebug, "debug", os.Getenv("DEBUG") == "true", "Enable debug mode")

	updateCmd.AddCommand(update.NewCmdUpdateProject())
	updateCmd.AddCommand(update.NewCmdUpdateClster())
	updateCmd.AddCommand(update.NewCmdUpdateSa())
	return updateCmd
}
