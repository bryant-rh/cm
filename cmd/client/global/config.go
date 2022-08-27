package global

import (
	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/environment"
	"os"

	"github.com/go-courier/helmx/spec"
)

//全局配置
var (
	// CM_SERVER_BASEURL = os.Getenv("CM_SERVER_BASEURL")
	// CM_SERVER_USERNAME = os.Getenv("CM_SERVER_USERNAME")
	// CM_SERVER_PASSWORD = os.Getenv("CM_SERVER_PASSWORD")
	CM_SERVER_BASEURL  string
	CM_SERVER_USERNAME string
	CM_SERVER_PASSWORD string
	Token              string
	TokenFile          = "token.json"
	CMClient           *client.CMClient
	ProxyClient        *client.CMClient
	Paths              environment.Paths
)

//命令参数
var (
	ProjectName     string
	New_ProjectName string

	ClusterName     string
	New_ClusterName string

	Description     string
	New_Description string

	SaName string

	SaToken string

	NameSpace string

	Label string

	EnableDebug bool

	NoFormat bool
)

//hx变量
var (
	KUBE_TUNNEL_GATEWAY_HOST string
//	KUBE_PROXY_CLUSTER       string
	KubeBearerToken          string
	//KubeBearerToken          = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImVCRGNaanExanNmZVl3enE1MmxzUEhXTUpDamgtTjJDYkwyNmU1SU82Z28ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlLXByb3h5LXRva2VuLW5yaHpyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6Imt1YmUtcHJveHkiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJjZGY3NDJkMC03YTIzLTRjN2ItYWFhOS00YmY4ODNmODA2NDUiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06a3ViZS1wcm94eSJ9.F5pzq2AtqlxWevqOjNET8Et652G-_3uxZCTXcSnlm_T44-vbAalH0SO0NhebsNmQgxjgZMj0sbWj9sWxV12HsU3FfrL6CAjn6OL69njAiv8YfDpUHh3Ptke9GmuAO5OLsgjD-ktGLP_H7UQJsWA8l8cQJ8J8cSCQkNP8E3b9FFQIseYtlM1DTj1U4LqgHcTZE4XQy0WEctNXULfI9JvTw1xZ_8LQnvv7b5NO_VteMoIOEqkjrkqjs3MqpvNgfKQS02Yx2hhyxJU2imLJ2AiCm46DS0fQcsHjJN112bdvxzqaDdZ3wVkf4p_lCTZzkApcMsEyhO34654JrnkChM7xzQ"
	KUBE_PROXY_IMAGE     = "bryantrh/kube-tunnel-gateway:v1.0"
	KUBE_PROXY_NAMESPACE = "kube-proxy"

	Ctx               = &Context{}
	HelmxRootFile     string
	LivenessCheckSkip bool
	ImagePullSecret   = os.Getenv(spec.EnvKeyImagePullSecret)
	SkipLivenessCheck bool
)

//k8s相关
