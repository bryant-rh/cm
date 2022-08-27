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

func NewCmdUpdateClster() *cobra.Command {

	cmd := &cobra.Command{
		Use:               "cluster",
		Short:             "更新集群信息",
		PersistentPreRunE: global.PreRun,
		RunE: func(cmd *cobra.Command, args []string) error {

			util.CheckErr(complete_cluster(cmd, args))
			err := run_cluster(global.CMClient)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.Flags().StringVarP(&global.ClusterName, "cluster", "c", global.ClusterName, "指定集群名称")
	cmd.Flags().StringVarP(&global.New_ClusterName, "new_cluster", "", global.New_ClusterName, "新的集群名称")
	cmd.Flags().StringVarP(&global.New_Description, "new_description", "", global.New_Description, "新的集群描述信息")
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")
	return cmd
}

func complete_cluster(cmd *cobra.Command, args []string) error {
	if global.New_ClusterName == "" && global.New_Description == "" {
		return errors.New("请指定参数 --new_cluster , --new_description 以确定需要更新的内容")
	}
	return nil
}

func run_cluster(client *client.CMClient) error {
	klog.V(4).Infoln("Update Cluster")

	if !util.AskForConfirmation("update") {
		os.Exit(1)
	}
	res_id, err := client.Cluster_GetId(global.ProjectName, global.ClusterName)
	if err != nil {
		klog.Fatal(err)
		return err
	}
	res, err := client.Cluster_Update(res_id.Data, global.ProjectName, global.New_ClusterName, global.New_Description)
	if err != nil {
		klog.Fatal(err)
		return err
	}

	klog.V(4).Infoln("Update Cluster, 输出结果")

	fmt.Println(util.GreenColor(res.Msg))

	return nil
}
