# 1 创建Token

**devops_admin_token.yaml**

```yaml
# Create ServiceAccount
apiVersion: v1
kind: ServiceAccount
metadata:
  name: devops-admin
  namespace: default
---
# Create ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: devops-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: devops-admin
  namespace: default
```

kubectl apply -f devops_admin_token.yaml

# 2 获取Token

根据自己创建ClusterRoleBinding所选的Namespace去查找Secret，然后截取Token

```bash
TOKEN=$(kubectl describe secret `kubectl get secret -ndefault | grep devops-admin | awk '{print $1}'` | grep token: | awk '{print $NF}')
```

# 3 访问API

（1）查询pod日志

```bash
curl -H "Authorization: Bearer $TOKEN" https://192.168.1.1:6443/api/v1/namespaces/scr/pods/log-85bb949cd-m546t/log/ --insecure
```

（2）查询所有namespace

```bash
curl -H "Authorization: Bearer $TOKEN" https://192.168.1.1:6443/api/v1/namespaces/ --insecure
```

（3）查询一个namespace下所有的pod

```bash
curl -H "Authorization: Bearer $TOKEN" https://192.168.1.1:6443/api/v1/namespaces/scr/pods --insecure
```



**参考文章：** 

[（1）如何使用curl访问k8s的apiserver](https://developer.aliyun.com/article/706210)

[（2）K8s API概述](https://kubernetes.io/zh/docs/reference/using-api/#api-versioning)

