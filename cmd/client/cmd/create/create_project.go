package create

import (
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/util"
	"fmt"

	"github.com/bryant-rh/cm/pkg/output"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdCreateProject() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "project",
		Short:             "创建项目",
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
	cmd.Flags().StringVarP(&global.Description, "description", "d", global.Description, "项目描述信息")
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("description")
	return cmd
}

func run_project(client *client.CMClient) error {
	klog.V(4).Infoln("Create Project")

	res, err := client.Project_Create(global.ProjectName, global.Description)
	if err != nil {
		klog.Fatal(err)
		return err
	}
	fmt.Println(util.GreenColor(res.Msg))
	var data [][]string
	ss := util.StructToSlice(res.Data)
	data = append(data, ss)

	klog.V(4).Infoln("Create Project, 输出结果")

	output.Write(data, "project", global.NoFormat)
	return nil
}
