
# 1 部署

```bash
kubectl create ns argo
kubectl apply -f quick-start-mysql.yaml -n argo
```

获取token

```bash
# 在argo server的pod中获取token
kubectl exec -it argo-server-7b59b5f8cf-p8kkf -n argo -- argo auth token
```

部署ArgoCli

```bash
# Download the binary
curl -sLO https://github.com/argoproj/argo-workflows/releases/download/v3.3.8/argo-linux-amd64.gz

# Unzip
gunzip argo-linux-amd64.gz

# Make binary executable
chmod +x argo-linux-amd64

# Move binary to path
mv ./argo-linux-amd64 /usr/bin/argo

# Test installation
argo version
```

部署argo event

```bash
# create ns
kubectl create namespace argo-events

# 1. Deploy Argo Events SA, ClusterRoles, and Controller for Sensor, EventBus, and EventSource.
kubectl apply -f https://raw.githubusercontent.com/argoproj/argo-events/stable/manifests/install.yaml
# 安装认证admission控制器
kubectl apply -f https://raw.githubusercontent.com/argoproj/argo-events/stable/manifests/install-validating-webhook.yaml

 * On GKE, you may need to grant your account the ability to create new custom resource definitions and clusterroles        kubectl create clusterrolebinding YOURNAME-cluster-admin-binding --clusterrole=cluster-admin --user=YOUREMAIL@gmail.com * On OpenShift:
     - Make sure to grant `anyuid` scc to the service account.        oc adm policy add-scc-to-user anyuid system:serviceaccount:argo-events:default     - Add update permissions for the `deployments/finalizers` and `clusterroles/finalizers` of the argo-events-webhook ClusterRole(this is necessary for the validating admission controller)        - apiGroups:
          - rbac.authorization.k8s.io
          resources:
          - clusterroles/finalizers
          verbs:
          - update
        - apiGroups:
          - apps
          resources:
          - deployments/finalizers
          verbs:
          - update
          
# deploy
kubectl apply -n argo-events -f https://raw.githubusercontent.com/argoproj/argo-events/stable/examples/eventbus/native.yaml
```

# 2 Argo workflows demo

## 2.1 Argo Hello World

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: hello-world-
  namespace: argo
  labels:
    workflows.argoproj.io/archive-strategy: "false"
spec:
  entrypoint: rsqsay
  templates:
  - name: rsqsay
    container: 
      image: harbor.rsq.cn/busybox/busybox:1.33.1
      command: [echo]
      args: ["hello world rsq"]
  - name: gen-random-int
    container:
      image: harbor.rsq.cn/busybox/busybox:1.33.1
      command: [sh]
      source: |
        cat /dev/urandom | od -N2 -An -i | awk -v f=1 -v r=100 '{printf "%i\n", f + r * $1 / 65536}'
```

## 2.2 Argo Artifactory

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: artifact-passing-
spec:
  entrypoint: artifact-example
  templates:
  - name: artifact-example
    steps:
    - - name: generate-artifact
        template: whalesay
    - - name: consume-artifact
        template: print-message
        arguments:
          artifacts:
          - name: message
            from: "{{steps.generate-artifact.outputs.artifacts.hello-art}}"

  - name: whalesay
    container:
      image: docker/whalesay:latest
      command: [sh, -c]
      args: ["sleep 1; cowsay hello world | tee /tmp/hello_world.txt"]
    outputs:
      artifacts:
      - name: hello-art
        path: /tmp/hello_world.txt

  - name: print-message
    inputs:
      artifacts:
      - name: message
        path: /tmp/message
    container:
      image: alpine:latest
      command: [sh, -c]
      args: ["cat /tmp/message"]
```

## 2.3 Argo CICD

