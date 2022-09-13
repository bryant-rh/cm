# 1.简介
cm，编写的一个用于多集群管理的CICD工具，包括项目初始化，生成模板Dockerfile，自动渲染k8s 部署文件，管理多集群通道token，部署至多集群中等功能。方便快速实现项目容器化，并部署至k8s集群中。

# 2.痛点
项目CICD中可能存在如下痛点：

1. CI过程中，项目编译成镜像时，需要依赖docker.sock，而如今continerd做容器运行时，不再适配。
2. 如今适应国产化，以及降低成本，ARM架构的服务器已变得流行，所以编译镜像需支持双架构(amd64/arm64)或者多架构。
3. 需要维护复杂的k8s yaml部署文件。
4. 当没有配置中心时，项目中的配置信息，无法方便的以环境变量的方式在容器中进行调用，而且还存在不同分支，配置不同值的调试需求。
5. CD过程中，一个项目可能存在多个k8s集群，且不同团队可能需要对不同的namespace具有发布权限
6. 当需要部署的k8s集群跟gitlab内网不通，即无法管控远程k8s集群。

# 3.方案
## 3.1 说明
  + 针对上面的痛点1、2  我们采用在buildx 容器中使用docker buildx 工具即可解决

  + 针对上面的痛点3、4 我们采用通过封装一个工具，用一个简易的yml 文件，来渲染k8s资源，满足一些基本的需求，开发人员不用关注k8s yaml 复杂的编写方式、运维人员也不再为每一个项目都维护一个yaml文件。且在渲染过程中，通过配置文件的方式，自动将不同的分支的配置文件，注入到k8s 环境变量

  + 针对上面的痛点5、6 我们借助一个自行编写的kube-tunnel 服务，用于在管控集群和子集群直接建立一条隧道，来访问子集群，同时编写cm-server服务，来管理一个项目存在多个集群，一个集群下多个namespace的token，如图：

## 3.2 部署

### 3.2.1 部署cm-server
cm-server 是用于管理集群通道Token的服务端，可以创建项目，集群关联到项目等

```Bash
cd deploy/cm-server

#部署mysql
kubectl apply -f mysql5.7-deploy.yaml

# 导入sql
cm-server.sql

#部署cm-server
kubectl apply -f cm-server-deploy.yaml

```

### 3.2.2 部署kube-tunnel-gateway
kube-tunnel-gateway 是用于创建和需要管控k8s集群的服务端，用于建立通道

#### 1. 生成PrivateKey
```Bash
openssl genrsa -out private.key
echo $(cat private.key|base64)|base64
```

#### 2. 生成PublicKey
```Bash
openssl rsa -in private.key -pubout -outform PEM -out public.pem
cat public.pem|base64
```

#### 3. 用上面生成的PrivateKey 和PublicKey 替换部署文件中 Secret 的对应字段，部署
[备注] : 必须 https 访问

```Bash
cd deploy/kube-tunnel
kubectl apply -f kube-tunnel-gateway.yaml
```

### 3.2.3 创建项目、集群相关信息
#### 1. 配置cm工具
在用户家目录.cm 目录下创建cm.yaml，配置如下值，或者创建相应环境变量也可，程序会自动读取。
```BASH
#上面部署的cm-server的访问地址
CM_SERVER_BASEURL: "http://cm-server.example.com/api/v1"
#cm-server支持了鉴权，这里配置对应的账号密码，上面部署cm-server时，可配置初始账号和密码，服务启动会自动进行注册
CM_SERVER_USERNAME: "admin"
CM_SERVER_PASSWORD: "123456"
#对应上面部署的kube-tunel-gateway 服务，注意必须https
KUBE_TUNNEL_GATEWAY_HOST: "https://kube-gateway.example.com"
```
然后即可下载cm客户端进行使用

#### 2. 创建项目
```BASH
#创建
cm create project -p test -d 测试项目
#查看
cm get project
```

#### 3.项目下创建集群
```BASH
#创建
cm create cluster -p test -c vm-test -d 测试集群
#查看
cm get cluster -p test
```

