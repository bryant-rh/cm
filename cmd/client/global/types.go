package global

import (
	bytes "bytes"

	github_com_go_courier_helmx_constants "github.com/go-courier/helmx/constants"
	github_com_go_courier_helmx_kubetypes "github.com/go-courier/helmx/kubetypes"
	github_com_go_courier_helmx_spec "github.com/go-courier/helmx/spec"
	github_com_go_courier_sqlx_v2_datatypes "github.com/go-courier/sqlx/v2/datatypes"
	github_com_go_courier_statuserror "github.com/go-courier/statuserror"
)

type BranchAndTag struct {
	Data  []string `json:"data"`
	Total int32    `json:"total"`
}

type BytesBuffer = bytes.Buffer

type CommitSet struct {
	Data  []ProjectSpec `json:"data"`
	Total int32         `json:"total"`
}

type Configuration struct {
	Spec      GithubComGoCourierHelmxSpecSpec `json:"spec"`
	Templates map[string][]uint8              `json:"templates,omitempty"`
}

type GithubComGoCourierHelmxConstantsProtocol = github_com_go_courier_helmx_constants.Protocol

type GithubComGoCourierHelmxConstantsPullPolicy = github_com_go_courier_helmx_constants.PullPolicy

type GithubComGoCourierHelmxKubetypesConfigMapVolumeSource struct {
	GithubComGoCourierHelmxKubetypesKubeLocalObjectReference
}

type GithubComGoCourierHelmxKubetypesCronJobOpts = github_com_go_courier_helmx_kubetypes.CronJobOpts

type GithubComGoCourierHelmxKubetypesDeploymentOpts = github_com_go_courier_helmx_kubetypes.DeploymentOpts

type GithubComGoCourierHelmxKubetypesEmptyDirVolumeSource = github_com_go_courier_helmx_kubetypes.EmptyDirVolumeSource

type GithubComGoCourierHelmxKubetypesExecAction = github_com_go_courier_helmx_kubetypes.ExecAction

type GithubComGoCourierHelmxKubetypesHTTPGetAction = github_com_go_courier_helmx_kubetypes.HTTPGetAction

type GithubComGoCourierHelmxKubetypesHTTPHeader = github_com_go_courier_helmx_kubetypes.HTTPHeader

type GithubComGoCourierHelmxKubetypesHandler = github_com_go_courier_helmx_kubetypes.Handler

type GithubComGoCourierHelmxKubetypesHostPathVolumeSource = github_com_go_courier_helmx_kubetypes.HostPathVolumeSource

type GithubComGoCourierHelmxKubetypesJobOpts = github_com_go_courier_helmx_kubetypes.JobOpts

type GithubComGoCourierHelmxKubetypesKubeLocalObjectReference = github_com_go_courier_helmx_kubetypes.KubeLocalObjectReference

type GithubComGoCourierHelmxKubetypesKubeVolumeSource = github_com_go_courier_helmx_kubetypes.KubeVolumeSource

type GithubComGoCourierHelmxKubetypesPersistentVolumeClaimVolumeSource = github_com_go_courier_helmx_kubetypes.PersistentVolumeClaimVolumeSource

type GithubComGoCourierHelmxKubetypesPodOpts = github_com_go_courier_helmx_kubetypes.PodOpts

type GithubComGoCourierHelmxKubetypesProbeOpts = github_com_go_courier_helmx_kubetypes.ProbeOpts

type GithubComGoCourierHelmxKubetypesSecretVolumeSource = github_com_go_courier_helmx_kubetypes.SecretVolumeSource

type GithubComGoCourierHelmxKubetypesTCPSocketAction = github_com_go_courier_helmx_kubetypes.TCPSocketAction

type GithubComGoCourierHelmxSpecAction struct {
	GithubComGoCourierHelmxKubetypesHandler
}

type GithubComGoCourierHelmxSpecContainer struct {
	GithubComGoCourierHelmxSpecImage
	Args           []string                                 `json:"args,omitempty"`
	Command        []string                                 `json:"command,omitempty"`
	Envs           GithubComGoCourierHelmxSpecEnvs          `json:"envs,omitempty"`
	Lifecycle      *GithubComGoCourierHelmxSpecLifecycle    `json:"lifecycle,omitempty"`
	LivenessProbe  *GithubComGoCourierHelmxSpecProbe        `json:"livenessProbe,omitempty"`
	Mounts         []GithubComGoCourierHelmxSpecVolumeMount `json:"mounts,omitempty"`
	ReadinessProbe *GithubComGoCourierHelmxSpecProbe        `json:"readinessProbe,omitempty"`
	TTY            bool                                     `json:"tty,omitempty"`
	WorkingDir     string                                   `json:"workingDir,omitempty"`
}

type GithubComGoCourierHelmxSpecEnvs = github_com_go_courier_helmx_spec.Envs

