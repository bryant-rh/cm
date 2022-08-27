package get

import (
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/util"

	"github.com/bryant-rh/cm/pkg/output"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdGetProject() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "project",
		Short:             "查看项目",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := run_project(global.CMClient)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称(不指定默认查看所有)")
	return cmd
}

func run_project(client *client.CMClient) error {
	klog.V(4).Infoln("Get Project")

	res, err := client.Project_List(global.ProjectName)
	if err != nil {
		klog.Fatal(err)
	}

	var data [][]string
	for _, v := range res.Data {
		ss := util.StructToSlice(v)
		data = append(data, ss)

	}

	klog.V(4).Infoln("Get Project, 输出结果")
	
	output.Write(data, "project", global.NoFormat)
	return nil
}
