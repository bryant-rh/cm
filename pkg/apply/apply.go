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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

		gvr := schema.GroupVersionResource{
			Group:    "",
			Version:  "v1",
			Resource: "namespaces",
		}
		obj, err := dynamicClient.Resource(gvr).Get(context.TODO(), ns, metav1.GetOptions{})
		if err != nil && obj == nil {
			s := getObject("v1", "Namespace", ns)
			_, err := dynamicClient.Resource(gvr).Create(context.TODO(), s, metav1.CreateOptions{})
			if err != nil {
				klog.Fatalf("create Namespace: [%s] error: %v", ns, err)

			}
		}

		fmt.Println(applyStr)
		applyOptions := kube.NewApplyOptions(dynamicClient, discoveryClient)
		if err := applyOptions.Apply(context.TODO(), []byte(applyStr), ns); err != nil {
			klog.Fatalf("apply error: %v", err)
		} else {
			fmt.Println(util.GreenColor("服务部署成功"))

		}
	}
}

func getObject(version, kind, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}
}
