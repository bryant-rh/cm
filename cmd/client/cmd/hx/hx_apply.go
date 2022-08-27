package hx

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/bryant-rh/cm/pkg/helmx"
	"github.com/bryant-rh/cm/pkg/kube"

	"github.com/bryant-rh/cm/pkg/util"
)

var (
	deployOpt = helmx.DeployOpt{}
)

func NewCmdHxApply() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "通过文件名或标准输入流(stdin)对资源进行配置",
		Run: func(cmd *cobra.Command, args []string) {
			klog.V(4).Infoln("hx apply")
			ns, data := helmx.CommandForDeploy(global.Ctx)
			if deployOpt.DryRun {
				fmt.Println(util.GreenColor(fmt.Sprintf("namespace:[%s], 渲染yaml文件如下:", ns)))
				fmt.Println(data)
			} else {
				if global.ProjectName == "" || global.ClusterName == "" || global.NameSpace == "" {
					cmd.Help()
					klog.Fatal(util.RedColor("需要 -p 指定项目名称, -c 指定集群名称, -n 指定namespace"))
					//fmt.Println(util.RedColor("需要 -p 指定项目名称, -c 指定集群名称, -n 指定namespace"))

				}
				//获取集群名称和token
				err := global.PreRun(cmd, args)
				if err != nil {
					klog.Fatal(err)
				}

				klog.V(4).Infoln("Get ServiceAccount Token")

				res, err := global.CMClient.Sa_GetToken(global.ProjectName, global.ClusterName, global.NameSpace)
				if err != nil {
					klog.Fatal(err)
				}
				global.KubeBearerToken = res.Data

				//去除空行
				re := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
				s := re.ReplaceAllString(data, "")

				applyStr := ""
				if strings.HasPrefix(s, "---") {
					applyStr = strings.Replace(s, "---", "", 1)

				}
				//判断namespace
				if global.NameSpace != "" {
					ns = global.NameSpace
				}

				proxy_clustername := fmt.Sprintf("%s_%s", global.ProjectName, global.ClusterName)
				config, err := kube.RestConfig(proxy_clustername)
				if err != nil {
					klog.Fatal(err)
				}
				dynamicClient, err := dynamic.NewForConfig(config)
				if err != nil {
					panic(err.Error())
				}
				discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
				if err != nil {
					panic(err.Error())
				}

				fmt.Println(applyStr)
				applyOptions := kube.NewApplyOptions(dynamicClient, discoveryClient)
				if err := applyOptions.Apply(context.TODO(), []byte(applyStr), ns); err != nil {
					klog.Fatalf("apply error: %v", err)
				} else {
					fmt.Println(util.GreenColor("服务部署成功"))

				}

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
