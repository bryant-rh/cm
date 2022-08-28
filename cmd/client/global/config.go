package global

import (
	"os"

	"github.com/bryant-rh/cm/pkg/client"
	"github.com/bryant-rh/cm/pkg/environment"

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
	KubeBearerToken          string
	KUBE_PROXY_IMAGE         = "bryantrh/kube-tunnel-gateway:v1.0"
	KUBE_PROXY_NAMESPACE     = "kube-proxy"

	Ctx               = &Context{}
	HelmxRootFile     string
	LivenessCheckSkip bool
	ImagePullSecret   = os.Getenv(spec.EnvKeyImagePullSecret)
	SkipLivenessCheck bool
)

//k8s相关