template workflow: 
```yaml
apiVersion: argoproj.io/v1alpha1
kind: WorkflowTemplate
metadata:
  annotations:
    workflows.argoproj.io/description: |
      Checkout out from Git, build and deploy application.
    workflows.argoproj.io/maintainer: '@rsq'
    workflows.argoproj.io/tags: go, git
    workflows.argoproj.io/version: '>= 2.9.0'
  name: devops-go 
spec:
  entrypoint: main
  arguments:
    parameters:
      # CI 代码仓库地址
      - name: repo
        value: 172.16.0.32/devops/goproj.git
      # CI 代码仓库branch
      - name: branch
        value: main
      # CI要制作的镜像名
      - name: image
        value: harbor.rsq.cn/go/goproj:20220728
      - name: cache-image
        value: harbor.rsq.cn/go/golang:1.18.4-centos7.9
      # CI 指定git仓库中Dockerfile的名字
      - name: dockerfile
        value: Dockerfile
      # CD 的时候指定的仓库
      - name: devops-cd-repo
        value: gitlab-test.coolops.cn/root/devops-cd.git
      - name: gitlabUsername
        value: devops_001
      - name: gitlabPassword
        value: 123456
  templates:
    - name: main
      steps:
      # 在step类型中，[--]代表顺序执行; [-]代表并行执行
      - - name: Checkout
          template: Checkout
      - - name: Build
          template: Build
      - - name: BuildImage
          template: BuildImage
      #- - name: Deploy
      #    template: Deploy
    # 拉取代码
    - name: Checkout
      script:
        image: harbor.rsq.cn/go/golang:1.18.4-centos7.9
        workingDir: /work
        command:
        - sh
        source: |
          git clone --branch {{workflow.parameters.branch}} http://{{workflow.parameters.gitlabUsername}}:{{workflow.parameters.gitlabPassword}}@{{workflow.parameters.repo}} .
        volumeMounts:
        - mountPath: /work
          name: work
    # 编译打包  
    - name: Build
      script:
        image: harbor.rsq.cn/go/golang:1.18.4-centos7.9
        workingDir: /work
        command:
        - sh
        source: go build -o goproj 
        volumeMounts:
        - mountPath: /work
          name: work
    # 构建镜像  
    - name: BuildImage
      container:
        image: harbor.rsq.cn/go/kaniko-executor:v1.5.0
        workingDir: /work
        command:
          - executor
        args:
          - --context=.
          - --dockerfile={{workflow.parameters.dockerfile}}
          - --destination={{workflow.parameters.image}}
          - --skip-tls-verify
          - --reproducible
          - --cache=true
          - --cache-repo={{workflow.parameters.cache-image}}
        volumeMounts:
        - mountPath: /work
          name: work
        - name: docker-config
          mountPath: /kaniko/.docker/
      volumes:
      - name: docker-config
        secret:
          secretName: docker-config
    # 部署  
    #- name: Deploy
    #  script:
    #    image: registry.cn-hangzhou.aliyuncs.com/rookieops/kustomize:v3.8.1
    #    workingDir: /work
    #    command:
    #    - sh
    #    source: |
    #       git remote set-url origin http://{{workflow.parameters.gitlabUsername}}:{{workflow.parameters.gitlabPassword}}@{{workflow.parameters.devops-cd-repo}}
    #       git config --global user.name "Administrator"
    #       git config --global user.email "songqi.ren@rsq.com"
    #       git clone http://{{workflow.parameters.gitlabUsername}}:{{workflow.parameters.gitlabPassword}}@{{workflow.parameters.devops-cd-repo}} /work/devops-cd
    #       cd /work/devops-cd
    #       git pull
    #       cd /work/devops-cd/devops-simple-go
    #       kustomize edit set image {{workflow.parameters.image}}
    #       git commit -am 'image update'
    #       git push origin master
    #    volumeMounts:
    #      - mountPath: /work
    #        name: work
  volumeClaimTemplates:
  - name: work
    metadata:
      name: work
    spec:
      storageClassName: nfs-client
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi 
```

创建docker-config的secret，这个secret要挂载到BuildImage的容器中，为build镜像准备

```bash
# 1 auth生成base64加密，此用户密码为docker镜像仓库的账号密码
echo -n "<user>:<password>" | base64

# 2 生成config.json
mkdir ${HOME}/.docker && cd ${HOME}/.docker
vim config.json
{
    "auths": {
         "http://harbor.rsq.cn": {
            "auth": "dXNlcjpwYXNzd29yZAo="
         }
    }
}

# 3 创建secret
kubectl create secret generic docker-config --from-file=/root/.docker/config.json -n argo
```

创建template，或者直接在web界面点击Workflow Templates -->CREATE NEW TEMPLATE WORKFLOW

```bash
argo template create -n argo cicd.yaml
```


使用如下命令创建并观察workflow:

```bash
$ argo submit -n argo helloworld.yaml --watch
```

还可以通过argo list来查看状态，如下：
```bash
# argo list -n argo
NAME                STATUS      AGE   DURATION   PRIORITY
hello-world-9pw7v   Succeeded   1m    10s        0
使用argo logs来查看具体的日志，如下：
# argo logs -n argo hello-world-9pw7v
hello-world-9pw7v:  _____________
hello-world-9pw7v: < hello world >
hello-world-9pw7v:  -------------
hello-world-9pw7v:     \
hello-world-9pw7v:      \
hello-world-9pw7v:       \     
hello-world-9pw7v:                     ##        .            
hello-world-9pw7v:               ## ## ##       ==            
hello-world-9pw7v:            ## ## ## ##      ===            
hello-world-9pw7v:        /""""""""""""""""___/ ===        
hello-world-9pw7v:   ~~~ {~~ ~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~   
hello-world-9pw7v:        \______ o          __/            
hello-world-9pw7v:         \    \        __/             
hello-world-9pw7v:           \____\______/   
```

# 参考文章
[Argo Workflows部署](https://blog.csdn.net/qq_29062169/article/details/125486646)
[Kubernetes 原生 CI/CD 构建框架 Argo](https://www.isolves.com/it/cxkf/kj/2021-11-30/46743.html)
[Argo Workflows —— Kubernetes的工作流引擎入门](https://blog.csdn.net/a772304419/article/details/125463627)
[给Workflow添加参数](https://www.jianshu.com/p/f7276c61f072)