#### 4. 生成kube-proxy部署文件，在子集群进行部署
```BASH
#对应上面创建的项目，集群生成部署文件
cm proxy -p test -c vm-test
```
输出内容：
```
kube-proxy服务, 部署yaml文件如下:
 
 
---
apiVersion: v1
kind: Service
metadata:
  name: kube-proxy
  namespace: kube-proxy
spec:
  selector:
    app: kube-proxy
  type: ClusterIP
  ports:
    - name: http-80
      port: 80
      targetPort: 80
      protocol: TCP
    - name: http-9000
      port: 9000
      targetPort: 9000
      protocol: TCP
---
 
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-proxy
  namespace: kube-proxy
spec:
  selector:
    matchLabels:
      app: kube-proxy
  replicas: 2
  template:
    metadata:
      labels:
        app: kube-proxy
    spec:
      containers:
        - name: kube-proxy
          readinessProbe:
            httpGet:
              port: 9000
              path: /health
              scheme: HTTP
            initialDelaySeconds: 5
            timeoutSeconds: 1
            periodSeconds: 5
            successThreshold: 1
            failureThreshold: 3
          livenessProbe:
            httpGet:
              port: 9000
              path: /health
              scheme: HTTP
            initialDelaySeconds: 5
            timeoutSeconds: 1
            periodSeconds: 5
            successThreshold: 1
            failureThreshold: 3
          image: bryantrh/kube-tunnel-gateway:v1.0
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
              protocol: TCP
            - containerPort: 9000
              protocol: TCP
          env:
            - name: KUBE_PROXY__LogConfig_Level
              value: "INFO"
            - name: KUBE_PROXY__Opt_AllowNamespaces
              value: "*"
            - name: KUBE_PROXY__Opt_BearerToken
              value: "eyJhbGciOiJSUzI1NiIsImtpZCI6Ik9NOEtCclFxdW53IiwidHlwIjoiSldUIn0.eyJhdWQiOlsia3ViZS10dW5uZWwiXSwiaWF0IjoxNjYyOTcwNDA5LCJpc3MiOiJrdWJlLXR1bm5lbCIsImp0aSI6IjA1YzM5NDIxLWUzMDYtNGIxMC1hMGZmLWE3NWZmZTA5ZGUyOSIsInN1YiI6Imt1YmUtcHJveHkifQ.yZkeTGD0LVNtpGz7_PZoSjma9d7-HGcQodDmQeJoakbeViYYvseqD5-lHwuVjtgPzonNv1EV9RneK5vd10ptqlgVCk8E2yw1lVoE2vCyu34T6YbuDAoN24bvKTlPQvcv_NEcQc41zjjyaj_JrXX-UOtQ9H0gpIzJwanoqOgst_yu6lAQAYyTkKzrbDw2neEOCowYz9BKbCWo3IZIU_U_YNpFhLjubeZrDXULRHym7vVfoqRyHyMu4t8einD-3plfA-rAFIoj-duZst0JaOfTKfWddVKVCpWm-boZ8xCQMVa_-WrRbTVIeTPxCnTqm1VdSpyfKnxiiEuyh2f2jBw"
            - name: KUBE_PROXY__Opt_Cluster
              value: "test_vm-test"
            - name: KUBE_PROXY__Opt_GatewayAddress
              value: "https://kube-gateway.example.com"
            - name: KUBE_PROXY__Opt_RetryInterval
              value: "5s"
            - name: KUBE_PROXY__Opt_Secure
              value: "true"
            - name: KUBE_PROXY__Opt_DataTunnelCount
              value: "50"
      restartPolicy: Always
```

#### 5. 子集群创建ServiceAccount
这里为了方便演示，我们直接创建admin ClusterRole，后续使用可根据实际情况，创建sa ，管控不同ns的权限
```BASH
#clusterRoleBinding.yaml
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kube-proxy
  labels:
    app: kube-proxy
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: rbac.authorization.k8s.io
subjects:
- kind: ServiceAccount
  name: kube-proxy
  namespace: kube-system
 
#serviceAccount.yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube-proxy
  namespace: kube-system
  labels:
    app: kube-proxy
```

创建完成后，获取对应token：
```Bash
kubectl get secret ${secret_name} -n {namespace} -o jsonpath={".data.token"} | base64 -d
```

#### 6. 根据上面获取到的token，配置sa
```Bash
#创建
cm create sa -s vmtest-token -p test -c vm-test -n infra,infra--staging -t eyJhbGciOiJSUzI1NiIsImtpZCI6ImVCRGNaanExanNmZVl3enE1MmxzUEhXTUpDamgtTjJDYkwyNmU1SU82Z28ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlLXByb3h5LXRva2VuLW5yaHpyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6Imt1YmUtcHJveHkiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJjZGY3NDJkMC03YTIzLTRjN2ItYWFhOS00YmY4ODNmODA2NDUiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06a3ViZS1wcm94eSJ9.F5pzq2AtqlxWevqOjNET8Et652G-_3uxZCTXcSnlm_T44-vbAalH0SO0NhebsNmQgxjgZMj0sbWj9sWxV12HsU3FfrL6CAjn6OL69njAiv8YfDpUHh3Ptke9GmuAO5OLsgjD-ktGLP_H7UQJsWA8l8cQJ8J8cSCQkNP8E3b9FFQIseYtlM1DTj1U4LqgHcTZE4XQy0WEctNXULfI9JvTw1xZ_8LQnvv7b5NO_VteMoIOEqkjrkqjs3MqpvNgfKQS02Yx2hhyxJU2imLJ2AiCm46DS0fQcsHjJN112bdvxzqaDdZ3wVkf4p_lCTZzkApcMsEyhO34654JrnkChM7xzQ
 
#查看
cm get sa -p test -c vm-test
```

