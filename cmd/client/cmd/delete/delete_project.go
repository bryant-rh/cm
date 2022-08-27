package delete

import (
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/util"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdDelProject() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "project",
		Short:             "删除项目",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := run_project(global.CMClient)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.MarkFlagRequired("project")
	return cmd
}

func run_project(client *client.CMClient) error {
	klog.V(4).Infoln("Delete Project")

	if !util.AskForConfirmation("delete") {
		os.Exit(1)
	}
	res, err := client.Project_Delete(global.ProjectName)
	if err != nil {
		klog.Fatal(err)
	}

	klog.V(4).Infoln("Delete Project, 输出结果")

	fmt.Println(util.GreenColor(res.Msg))

	return nil
}
