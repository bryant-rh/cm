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

func NewCmdCreateCluster() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "cluster",
		Short:             "创建集群",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := run_cluster(global.CMClient)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.Flags().StringVarP(&global.ClusterName, "cluster", "c", global.ClusterName, "指定集群名称")
	cmd.Flags().StringVarP(&global.Description, "description", "d", global.Description, "集群描述信息")
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("description")
	return cmd
}

func run_cluster(client *client.CMClient) error {
	klog.V(4).Infoln("Create Cluster")

	res, err := client.Cluster_Create(global.ProjectName, global.ClusterName, global.Description)
	if err != nil {
		klog.Fatal(err)
		return err
	}
	fmt.Println(util.GreenColor(res.Msg))
	var data [][]string
	ss := util.StructToSlice(res.Data)
	data = append(data, ss)

	klog.V(4).Infoln("Create Cluster, 输出结果")

	output.Write(data, "cluster", global.NoFormat)
	return nil
}
