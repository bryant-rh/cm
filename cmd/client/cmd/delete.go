package cmd

import (
	"github.com/bryant-rh/cm/cmd/client/cmd/delete"
	"github.com/bryant-rh/cm/cmd/client/global"
	"os"

	"github.com/spf13/cobra"
)

func NewCmdDelete() *cobra.Command {
	// certCmd represents the cert command
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "删除一个或更多 resources",
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("cert called")
			_ = cmd.Help()
		},
	}
	deleteCmd.PersistentFlags().BoolVar(&global.EnableDebug, "debug", os.Getenv("DEBUG") == "true", "Enable debug mode")

	deleteCmd.AddCommand(delete.NewCmdDelProject())
	deleteCmd.AddCommand(delete.NewCmdDelCluster())
	deleteCmd.AddCommand(delete.NewCmdDelSa())
	deleteCmd.AddCommand(delete.NewCmdDelLabel())
	deleteCmd.AddCommand(delete.NewCmdDelNs())
	return deleteCmd
}
