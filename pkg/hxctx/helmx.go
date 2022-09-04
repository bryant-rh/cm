package hxctx

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bryant-rh/cm/cmd/client/global"
	"github.com/go-courier/helmx"
	"github.com/go-courier/helmx/kubetypes"
	"github.com/go-courier/helmx/spec"
	"github.com/sirupsen/logrus"
)

type DeployOpt struct {
	DryRun bool
}

func CommandForDeploy(w *Workspace) (namespace, data string) {
	namespace = strings.ToLower(w.Project.Group)

	if w.DeployEnv != "" && w.DeployEnv != "online" {
		namespace = namespace + "--" + w.DeployEnv
	}

	namespace = strings.ToLower(namespace)

	// ret, _ := exec.Command("sh", "-c", fmt.Sprintf("kubectl get namespace %s -o yaml", namespace)).Output()

	// istioEnabled := bytes.Contains(ret, []byte("istio-injection: enabled"))
	istioEnabled := false
	buf := bytes.NewBuffer(nil)

	executeAll(buf, w, istioEnabled)

	if buf.Len() == 0 {
		return namespace, `echo nothing to deploy`
	}

	data = quoteWith(buf.String())
	return namespace, data

}

// 替换所有gitlab-ci内置的CI变量
func replaceWithGitlabEnv(s string) string {
	stemp := s

	re, _ := regexp.Compile("(CI_[A-Z_]*)|GITLAB_USER_NAME")

	for _, env := range re.FindAllString(stemp, -1) {
		e := os.Getenv(env)

		if e != "" {
			stemp = strings.ReplaceAll(stemp, env, replaceControlS(e))
		}
	}
	return stemp
}

// 替换所有特殊字符
func replaceControlS(s string) string {
	reControl, err := regexp.Compile("[\\\\\"']")
	if err != nil {
		logrus.Errorln(err.Error())
		return s
	}
	return reControl.ReplaceAllString(s, "*")
}

func quoteWith(s string) string {
	buf := bytes.NewBuffer(nil)

	for _, b := range s {
		if b == '\'' {
			buf.WriteRune('\\')
			buf.WriteRune('\'')
			continue
		}

		buf.WriteRune(b)
	}
	return replaceWithGitlabEnv(buf.String())
}

func executeAll(writer io.Writer, w *Workspace, istioEnabled bool) {
	hx := helmx.NewHelmX()
	hx.Spec = w.Spec

	hx.Envs = map[string]string{}

	// set default value
	// set job default values
	// TTLSecondsAfterFinished is disabled before v1.21
	// https://v1-18.docs.kubernetes.io/docs/concepts/workloads/controllers/ttlafterfinished/
	// for key := range hx.Spec.Jobs {
	// 	job := hx.Spec.Jobs[key]
	// 	if job.TTLSecondsAfterFinished == nil {
	// 		job.TTLSecondsAfterFinished = ptrInt32(10)

	// 	}
	// 	if job.BackoffLimit == nil {
	// 		job.BackoffLimit = ptrInt32(1)
	// 	}
	// 	hx.Spec.Jobs[key] = job
	// }

	for k, v := range w.EnvVarSet.Values {
		hx.Envs[k] = v
	}
	setDefaults(&hx.Spec)

	hx.AddTemplate("pullSecret", pullSecret)
	hx.AddTemplate("serviceAccount", serviceAccount)
	hx.AddTemplate("service", service)

	if istioEnabled {
		hx.AddTemplate("virtualService", virtualService)
	} else {
		hx.AddTemplate("ingress", ingress)
	}

	hx.AddTemplate("deployment", deployment)
	hx.AddTemplate("job", job)
	hx.AddTemplate("cronJob", cronJob)

	templatePath := "templates"

	if isPathExist(templatePath) {

		files, err := ioutil.ReadDir(templatePath)
		if err != nil {
			panic(err)
		}

		w.Templates = make(map[string][]byte)

		for _, f := range files {
			if !f.IsDir() {
				data, err := ioutil.ReadFile(filepath.Join(templatePath, f.Name()))
				if err != nil {
					panic(err)
				}
				hx.AddTemplate(strings.Replace(filepath.Base(f.Name()), filepath.Ext(f.Name()), "", -1), string(data))
				// 存取自定义模版
				w.Templates[f.Name()] = data
			}
		}
	}

	err := hx.ExecuteAll(writer, &hx.Spec)
	if err != nil {
		panic(err)
	}

	// skip liveness check
	skipLivenessCheck(&hx.Spec)
}