type GithubComGoCourierHelmxSpecImage = github_com_go_courier_helmx_spec.Image

type GithubComGoCourierHelmxSpecImagePullSecret = github_com_go_courier_helmx_spec.ImagePullSecret

type GithubComGoCourierHelmxSpecIngressRule = github_com_go_courier_helmx_spec.IngressRule

type GithubComGoCourierHelmxSpecJob struct {
	GithubComGoCourierHelmxSpecPod
	GithubComGoCourierHelmxKubetypesJobOpts
	Cron *GithubComGoCourierHelmxKubetypesCronJobOpts
}

type GithubComGoCourierHelmxSpecLifecycle = github_com_go_courier_helmx_spec.Lifecycle

type GithubComGoCourierHelmxSpecPod struct {
	GithubComGoCourierHelmxSpecContainer
	GithubComGoCourierHelmxKubetypesPodOpts
	Initials []GithubComGoCourierHelmxSpecContainer `json:"initials,omitempty"`
}

type GithubComGoCourierHelmxSpecPort = github_com_go_courier_helmx_spec.Port

type GithubComGoCourierHelmxSpecProbe struct {
	GithubComGoCourierHelmxKubetypesProbeOpts
	Action GithubComGoCourierHelmxSpecAction `json:"action"`
}

type GithubComGoCourierHelmxSpecProject = github_com_go_courier_helmx_spec.Project

type GithubComGoCourierHelmxSpecRequestAndLimit = github_com_go_courier_helmx_spec.RequestAndLimit

type GithubComGoCourierHelmxSpecService struct {
	GithubComGoCourierHelmxSpecPod
	GithubComGoCourierHelmxKubetypesDeploymentOpts
	Ingresses []GithubComGoCourierHelmxSpecIngressRule `json:"ingresses,omitempty"`
	Ports     []GithubComGoCourierHelmxSpecPort        `json:"ports,omitempty"`
}

type GithubComGoCourierHelmxSpecSpec = github_com_go_courier_helmx_spec.Spec

type GithubComGoCourierHelmxSpecVersion = github_com_go_courier_helmx_spec.Version

type GithubComGoCourierHelmxSpecVolume struct {
	GithubComGoCourierHelmxKubetypesKubeVolumeSource
}

type GithubComGoCourierHelmxSpecVolumeMount = github_com_go_courier_helmx_spec.VolumeMount

type GithubComGoCourierHelmxSpecVolumes = github_com_go_courier_helmx_spec.Volumes

type GithubComGoCourierSqlxV2DatatypesMySQLTimestamp = github_com_go_courier_sqlx_v2_datatypes.MySQLTimestamp

type GithubComGoCourierStatuserrorErrorField = github_com_go_courier_statuserror.ErrorField

type GithubComGoCourierStatuserrorErrorFields = github_com_go_courier_statuserror.ErrorFields

type GithubComGoCourierStatuserrorStatusErr = github_com_go_courier_statuserror.StatusErr

type Group struct {
	PrimaryID
	RefGroupName
	OperationTime
}

type GroupSet struct {
	Data  []Group `json:"data"`
	Total int32   `json:"total"`
}

type OperationTime struct {
	CreatedAt GithubComGoCourierSqlxV2DatatypesMySQLTimestamp `json:"createdAt"`
	UpdatedAt GithubComGoCourierSqlxV2DatatypesMySQLTimestamp `json:"updatedAt"`
}

type PrimaryID struct {
}

type Project struct {
	PrimaryID
	ProjectInfo
	OperationTime
}

type ProjectInfo struct {
	RefGroupName
	RefProjectName
}

type ProjectSet struct {
	Data  []Project `json:"data"`
	Total int32     `json:"total"`
}

type ProjectSpec struct {
	PrimaryID
	RefProjectSpec
	ProjectSpecInfo
	OperationTime
}

type ProjectSpecInfo struct {
	RefGroupName
	RefProjectName
	Ref
	// Commit号
	CommitSha string `json:"commitSha"`
	// 默认测试环境部署的配置
	Configuration Configuration `json:"configuration"`
	// 提交的备注信息
	Desc string `json:"desc,omitempty"`
	// 可部署的镜像名称
	Image string `json:"image"`
	// 该服务当前版本依赖的服务包版本
	Refers Refers `json:"refers,omitempty"`
}

type Ref struct {
	RefName string `json:"refName,omitempty"`
}

type RefGroupName struct {
	// 组名称
	GroupName string `json:"groupName"`
}

type RefProjectName struct {
	// 服务名称
	ProjectName string `json:"projectName"`
}

type RefProjectSpec struct {
	// 唯一ID
	ProjectSpecId string `json:"projectSpecId"`
}

type Refer struct {
	ReferName    string `json:"referName"`
	ReferVersion string `json:"referVersion"`
}

type Refers []Refer
