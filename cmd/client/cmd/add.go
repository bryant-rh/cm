package cmd

import (
	"github.com/bryant-rh/cm/cmd/client/cmd/add"
	"github.com/bryant-rh/cm/cmd/client/global"
	"os"

	"github.com/spf13/cobra"
)

func NewCmdAdd() *cobra.Command {
	// certCmd represents the cert command
	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "添加一个或更多 resources",
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("cert called")
			_ = cmd.Help()
		},
	}

	addCmd.PersistentFlags().BoolVar(&global.EnableDebug, "debug", os.Getenv("DEBUG") == "true", "Enable debug mode")
	addCmd.PersistentFlags().BoolVar(&global.NoFormat, "no-format", global.NoFormat, "If present, print output without format table")

	addCmd.AddCommand(add.NewCmdAddNs())
	addCmd.AddCommand(add.NewCmdAddLabel())

	return addCmd
}
