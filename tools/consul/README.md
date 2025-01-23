# Usage

注意需要使用到 `ibm-spectrum-scale-csi-lt` csi, 先确保该csi已安装,如未安装请使用其它存储

## Deploy consul

```bash
kubectl apply -f consul-deploy.yaml
```

## Deploy tensuns

如果需要第三方集成平台可部署tensuns,按需配置

```bash
kubectl apply -f tensuns-deploy.yaml
```

