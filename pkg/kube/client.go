package kube

import (
	"github.com/bryant-rh/cm/cmd/client/global"
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

// RestConfig returns a complete rest client config.
func RestConfig(cluster_name string) (*rest.Config, error) {
	if len(global.KUBE_TUNNEL_GATEWAY_HOST) == 0 {
		klog.Fatalf("请在目录:[%s] 创建配置文件: [cm.yaml], 配置 KUBE_TUNNEL_GATEWAY_HOST 或者配置对应环境变量\n", global.Paths.BasePath())

		return nil, rest.ErrNotInCluster
	}

	host := fmt.Sprintf("%s/proxies/%s", global.KUBE_TUNNEL_GATEWAY_HOST, cluster_name)
	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host:        host,
		BearerToken: global.KubeBearerToken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}, nil

}
