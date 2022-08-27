package add

import (
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/util"
	"fmt"
	"strings"

	"github.com/bryant-rh/cm/pkg/output"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdAddLabel() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "label",
		Short:             "添加 Cluster下关联的 Label",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := run_label(global.CMClient)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.Flags().StringVarP(&global.ClusterName, "cluster", "c", global.ClusterName, "指定集群名称")
	cmd.Flags().StringVarP(&global.Label, "label", "l", global.Label, "Selector (label query) to filter on, supports '=' ,.(e.g. -l key1=value1)")
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("label")

	return cmd
}

func run_label(client *client.CMClient) error {
	klog.V(4).Infoln("Add Cluster's Label")
	var label_key, label_value string
	if global.Label != "" {
		s := strings.Split(global.Label, "=")
		label_key = s[0]
		label_value = s[1]

	}
	res, err := client.Label_Create(global.ProjectName, global.ClusterName, label_key, label_value)
	if err != nil {
		klog.Fatal(err)
	}

	fmt.Println(util.GreenColor(res.Msg))
	var data [][]string
	for _, v := range res.Data {
		ss := util.StructToSlice(v)
		data = append(data, ss)

	}
	klog.V(4).Infoln("Add Cluster's Label, 输出结果")

	output.Write(data, "label", global.NoFormat)
	return nil
}
