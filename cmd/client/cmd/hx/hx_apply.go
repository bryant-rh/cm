package hx

import (
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/bryant-rh/cm/pkg/apply"
	"github.com/bryant-rh/cm/pkg/hxctx"
	"github.com/bryant-rh/cm/pkg/util"
)

var (
	deployOpt = hxctx.DeployOpt{}
)

func NewCmdHxApply(ctx *hxctx.Context) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "通过文件名或标准输入流(stdin)对资源进行配置",
		Run: func(cmd *cobra.Command, args []string) {
			klog.V(4).Infoln("hx apply")
			if !deployOpt.DryRun {
				if global.ProjectName == "" || global.ClusterName == "" {
					cmd.Help()
					klog.Fatal(util.RedColor("需要 -p 指定项目名称, -c 指定集群名称 "))
					//fmt.Println(util.RedColor("需要 -p 指定项目名称, -c 指定集群名称 "))

				}
				//获取集群名称和token
				err := global.PreRun(cmd, args)
				if err != nil {
					klog.Fatal(err)
				}

				ctx.Range(func(w *hxctx.Workspace) {
					apply.RunApply(w, deployOpt.DryRun)

				})

			} else {
				ctx.Range(func(w *hxctx.Workspace) {
					apply.RunApply(w, deployOpt.DryRun)

				})
			}

		},
	}
	cmd.PersistentFlags().BoolVarP(&global.SkipLivenessCheck, "skip-liveness-check", "", false, "skip liveness check")
	cmd.Flags().BoolVarP(&deployOpt.DryRun, "dry-run", "", false, "dry run used like helm debug")
	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.Flags().StringVarP(&global.ClusterName, "cluster", "c", global.ClusterName, "指定集群名称")
	cmd.Flags().StringVarP(&global.NameSpace, "namespace", "n", global.NameSpace, "指定namespace")

	return cmd
}
