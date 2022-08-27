package global

import (
	"github.com/bryant-rh/cm/pkg/util"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

func PreRun(cmd *cobra.Command, _ []string) error {

	if len(CM_SERVER_BASEURL) == 0 || len(CM_SERVER_USERNAME) == 0 || len(CM_SERVER_PASSWORD) == 0 {
		klog.Fatalf("请在目录:[%s] 创建配置文件: [cm.yaml], 配置CM_SERVER_BASEURL、CM_SERVER_USERNAME、CM_SERVER_PASSWORD 或者配置对应环境变量\n", Paths.BasePath())

	}
	if EnableDebug { // Enable debug mode if `--enableDebug=true` or `DEBUG=true`.
		CMClient.SetDebug(true)
	}
	klog.V(4).Infof("CM_SERVER_BASEURL: %s; CM_SERVER_USERNAME: %s; CM_SERVER_PASSWORD: %s ", CM_SERVER_BASEURL, CM_SERVER_USERNAME, CM_SERVER_PASSWORD)

	authFile := fmt.Sprintf("%s/%s", Paths.BasePath(), TokenFile)
	if !util.Exists(authFile) {
		klog.V(4).Infoln("登录生成token")

		res, err := CMClient.User_Login(CM_SERVER_USERNAME, CM_SERVER_PASSWORD)
		if err != nil {
			klog.Fatal(err)
		}
		klog.V(4).Infof("登录生成token,并保存至文件: [%s]", authFile)

		err = ioutil.WriteFile(authFile, []byte(res.Data), 0644)
		if err != nil {
			klog.Fatal(err)
		}
		Token = res.Data
	} else {
		klog.V(4).Infoln("检测token是否过期")
		auth_token, err := ioutil.ReadFile(authFile)
		if err != nil {
			klog.Fatal(err)
		}
		_, err = CMClient.User_VerifyToken(string(auth_token))
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				klog.V(4).Infoln("检测token已过期,重新登录生成token")
				new_res, err := CMClient.User_Login(CM_SERVER_USERNAME, CM_SERVER_PASSWORD)
				if err != nil {
					klog.Fatal(err)
				}
				klog.V(4).Infof("生成token,并保存至文件: [%s]", authFile)
				err = ioutil.WriteFile(authFile, []byte(new_res.Data), 0644)
				if err != nil {
					klog.Fatal(err)
				}
				Token = new_res.Data
			} else {
				klog.Fatal(err)

			}
		} else {
			klog.V(4).Infoln("token未过期")
			Token = string(auth_token)

		}

	}
	klog.V(4).Infof("token: %s", Token)
	klog.V(4).Infoln("执行 LoginWithToken")

	CMClient.LoginWithToken(Token)

	return nil
}

func InClusterConfig(cluster_name, bearertoken string) (*rest.Config, error) {

	if len(KUBE_TUNNEL_GATEWAY_HOST) == 0 {
		klog.Fatalf("请在目录:[%s] 创建配置文件: [cm.yaml], 配置 KUBE_TUNNEL_GATEWAY_HOST 或者配置对应环境变量\n", Paths.BasePath())

		return nil, rest.ErrNotInCluster
	}

	host := fmt.Sprintf("%s/proxies/%s", KUBE_TUNNEL_GATEWAY_HOST, cluster_name)
	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host:        host,
		BearerToken: string(bearertoken),
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}, nil
}
