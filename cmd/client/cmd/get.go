package cmd

import (
	"os"

	"github.com/bryant-rh/cm/cmd/client/cmd/get"
	"github.com/bryant-rh/cm/cmd/client/global"

	"github.com/spf13/cobra"
)

func NewCmdGet() *cobra.Command {
	// certCmd represents the cert command
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "显示一个或更多 resources",
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("cert called")
			_ = cmd.Help()
		},
	}

	getCmd.PersistentFlags().BoolVar(&global.EnableDebug, "debug", os.Getenv("DEBUG") == "true", "Enable debug mode")
	getCmd.PersistentFlags().BoolVar(&global.NoFormat, "no-format", global.NoFormat, "If present, print output without format table")

	getCmd.AddCommand(get.NewCmdGetProject())
	getCmd.AddCommand(get.NewCmdGetCluster())
	getCmd.AddCommand(get.NewCmdGetSa())
	getCmd.AddCommand(get.NewCmdGetLabel())
	getCmd.AddCommand(get.NewCmdGetToken())

	return getCmd
}