#### 7. 可生成一个kubeconfig文件，即可使用kubectl --kubeconfig 管理子集群
```Bash
cm kubeconfig -p test -c vm-test -n infra
 
生成kubeconfig: [test_vm-test.config] 成功
可通过 kubectl --kubeconfig test_vm-test.config  访问k8s集群资源

```

# 4. Demo
这里创建一个演示项目：
只包含一个main.go
```
package main
 
import (
    "net/http"
 
    "github.com/gin-gonic/gin"
)
 
func setupRouter() *gin.Engine {
 
    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, "ok")
    })
    // Ping test
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": 0,
            "msg":    "success",
            "data": gin.H{
                "content": "pong",
            },
        })
    })
    return r
}
 
func main() {
    r := setupRouter()
    r.Run(":80")
}
```
1. 在当前目录下执行如下命令，进行初始化
```
cm hx init
```
根据提示输入项目名称、项目描述、项目所属部门(通常为gitlab组名)，项目版本号

即可生成helmx.yml、dockerfile.default、Makefile.default，如图所示：

目前可支持python，java，golang，生成模板Dockerfile

java 是根据当前目录是否存在pom.xml进行判断

python 是根据当前目录是否存在requirements.txt进行判断

golang 是根据当前目录是否存在go.mod进行判断

然后我们可根据项目的实际情况，修改Dockerfile.default、Makefile.default、helmx.yml 等文件。

其中helmx.yml 便是用来生成最终k8s部署文件的模板文件，通常需要以下内容即可：

```
#helmx.yml
 
service:
  ports:
    - "80"
  readinessProbe:
    action: http://:80
    initialDelaySeconds: 5
    periodSeconds: 5
  livenessProbe:
    action: http://:80
    initialDelaySeconds: 5
    periodSeconds: 5
resources:
  cpu: 100/1000m
  memory: 200/1000M
```
包含服务端口，健康检查，资源限制。

2. 如果项目包含配置文件，在项目根目录创建config 目录，根据分支创建配置文件，最终apply时候，会自动根据分支注入至k8s部署文件中，
   
3. 本地可以通过如下命令，查看渲染出来的k8s 部署文件
```
cm hx apply --dry-run
```

4. 创建.gitlab-ci.yml 文件
这里只是演示，所以只包含了ship 和apply两步，即编译镜像，部署至k8s集群，两个步骤，实际项目中，我们可能包含代码检查，单元测试等步骤，我们可采用模板文件，其他项目通过include引用即可。

```
.base: &base
  only:
    - master
    - main
    - develop
    - /^feat(ure)?\/.*$/
    - /^release\/.*$/
    - /^test\/.*$/
 
stages:
#  - sast
#  - lint
  - ship
  - apply
 
ship:
  <<: *base
  stage: ship
  variables:
    MULTI_ARCH_BUILDER: 1
  image: ${BUILDX_IMAGE}
  before_script:
    - echo "${DOCKER_PASSWORD}" | docker login "${DOCKER_REGISTRY}" -u="${DOCKER_USERNAME}" --password-stdin
    #- BUILDKIT_NAME=golang-buildkit cm hx ci-setup
    - cm hx ci-setup
  script:
    - cm hx buildx --with-builder --push --platform linux/amd64
 
apply:
  <<: *base
  stage: apply
  script:
    #- cm hx apply --dry-run
    - cm hx apply -p ${PROJECT_NAME} -c ${CLUSTER_NAME}
  dependencies:
    - ship
```


# 5. CM使用
```
cm is a tool for managing k8s clusters.
You can invoke cm through comand: "cm [command]..."
 
Usage:
  cm [command]
 
Available Commands:
  add         添加一个或更多 resources
  create      创建一个或更多 resources
  delete      删除一个或更多 resources
  get         显示一个或更多 resources
  help        Help about any command
  hx          用于初始化项目，部署服务
  kubeconfig  生成kubeConfig文件
  proxy       生成kube_proxy部署文件
  update      更新一个或更多 resources
  version     version for cm
 
Flags:
  -h, --help      help for cm
  -v, --v Level   number for the log level verbosity
 
Use "cm [command] --help" for more information about a command.
```

使用命令行，可以对项目以及集群进行增删查改等操作，hx子命令主要是用于CICD过程中，通过-h 即可查看每个子命令的用法