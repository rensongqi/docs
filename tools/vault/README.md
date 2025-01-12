

# Vault部署
参考文章：
[vault-in-kubernetes](https://devopscube.com/vault-in-kubernetes/)

## k8s deploy
```bash
# git clone
git clone https://github.com/scriptcamp/kubernetes-vault.git


# 1 创建rbac
cd vault-manifests
kubectl create -f rbac.yaml

# 2 配置vault storage 
# 参考文章：https://medium.com/@holden.omans/use-minio-or-s3-as-a-vault-backend-on-k8s-335312a6349a
# 2.1 minio oss
# 参考如下配置
    storage "s3" {
      endpoint    = "ossapi.rsq.cn:9000"
      access_key  = "vault"
      secret_key  = "dcwa8goj0e6EiGsxxxxxxxxx"
      bucket      = "vault"
      region      = "oss-stg"
      disable_ssl = "true"
      s3_force_path_style = "true"
    }
# 2.2 local storage
# 参考如下配置
    storage "file" {
      path = "/vault/data"
    }
kubecmltl apply -f configmap.yaml
 
# 3 部署service
kubectl apply -f services.yaml
 
# 4 部署statefulset
kubectl apply -f statefulset.yaml
```

## 初始化vault
```bash
# Generate Token
kubectl exec vault-0 -n default  -- vault operator init -key-shares=1 -key-threshold=1 -format=json > keys.json

# 获取root key
VAULT_UNSEAL_KEY=$(cat keys.json | jq -r ".unseal_keys_b64[]")
echo $VAULT_UNSEAL_KEY

VAULT_ROOT_KEY=$(cat keys.json | jq -r ".root_token")
echo $VAULT_ROOT_KEY

# Unseal is the state at which the vault can construct keys that are required to decrypt the data stored inside it.
kubectl exec vault-0 -n default -- vault operator unseal $VAULT_UNSEAL_KEY
```

## 日志审计
```bash
# 对某个path开启审计
vault audit enable -path=attendance file file_path=/vault/logs/audit1.log

# 查看所有审计日志
/ $ vault audit list
Path           Type    Description
----           ----    -----------
attendance/    file    n/a
file/          file    n/a

# 关闭某一个path的日志审计
vault audit disable attendance
```

# AppRole

1. 开启approle
```bash
vault auth enable approle
```

2. 创建kv对
```bash
vault kv put -mount=devops jenkins username=admin
vault kv put -mount=devops jenkins password=123456
vault kv get devops/jenkins/gitlab
```

3. 创建policy
jenkins-policy.hcl
```
# path 为 vault kv get devops/jenkins/gitlab 获取的 Secret Path
path "devops/data/jenkins/*" {
    capabilities= ["read"]
}
```
create
```bash
vault policy write jenkins-policy ./jenkins-policy.hcl
```

4. 创建role绑定policy
```bash
vault write auth/approle/role/jenkins \
    token_type=batch \
    secret_id_ttl=60m \
    token_ttl=60m \
    token_max_ttl=60m \
    secret_id_num_uses=40 \
    token_policies=jenkins-policy
```

5. 创建secret-id
```bash
vault write -f auth/approle/role/jenkins/secret-id
```

6. read role-id
```bash
vault read auth/approle/role/jenkins/role-id
```
