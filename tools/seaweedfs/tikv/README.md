

# 使用Tikv作为Filer 的store

> 由于redis内存型数据库，随着数据量的增长占用内存也越多，需要另一种kv数据库来支持

- [使用 TiUP 部署 TiDB 集群](https://docs.pingcap.com/zh/tidb/dev/production-deployment-using-tiup/)


如果网络受限的情况下推荐使用离线部署很方便，注意事项
1. Tikv集群之间需要进行ssh免密，主控机建议选择Tikv集群中某个节点即可

`
> 在[官方下载页面](https://cn.pingcap.com/product-community/?_gl=1*20ofaf*_gcl_au*ODU5MzQ5MTcwLjE3NDIyODE1MTk.*_ga*Nzc3NDg2Njg0LjE3NDIyODE1MTk.*_ga_3JVXJ41175*MTc0MjQ1NjU1Ny40LjEuMTc0MjQ1NjU2NC41My4wLjE1ODU1NDQ2Nzc.*_ga_CPG2VW1Y41*MTc0MjQ1NjYwMy42LjAuMTc0MjQ1NjYwMy4wLjAuMA..)选择对应版本的 TiDB server 离线镜像包（包含 TiUP 离线组件包）。需要同时下载 TiDB-community-server 软件包和 TiDB-community-toolkit 软件包。

如下操作均在主控机执行

```bash
# 执行以下命令安装 TiUP 组件
# local_install.sh 脚本会自动执行 tiup mirror set tidb-community-server-${version}-linux-amd64 命令将当前镜像地址设置为 tidb-community-server-${version}-linux-amd64。
version=v8.5.1
tar xzvf tidb-community-server-${version}-linux-amd64.tar.gz && \
sh tidb-community-server-${version}-linux-amd64/local_install.sh && \
source /home/<user>/.bash_profile
```
合并离线包
```bash
tar xf tidb-community-toolkit-${version}-linux-amd64.tar.gz
ls -ld tidb-community-server-${version}-linux-amd64 tidb-community-toolkit-${version}-linux-amd64
cd tidb-community-server-${version}-linux-amd64/
cp -rp keys ~/.tiup/
tiup mirror merge ../tidb-community-toolkit-${version}-linux-amd64
```

初始化集群配置文件
```bash
tiup cluster template > topology.yaml
```

如果仅安装Tikv可参考如下配置 `topology.yaml`
> Tiup工具会在配置文件中不同服务器上安装指定的服务应用
```yaml
global:
  user: "root"
  ssh_port: 22
  deploy_dir: "/data/tikv/deploy"
  data_dir: "/data/tikv/data"

pd_servers:
  - host: 172.16.9.175
  - host: 172.16.9.177
  - host: 172.16.9.178

tikv_servers:
  - host: 172.16.9.175
  - host: 172.16.9.177
  - host: 172.16.9.178

monitoring_servers:
  - host: 172.16.9.175

grafana_servers:
  - host: 172.16.9.177

alertmanager_servers:
  - host: 172.16.9.178
```

Tikv集群所有节点均安装numactl
```bash
apt install numactl -y
```

检查集群存在的潜在风险：
```bash
tiup cluster check ./topology.yaml --user root [-p] [-i /home/root/.ssh/gcp_rsa]
```

自动修复集群存在的潜在风险：
```bash
tiup cluster check ./topology.yaml --apply --user root [-p] [-i /home/root/.ssh/gcp_rsa]
```

部署 Tikv 集群：
```bash
tiup cluster deploy tidb-test v8.5.0 ./topology.yaml --user root [-p] [-i /home/root/.ssh/gcp_rsa]
```

查看 TiUP 管理的集群情况
```bash
tiup cluster list
```

执行如下命令检查 tidb-test 集群情况
```bash
tiup cluster display tidb-test
```

安全启动集群
```bash
tiup cluster start tidb-test --init
```

验证集群运行状态
```bash
tiup cluster display tidb-test
```

# SeaweedFS配置
filer.toml
```
[tikv]
enabled = true
pdaddrs = "172.16.9.175:2379,172.16.9.177:2379,172.16.9.178:2379"
deleterange_concurrency = 8
enable_1pc = false
ca_path=""
cert_path=""
key_path=""
verify_cn=""
```


# 迁移Redis数据至Tikv store

执行如下命令，将 172.16.9.176:8888 filer对应的数据同步至 /data/seaweedfs/config/filer.toml 对应的存储中
```bash
weed filer.meta.backup -config=/data/seaweedfs/config/filer.toml  -filer="172.16.9.176:8888" -restart
```

- `-restart`参数可以支持异步增量同步，前提是已经全量同步过一次数据，但是该参数并不能确保增量数据完全同步，所以为了尽量减少对现有存储的影响，应在老存储的Filer切到Tikv之后再执行上述命令进行全量元数据同步