package create

import (
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/output"
	"github.com/bryant-rh/cm/pkg/util"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdCreateSa() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "sa",
		Short:             "创建SaToken",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := run_sa(global.CMClient)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.Flags().StringVarP(&global.ClusterName, "cluster", "c", global.ClusterName, "指定集群名称")
	cmd.Flags().StringVarP(&global.SaName, "saname", "s", global.SaName, "指定token名称")
	cmd.Flags().StringVarP(&global.SaToken, "satoken", "t", global.SaToken, "指定token内容")
	cmd.Flags().StringVarP(&global.NameSpace, "namespace", "n", global.NameSpace, "token管理的namespace")
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("saname")
	cmd.MarkFlagRequired("satoken")
	cmd.MarkFlagRequired("namespace")
	return cmd
}

func run_sa(client *client.CMClient) error {
	klog.V(4).Infoln("Create ServiceAccount Token")

	res, err := client.Sa_Create(global.ProjectName, global.ClusterName, global.SaName, global.SaToken, global.NameSpace)
	if err != nil {
		klog.Fatal(err)
		return err
	}
	fmt.Println(util.GreenColor(res.Msg))
	var data [][]string
	for _, v := range res.Data {
		ss := util.StructToSlice(v)
		data = append(data, ss)

	}

	klog.V(4).Infoln("Create ServiceAccount Token, 输出结果")

	output.Write(data, "sa", global.NoFormat)
	return nil
}
