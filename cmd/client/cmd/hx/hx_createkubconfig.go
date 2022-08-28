package hx

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/bryant-rh/cm/pkg/util"
)

type kubeconfig struct {
	ServerUrl   string
	ClusterName string
	UserName    string
	Token       string
}

func NewCmdHxCreateKubeConfig() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "生成kubeConfig文件",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			err := global.PreRun(cmd, args)
			if err != nil {
				klog.Fatal(err)
			}

		},
		RunE: func(cmd *cobra.Command, args []string) error {

			err := run_kube()
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
	cmd.Flags().StringVarP(&global.NameSpace, "namespace", "n", global.NameSpace, "指定namespace")
	cmd.MarkFlagRequired("project")
	cmd.MarkFlagRequired("cluster")
	cmd.MarkFlagRequired("namespace")

	return cmd
}
func run_kube() error {

	klog.V(4).Infoln("Get ServiceAccount Token")

	res, err := global.CMClient.Sa_GetToken(global.ProjectName, global.ClusterName, global.NameSpace)
	if err != nil {
		klog.Fatal(err)
	}

	token := res.Data
	clusterName := fmt.Sprintf("%s_%s", global.ProjectName, global.ClusterName)
	serverUrl := fmt.Sprintf("%s/proxies/%s", global.KUBE_TUNNEL_GATEWAY_HOST, clusterName)
	userName := fmt.Sprintf("%s-%s-rw", global.ProjectName, global.ClusterName)

	fmt.Println(serverUrl)
	kc := kubeconfig{}
	kc.SetDefaults(token, clusterName, userName, serverUrl)
	output, err := kc.GetKubeConfig()
	if err != nil {
		klog.Fatal(err)
	}
	configname := fmt.Sprintf("%s.config", clusterName)

	_ = util.WriteToFile(configname, output, func(v interface{}) ([]byte, error) {
		return v.([]byte), nil
	})
	fmt.Println(util.GreenColor(fmt.Sprintf("生成kubeconfig: [%s] 成功", configname)))
	fmt.Println(util.GreenColor(fmt.Sprintf("可通过 kubectl --kubeconfig %s  访问k8s集群资源", configname)))

	return nil
}

func (c *kubeconfig) SetDefaults(token, clusterName, userName, serverUrl string) {
	c.ClusterName = clusterName
	c.ServerUrl = serverUrl
	c.Token = token
	c.UserName = userName

}

func (c *kubeconfig) GetKubeConfig() ([]byte, error) {
	klog.V(4).Infoln("Create kubeconfig")
	t, err := template.New("kubeconfig").Parse(kubeconfigTemplate)
	if err != nil {
		klog.Errorf("Read Template Error %s", err.Error())
		return nil, err
	}
	buff := bytes.NewBufferString("")
	err = t.ExecuteTemplate(buff, "kubeconfig", c)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

var kubeconfigTemplate = `
apiVersion: v1
clusters:
  - cluster:
      insecure-skip-tls-verify: true
      # https://{gatewayUrl}/proxies/{cluster}
      server: {{ .ServerUrl }}
    name: {{ .ClusterName }}
contexts:
  - context:
      cluster: {{ .ClusterName }}
      user: {{ .UserName }}
    name: token-context
current-context: token-context
kind: Config
users:
  - name: {{ .UserName }}
    user:
      token: {{ .Token }}

`
