package cmd

import (
	"github.com/bryant-rh/cm/cmd/client/cmd/create"
	"github.com/bryant-rh/cm/cmd/client/global"
	"os"

	"github.com/spf13/cobra"
)

func NewCmdCreate() *cobra.Command {
	// certCmd represents the cert command
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "创建一个或更多 resources",
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("cert called")
			_ = cmd.Help()
		},
	}

	createCmd.PersistentFlags().BoolVar(&global.EnableDebug, "debug", os.Getenv("DEBUG") == "true", "Enable debug mode")
	createCmd.PersistentFlags().BoolVar(&global.NoFormat, "no-format", global.NoFormat, "If present, print output without format table")

	createCmd.AddCommand(create.NewCmdCreateProject())
	createCmd.AddCommand(create.NewCmdCreateCluster())
	createCmd.AddCommand(create.NewCmdCreateSa())

	return createCmd
}
