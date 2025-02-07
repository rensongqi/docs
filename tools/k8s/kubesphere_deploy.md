KubeSphere部署

机器列表：

|主机IP|主机名|备注|
|---|---|---|
|172.16.0.50|master||
|172.16.0.51|node01||
|172.16.0.52|node02|nfs server|
|172.16.0.53|node03||
|172.16.0.58|node08||
|172.16.0.59|node09||
|172.16.0.60|node10||
|172.16.0.60|node11||


## 1 部署k8s集群
这里我用的ansible来批量部署初始化集群，ansible配置在 芜湖R750 那台机器
系统的初始化包括防火墙、hosts文件都已经集成在ansible中

### 1.1 修改hosts文件
```bash
# vim /etc/ansible/hosts
[master]
172.16.0.50 node_name=master

[node]
172.16.0.51 node_name=node01
172.16.0.52 node_name=node02
172.16.0.53 node_name=node03
172.16.0.58 node_name=node08
172.16.0.59 node_name=node09
172.16.0.60 node_name=node10
172.16.0.61 node_name=node11

[nfs]
172.16.0.51 node_name=node01
172.16.0.53 node_name=node03
172.16.0.58 node_name=node08
172.16.0.59 node_name=node09
172.16.0.60 node_name=node10
172.16.0.61 node_name=node11
```

### 1.2 修改group_vars
安装k8s版本和网络插件自定义即可
```bash
# vim /etc/ansible/group_vars/all.yml
# k8s version [1.18.3/1.19.4/...]
# 支持大多数版本
k8s_version: '1.22.10'

# network plugin [flannel/calico]
network_plugin: 'calico'

# 集群网络规划
service_cidr: '10.96.0.0/12'
pod_cidr: '10.244.0.0/16'

# 集群IP规划，master暂时只支持一个，node可以有多个
hosts:
  master:
    - '172.16.0.50'
  node:
    - '172.16.0.51'
    - '172.16.0.52'
    - '172.16.0.53'
    - '172.16.0.58'
    - '172.16.0.59'
    - '172.16.0.60'
    - '172.16.0.61'
```

### 1.3 部署k8s集群
以root用户执行如下命令
```
# cd /etc/ansible/
ansible-playbook k8s_single_master_install.yml
```
## 2 配置nfs和storageClass
### 2.1 nfs server配置如下
```bash
# 1. 在每个机器内执行
yum install -y nfs-utils
# 用ansible操作命令如下
ansible all -m yum -a "name=nfs-utils state=latest"

# 2. 在nfs server上执行以下命令 
echo "/home/data/ *(insecure,rw,sync,no_root_squash)" > /etc/exports

# 3. 在nfs server执行
systemctl enable rpcbind
systemctl enable nfs-server
systemctl start rpcbind
systemctl start nfs-server

# 4. 在nfs server执行 使配置生效
exportfs -r

# 5. 在nfs server执行 检查配置是否生效
exportfs
```

### 2.2 nfs client配置如下
```bash
# 1. 没台client机器均需要创建挂在的目录
ansible all -m shell -a "mkdir /data"
# 2. 用ansible对所有的node节点进行nfs 挂载
ansible nfs -m shell -a "echo '172.16.0.52:/home/data /data nfs defaults 0 0' >>/etc/fstab"
ansible nfs -m shell -a "mount -a"
```

### 2.3 创建nfs storageclass
```bash
kubectl apply -f nfs.yaml
```
yaml文件如下:
```yaml
## 创建了一个存储类
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: nfs-storage
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: k8s-sigs.io/nfs-subdir-external-provisioner
parameters:
  archiveOnDelete: "true"  ## 删除pv的时候，pv的内容是否要备份

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nfs-client-provisioner
  labels:
    app: nfs-client-provisioner
  # replace with namespace where provisioner is deployed
  namespace: default
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: nfs-client-provisioner
  template:
    metadata:
      labels:
        app: nfs-client-provisioner
    spec:
      serviceAccountName: nfs-client-provisioner
      containers:
        - name: nfs-client-provisioner
          image: registry.cn-hangzhou.aliyuncs.com/lfy_k8s_images/nfs-subdir-external-provisioner:v4.0.2
          # resources:
          #    limits:
          #      cpu: 10m
          #    requests:
          #      cpu: 10m
          volumeMounts:
            - name: nfs-client-root
              mountPath: /persistentvolumes
          env:
            - name: PROVISIONER_NAME
              value: k8s-sigs.io/nfs-subdir-external-provisioner
            - name: NFS_SERVER
              value: 172.16.0.52 ## 指定自己nfs服务器地址
            - name: NFS_PATH  
              value: /home/data  ## nfs服务器共享的目录
      volumes:
        - name: nfs-client-root
          nfs:
            server: 172.16.0.52
            path: /home/data
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-client-provisioner
  # replace with namespace where provisioner is deployed
  namespace: default
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: nfs-client-provisioner-runner
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "update", "patch"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: run-nfs-client-provisioner
subjects:
  - kind: ServiceAccount
    name: nfs-client-provisioner
    # replace with namespace where provisioner is deployed
    namespace: default
roleRef:
  kind: ClusterRole
  name: nfs-client-provisioner-runner
  apiGroup: rbac.authorization.k8s.io
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-client-provisioner
  # replace with namespace where provisioner is deployed
  namespace: default
rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-client-provisioner
  # replace with namespace where provisioner is deployed
  namespace: default
subjects:
  - kind: ServiceAccount
    name: nfs-client-provisioner
    # replace with namespace where provisioner is deployed
    namespace: default
roleRef:
  kind: Role
  name: leader-locking-nfs-client-provisioner
  apiGroup: rbac.authorization.k8s.io
```

