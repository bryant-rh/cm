# CM
这是一个用于CICD 流程中，初始化项目，生成k8s部署文件，且管理k8s集群通道，进行部署的工具

# Installation

## 1.部署cm-server
cm-server 是用于管理集群通道Token的服务端，可以创建项目，集群关联到项目等

```Bash
cd deploy/cm-server

#部署mysql
kubectl apply -f mysql5.7-deploy.yaml

#部署cm-server
kubectl apply -f cm-server-deploy.yaml

```

## 2. 部署kube-tunnel-gateway
kube-tunnel-gateway 是用于创建和需要管控k8s集群的服务端，用于建立通道

### 1. 生成PrivateKey
```Bash
openssl genrsa -out private.key
echo $(cat private.key|base64)|base64
```

### 2. 生成PublicKey
```Bash
openssl rsa -in private.key -pubout -outform PEM -out public.pem
cat public.pem|base64
```

### 3. 用上面生成的PrivateKey 和PublicKey 替换部署文件中 Secret 的对应字段，部署
[备注] : 必须 https 访问

```Bash
cd deploy/kube-tunnel
kubectl apply -f kube-tunnel-gateway.yaml
```

