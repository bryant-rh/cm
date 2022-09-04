package apply

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/bryant-rh/cm/pkg/hxctx"
	"github.com/bryant-rh/cm/pkg/kube"
	"github.com/bryant-rh/cm/pkg/util"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
)

func RunApply(workspace *hxctx.Workspace, dryRun bool, dynamicClient dynamic.Interface, discoveryClient *discovery.DiscoveryClient) {
	ns, data := hxctx.CommandForDeploy(workspace)
	if dryRun {
		fmt.Println(util.GreenColor(fmt.Sprintf("namespace:[%s], 渲染yaml文件如下:", ns)))
		fmt.Println(data)
	} else {
		// if global.ProjectName == "" || global.ClusterName == "" || global.NameSpace == "" {
		// 	cmd.Help()
		// 	klog.Fatal(util.RedColor("需要 -p 指定项目名称, -c 指定集群名称, -n 指定namespace"))
		// 	//fmt.Println(util.RedColor("需要 -p 指定项目名称, -c 指定集群名称, -n 指定namespace"))

		// }
		// //获取集群名称和token
		// err := global.PreRun(cmd, args)
		// if err != nil {
		// 	klog.Fatal(err)
		// }

		// klog.V(4).Infoln("Get ServiceAccount Token")

		// res, err := global.CMClient.Sa_GetToken(global.ProjectName, global.ClusterName, global.NameSpace)
		// if err != nil {
		// 	klog.Fatal(err)
		// }
		// global.KubeBearerToken = res.Data

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

		// proxy_clustername := fmt.Sprintf("%s_%s", global.ProjectName, global.ClusterName)
		// config, err := kube.RestConfig(proxy_clustername)
		// if err != nil {
		// 	klog.Fatal(err)
		// }
		// dynamicClient, err := dynamic.NewForConfig(config)
		// if err != nil {
		// 	panic(err.Error())
		// }
		// discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
		// if err != nil {
		// 	panic(err.Error())
		// }

		fmt.Println(applyStr)
		applyOptions := kube.NewApplyOptions(dynamicClient, discoveryClient)
		if err := applyOptions.Apply(context.TODO(), []byte(applyStr), ns); err != nil {
			klog.Fatalf("apply error: %v", err)
		} else {
			fmt.Println(util.GreenColor("服务部署成功"))

		}
	}
}
