package delete

import (
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/util"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
)

func NewCmdDelSa() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "sa",
		Short:             "删除sa",
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
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("saname")
	return cmd
}

func run_sa(client *client.CMClient) error {
	klog.V(4).Infoln("Delete ServiceAccount Token")

	if !util.AskForConfirmation("delete") {
		os.Exit(1)
	}
	res, err := client.Sa_Delete(global.ProjectName, global.ClusterName, global.SaName)
	if err != nil {
		klog.Fatal(err)
	}

	klog.V(4).Infoln("Delete ServiceAccount Token, 输出结果")

	fmt.Println(util.GreenColor(res.Msg))

	return nil
}
