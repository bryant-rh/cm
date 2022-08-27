package update

import (
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/util"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdUpdateProject() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "project",
		Short:             "更新项目信息",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {

			util.CheckErr(complete_project(cmd, args))
			err := run_project(global.CMClient)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.Flags().StringVarP(&global.New_ProjectName, "new_project", "", global.New_ProjectName, "新的项目名称")
	cmd.Flags().StringVarP(&global.New_Description, "new_description", "", global.New_Description, "新的项目描述信息")
	cmd.MarkFlagRequired("project")
	return cmd
}

func complete_project(cmd *cobra.Command, args []string) error {
	if global.New_ProjectName == "" && global.New_Description == "" {
		return errors.New("请指定参数 --new_project , --new_description 以确定需要更新的内容")
	}
	return nil
}

func run_project(client *client.CMClient) error {
	klog.V(4).Infoln("Update Project")

	if !util.AskForConfirmation("update") {
		os.Exit(1)
	}
	res_id, err := client.Project_GetId(global.ProjectName)
	if err != nil {
		klog.Fatal(err)
		return err
	}
	res, err := client.Project_Update(res_id.Data, global.New_ProjectName, global.New_Description)
	if err != nil {
		klog.Fatal(err)
		return err
	}

	klog.V(4).Infoln("Update Project, 输出结果")

	fmt.Println(util.GreenColor(res.Msg))

	return nil
}