## 3 部署metrics-server
```bash
kubectl apply -f metrics-server.yaml
```
yaml文件如下:
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    k8s-app: metrics-server
  name: metrics-server
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    k8s-app: metrics-server
    rbac.authorization.k8s.io/aggregate-to-admin: "true"
    rbac.authorization.k8s.io/aggregate-to-edit: "true"
    rbac.authorization.k8s.io/aggregate-to-view: "true"
  name: system:aggregated-metrics-reader
rules:
- apiGroups:
  - metrics.k8s.io
  resources:
  - pods
  - nodes
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    k8s-app: metrics-server
  name: system:metrics-server
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - nodes
  - nodes/stats
  - namespaces
  - configmaps
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  labels:
    k8s-app: metrics-server
  name: metrics-server-auth-reader
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: extension-apiserver-authentication-reader
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    k8s-app: metrics-server
  name: metrics-server:system:auth-delegator
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    k8s-app: metrics-server
  name: system:metrics-server
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:metrics-server
subjects:
- kind: ServiceAccount
  name: metrics-server
  namespace: kube-system
---
apiVersion: v1
kind: Service
metadata:
  labels:
    k8s-app: metrics-server
  name: metrics-server
  namespace: kube-system
spec:
  ports:
  - name: https
    port: 443
    protocol: TCP
    targetPort: https
  selector:
    k8s-app: metrics-server
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    k8s-app: metrics-server
  name: metrics-server
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: metrics-server
  strategy:
    rollingUpdate:
      maxUnavailable: 0
  template:
    metadata:
      labels:
        k8s-app: metrics-server
    spec:
      containers:
      - args:
        - --cert-dir=/tmp
        - --kubelet-insecure-tls
        - --secure-port=4443
        - --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname
        - --kubelet-use-node-status-port
        image: registry.cn-hangzhou.aliyuncs.com/lfy_k8s_images/metrics-server:v0.4.3
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /livez
            port: https
            scheme: HTTPS
          periodSeconds: 10
        name: metrics-server
        ports:
        - containerPort: 4443
          name: https
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /readyz
            port: https
            scheme: HTTPS
          periodSeconds: 10
        securityContext:
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          runAsUser: 1000
        volumeMounts:
        - mountPath: /tmp
          name: tmp-dir
      nodeSelector:
        kubernetes.io/os: linux
      priorityClassName: system-cluster-critical
      serviceAccountName: metrics-server
      volumes:
      - emptyDir: {}
        name: tmp-dir
---
apiVersion: apiregistration.k8s.io/v1
kind: APIService
metadata:
  labels:
    k8s-app: metrics-server
  name: v1beta1.metrics.k8s.io
spec:
  group: metrics.k8s.io
  groupPriorityMinimum: 100
  insecureSkipTLSVerify: true
  service:
    name: metrics-server
    namespace: kube-system
  version: v1beta1
  versionPriority: 100
```

## 4 部署Kubesphere
1. 配置文件如下：
- [kubesphere-installer.yaml](https://lpq-k8s.oss-cn-hangzhou.aliyuncs.com/kubesphere-installer.yaml)
- [cluster-configuration.yaml](https://lpq-k8s.oss-cn-hangzhou.aliyuncs.com/cluster-configuration.yaml)

2. 修改cluster-configuration.yaml

将文件内172.31.0.4 全局替换为自己master节点的IP地址

3. 部署KubeSphere
```bash
kubectl apply -f kubesphere-installer.yaml
kubectl apply -f cluster-configuration.yaml
```

4. 查看安装进度
```bash
kubectl logs -n kubesphere-system $(kubectl get pod -n kubesphere-system -l app=ks-install -o jsonpath='{.items[0].metadata.name}') -f
```

5. 如果有etcd监控证书找不到的bug，执行如下命令
```bash
kubectl -n kubesphere-monitoring-system create secret generic kube-etcd-client-certs  --from-file=etcd-client-ca.crt=/etc/kubernetes/pki/etcd/ca.crt  --from-file=etcd-client.crt=/etc/kubernetes/pki/apiserver-etcd-client.crt  --from-file=etcd-client.key=/etc/kubernetes/pki/apiserver-etcd-client.key
```

6. 部署成功后log输出如下
```bash
Collecting installation results ...
#####################################################
###              Welcome to KubeSphere!           ###
#####################################################

Console: http://172.16.0.50:30880
Account: admin
Password: P@88w0rd

NOTES：
  1. After you log into the console, please check the
     monitoring status of service components in
     "Cluster Management". If any service is not
     ready, please wait patiently until all components 
     are up and running.
  2. Please change the default password after login.

#####################################################
https://kubesphere.io             2022-05-30 13:50:14
#####################################################
```

参考博客：[k8s与kubesphere安装](https://blog.csdn.net/baidu_41860619/article/details/124920503)
