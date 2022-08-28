package get

import (
	"strings"

	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/util"

	"github.com/bryant-rh/cm/pkg/output"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdGetCluster() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "cluster",
		Short:             "查看集群",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := run_cluster(global.CMClient)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称(不指定默认查看所有)")
	cmd.Flags().StringVarP(&global.ClusterName, "cluster", "c", global.ClusterName, "指定集群名称(不指定默认查看所有)")
	cmd.Flags().StringVarP(&global.Label, "label", "l", global.Label, "Selector (label query) to filter on, supports '=' ,.(e.g. -l key1=value1)")
	return cmd
}

func run_cluster(client *client.CMClient) error {
	var data [][]string
	if global.Label == "" {
		klog.V(4).Infoln("Get Cluster")
		res, err := client.Cluster_List(global.ProjectName, global.ClusterName)
		if err != nil {
			klog.Fatal(err)
		}
		for _, v := range res.Data {
			ss := util.StructToSlice(v)
			data = append(data, ss)
		}
	} else {
		klog.V(4).Infoln("Get Cluster For label")
		var label_key, label_value string
		s := strings.Split(global.Label, "=")
		label_key = s[0]
		label_value = s[1]

		res, err := client.Cluster_label(global.ProjectName, label_key, label_value)
		if err != nil {
			klog.Fatal(err)
		}
		for _, v := range res.Data {
			ss := util.StructToSlice(v)
			data = append(data, ss)
		}
	}

	klog.V(4).Infoln("Get Cluster, 输出结果")

	output.Write(data, "cluster", global.NoFormat)
	return nil
}
