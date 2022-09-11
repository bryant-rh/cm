package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/bryant-rh/cm/pkg/util"
)

func NewCmdHxCreateProxy() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "生成kube_proxy部署文件",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if global.EnableDebug { // Enable debug mode if `--enableDebug=true` or `DEBUG=true`.
				global.ProxyClient.SetDebug(true)
			}
			err := global.PreRun(cmd, args)
			if err != nil {
				klog.Fatal(err)
			}

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			//获取集群名称

			err := run_proxy()
			if err != nil {
				return err
			}
			return nil

		},
	}
	cmd.PersistentFlags().BoolVar(&global.EnableDebug, "debug", os.Getenv("DEBUG") == "true", "Enable debug mode")
	cmd.Flags().StringVarP(&global.ProjectName, "project", "p", global.ProjectName, "指定项目名称")
	cmd.Flags().StringVarP(&global.ClusterName, "cluster", "c", global.ClusterName, "指定集群名称(必须指定)")
	cmd.Flags().StringVarP(&global.KUBE_TUNNEL_GATEWAY_HOST, "tunnel_gateway_host", "", global.KUBE_TUNNEL_GATEWAY_HOST, "指定tunnel_gateway 地址(不指定使用默认值)")
	cmd.Flags().StringVarP(&global.KUBE_PROXY_IMAGE, "image", "i", global.KUBE_PROXY_IMAGE, "指定PROXY镜像(不指定使用默认值)")
	cmd.Flags().StringVarP(&global.KUBE_PROXY_NAMESPACE, "namespace", "n", global.KUBE_PROXY_NAMESPACE, "指定kube-proxy部署的namespace (不指定使用默认值)")

	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")

	return cmd
}
func run_proxy() error {

	klog.V(4).Infoln("Get Cluster")

	res, err := global.CMClient.Cluster_List(global.ProjectName, global.ClusterName)
	if err != nil {
		klog.Fatal(err)
	}
	if len(res.Data) == 0 {
		klog.Fatal(res.Msg)

	}

	klog.V(4).Infoln("Create kube-proxy Manifest")
	proxy_clustername := fmt.Sprintf("%s_%s", global.ProjectName, global.ClusterName)

	res_proxy, err := global.ProxyClient.Proxy_Create(global.KUBE_PROXY_IMAGE, proxy_clustername, global.KUBE_TUNNEL_GATEWAY_HOST, global.KUBE_PROXY_NAMESPACE)
	if err != nil {
		klog.Fatal(err)
	}

	klog.V(4).Infoln("Create kube-proxyc Manifest, 输出结果")
	fmt.Println(util.GreenColor("kube-proxy服务, 部署yaml文件如下:"))

	fmt.Println(res_proxy)

	return nil
}
