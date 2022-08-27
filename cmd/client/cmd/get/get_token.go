package get

import (
	"github.com/bryant-rh/cm/pkg/client"

	"github.com/bryant-rh/cm/pkg/output"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdGetToken() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "token",
		Short:             "查看ServiceAccount Token",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := run_token(global.CMClient)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.Flags().StringVarP(&global.ClusterName, "cluster", "c", global.ClusterName, "指定集群名称")
	cmd.Flags().StringVarP(&global.NameSpace, "namespace", "n", global.NameSpace, "指定namespace")
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func run_token(client *client.CMClient) error {
	klog.V(4).Infoln("Get ServiceAccount Token")

	res, err := client.Sa_GetToken(global.ProjectName, global.ClusterName, global.NameSpace)
	if err != nil {
		klog.Fatal(err)
	}

	var data [][]string
	var ss []string
	ss = append(ss, res.Data)
	data = append(data, ss)

	klog.V(4).Infoln("Get ServiceAccount Token, 输出结果")

	output.Write(data, "token", global.NoFormat)
	return nil
}
