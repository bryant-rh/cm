package get

import (
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/util"

	"github.com/bryant-rh/cm/pkg/output"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdGetSa() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "sa",
		Short:             "查看ServiceAccount",
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
	cmd.Flags().StringVarP(&global.SaName, "saname", "s", global.SaName, "指定sa名称(不指定默认查看所有)")
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")

	return cmd
}

func run_sa(client *client.CMClient) error {
	klog.V(4).Infoln("Get ServiceAccount")

	res, err := client.Sa_List(global.ProjectName, global.ClusterName, global.SaName)
	if err != nil {
		klog.Fatal(err)
	}

	var data [][]string
	for _, v := range res.Data {
		ss := util.StructToSlice(v)
		data = append(data, ss)

	}

	klog.V(4).Infoln("Get ServiceAccount, 输出结果")

	output.Write(data, "sa", global.NoFormat)
	return nil
}
