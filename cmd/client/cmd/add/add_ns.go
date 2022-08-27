package add

import (
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/util"

	"github.com/bryant-rh/cm/pkg/output"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdAddNs() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "ns",
		Short:             "添加 ServiceAccount下关联的 NameSpace",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := run_ns(global.CMClient)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.Flags().StringVarP(&global.ClusterName, "cluster", "c", global.ClusterName, "指定集群名称")
	cmd.Flags().StringVarP(&global.SaName, "saname", "s", global.SaName, "指定sa名称")
	cmd.Flags().StringVarP(&global.NameSpace, "namespace", "n", global.NameSpace, "指定namespace,可指定多个,逗号分割,example: kube1,kube2,kube3")
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("saname")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func run_ns(client *client.CMClient) error {
	klog.V(4).Infoln("Add ServiceAccount's NameSpace")

	res, err := client.Sa_AddNs(global.ProjectName, global.ClusterName, global.SaName, global.NameSpace)
	if err != nil {
		klog.Fatal(err)
	}

	var data [][]string
	for _, v := range res.Data {
		ss := util.StructToSlice(v)
		data = append(data, ss)

	}

	klog.V(4).Infoln("Add ServiceAccount's NameSpace, 输出结果")

	output.Write(data, "sa", global.NoFormat)
	return nil
}
