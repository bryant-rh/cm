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

func NewCmdUpdateSa() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "sa",
		Short:             "更新sa信息",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {

			util.CheckErr(complete_sa(cmd, args))
			err := run_sa(global.CMClient)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.Flags().StringVarP(&global.ClusterName, "cluster", "c", global.ClusterName, "指定集群名称")
	cmd.Flags().StringVarP(&global.SaName, "saname", "s", global.SaName, "token名称")
	cmd.Flags().StringVarP(&global.SaToken, "satoken", "t", global.SaToken, "token内容")
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("saname")
	return cmd
}

func complete_sa(cmd *cobra.Command, args []string) error {
	if global.SaToken == ""  {
		return errors.New("请指定参数 --satoken | -t 以确定需要更新的内容")
	}
	return nil
}

func run_sa(client *client.CMClient) error {
	klog.V(4).Infoln("Update ServiceAccount Token")

	if !util.AskForConfirmation("update") {
		os.Exit(1)
	}

	res, err := client.Sa_Update(global.ProjectName, global.ClusterName, global.SaName, global.SaToken)
	if err != nil {
		klog.Fatal(err)
		return err
	}

	klog.V(4).Infoln("Update ServiceAccount Token, 输出结果")

	fmt.Println(util.GreenColor(res.Msg))

	return nil
}