func skipLivenessCheck(spec *spec.Spec) {

	if global.SkipLivenessCheck {
		return
	}

	// spew.Dump(spec.Service.Pod.LivenessProbe)
	// spew.Dump(spec.Service.Pod.ReadinessProbe)

	if spec.Service.Pod.LivenessProbe == nil || spec.Service.Pod.ReadinessProbe == nil {
		log.Fatal(`################################################################
		livenessProbe 和 readinessProbe 必须存在。
		更多 CICD 流程， 
			参考文档: https://rockontrol.yuque.com/spec/eb4smu/cz2in7#pYSGw
		################################################################
		`)
	}

}
func setDefaults(s *spec.Spec) {
	if s.Service != nil && s.Service.DNSConfig == nil {
		s.Service.DNSConfig = &kubetypes.DNSConfig{
			Options: []kubetypes.KubeOption{
				{
					Name:  "ndots",
					Value: "2",
				},
			},
		}
		s.Service.DNSPolicy = "ClusterFirst"
		s.Service.Container.ImagePullPolicy = "IfNotPresent"

		if s.Service.ImagePullSecret == nil {
			v := os.Getenv("IMAGE_PULL_SECRET")
			if v != "" {
				s.Service.ImagePullSecret, _ = spec.ParseImagePullSecret(v)
			}
		}
	}

	// 检查资源属性
	if s.Resources == nil {
		s.Resources = spec.Resources{}
	}

	if _, ok := s.Resources["cpu"]; !ok {
		s.Resources["cpu"], _ = spec.ParseRequestAndLimit("10/500m")
	}

	if _, ok := s.Resources["memory"]; !ok {
		s.Resources["memory"], _ = spec.ParseRequestAndLimit("10/1024Mi")
	}
}

var (
	serviceAccount = `
{{ if ( len .Service.ServiceAccountRoleRules ) }}

--- 
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ ( .Service.ServiceAccountName ) }}
rules:
{{ spaces 2 | toYamlIndent ( toKubeRoleRules . )}}

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ ( .Service.ServiceAccountName ) }}

---

apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ ( .Service.ServiceAccountName ) }}
subjects:
  - kind: ServiceAccount
    name: {{ ( .Service.ServiceAccountName ) }}
roleRef:
  kind: Role
  name: {{ ( .Service.ServiceAccountName ) }}
  apiGroup: rbac.authorization.k8s.io

{{ end }}
`

	service = `
{{ if ( and ( exists .Service ) ( gt ( len .Service.Ports ) 0 ) ) }}

--- 
apiVersion: v1
kind: Service
metadata:
  name: {{ ( .Project.FullName ) }}
  annotations: 
    helmx/project: >-
      {{ toJson .Project }}
    helmx/upstreams: {{ join .Upstreams "," | quote }}
spec:
  selector:
    srv: {{ ( .Project.FullName ) }}
{{ spaces 2 | toYamlIndent ( toKubeServiceSpec . )  }}
{{ end }}
`

	deployment = `
{{ if ( exists .Service ) }}

---
apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    "git.querycap.com/last-commit-by": "CI_COMMIT_REF_SLUG/GITLAB_USER_NAME/CI_COMMIT_SHA/CI_COMMIT_MESSAGE"
  name: {{ ( .Project.FullName ) }}
  labels:
    app: {{ ( .Project.FullName ) }}
    version: {{ ( .Project.Version ) }}
spec:
  selector:
    matchLabels:
      srv: {{ ( .Project.FullName ) }}
{{ spaces 2 | toYamlIndent ( toKubeDeploymentSpec . )  }}
{{ end }}
`
	job = `
{{ $spec := .}}
{{ range $name, $job := .Jobs }}
{{ if (not (exists $job.Cron)) }}
---
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ ( $spec.Project.FullName ) }}--{{ $name }}
spec:
{{ spaces 2 | toYamlIndent ( toKubeJobSpec $spec $job )  }}
{{ end }}
{{ end }}
`

	cronJob = `
{{ $spec := .}}
{{ range $name, $job := .Jobs }}
{{ if (exists $job.Cron) }}
---
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: {{ ( $spec.Project.FullName ) }}--{{ $name }}
spec:
{{ spaces 2 | toYamlIndent ( toKubeCronJobSpec $spec $job )  }}
{{ end }}
{{ end }}
`

	ingress = `
{{ if ( gt ( len .Service.Ingresses ) 0 ) }}

--- 
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{ ( .Project.FullName ) }}
  annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
{{ spaces 2 | toYamlIndent ( toKubeIngressSpec . )}}
{{ end }}
`

	virtualService = `
{{ if ( and ( exists .Service ) ( gt ( len .Service.Ports ) 0 ) ) }}

---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{ ( .Project.FullName ) }}
spec:
  gateways:
    - default-gateway
  hosts:
    - "*"
  http:
    - route:
        - destination:
            host: {{ ( .Project.FullName ) }}
{{ end }}
`

	pullSecret = `
{{ if ( and ( exists .Service ) ( exists .Service.ImagePullSecret ) ) }}

--- 
apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: {{ ( .Service.ImagePullSecret.Name ) }}
data:
  .dockerconfigjson: "{{ ( .Service.ImagePullSecret.Base64EncodedDockerConfigJSON ) }}"

{{ end }}
`
)
