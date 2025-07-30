
# docker-compose 部署

```yaml
version: '3.7'

x-minio-common: &minio-common
  image: "harbor.rsq.cn/library/minio/minio:RELEASE.2023-10-24T04-42-36Z"
  command: server --console-address ":9001" http://minio{1...3}/data0 http://minio{1...3}/data1  http://minio{1...3}/data2  http://minio{1...3}/data3  http://minio{1...3}/data4  http://minio{1...3}/data5
  expose:
    - "9000"
    - "9001"
# 增加host映射，以便三个节点之间通过域名连通
  extra_hosts:
    minio1: 172.17.1.17
    minio2: 172.17.1.18
    minio3: 172.17.1.19
  depends_on:
    - nginx
  environment:
    MINIO_ROOT_USER: admin
    MINIO_ROOT_PASSWORD: xxxxxxxx
    MINIO_DOMAIN: ossapi.rsq.cn
    #MINIO_STORAGE_CLASS_STANDARD: "EC:3"
    #MINIO_STORAGE_CLASS_RRS: "EC:2"
    TZ: Asia/Shanghai
  healthcheck:
    test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
    interval: 30s
    timeout: 20s
    retries: 5

# 数据盘挂载目录，按需修改
services:
  minio1:
    <<: *minio-common
    container_name: minio1
    hostname: minio1
    volumes:
      - /minio0:/data0
      - /minio1:/data1
      - /minio2:/data2
      - /minio3:/data3
      - /minio4:/data4
      - /minio5:/data5
    ports:
      - "9000:9000"
      - "9001:9001"
    restart: always
    network_mode: host

  minio2:
    <<: *minio-common
    container_name: minio2
    hostname: minio2
    volumes:
      - /minio0:/data0
      - /minio1:/data1
      - /minio2:/data2
      - /minio3:/data3
      - /minio4:/data4
      - /minio5:/data5
    ports:
      - "9000:9000"
      - "9001:9001"

  minio3:
    <<: *minio-common
    container_name: minio3
    hostname: minio3
    volumes:
      - /minio0:/data0
      - /minio1:/data1
      - /minio2:/data2
      - /minio3:/data3
      - /minio4:/data4
      - /minio5:/data5
    ports:
      - "9000:9000"
      - "9001:9001"
# 三个节点都安装nginx，并且负载到三个节点，nginx暴露端口按需修改19000/19001
  nginx:
    image: "harbor.rsq.cn/library/nginx:1.19.2-alpine"
    container_name: nginx-minio
    hostname: nginx
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    ports:
      - "30009:9000"
      - "30099:9001"
    extra_hosts:
      minio1: 172.16.100.107
      minio2: 172.16.100.108
      minio3: 172.16.100.109
```

# mc 命令
> [mc下载](https://github.com/minio/mc/releases)

```bash
# 设置别名
mc alias set local-oss http://ossapi.rsq.cn:9000 <secretKey> <secretPass>

# 拷贝本地文件至oss
mc cp --recursive /PATH/TO/FILE_OR_DIR local-oss/<bucket>/<path>

# 拷贝oss数据至本地
mc cp --recursive local-oss/<bucket>/<path> /PATH/TO/DIR

# 跟踪bucket
mc admin trace -v --path <bucket>/* local-oss

# 冷热分层取回数据
mc ilm restore --days 9999 local-oss/<bucket>/<file_path>

# 删除Tier
mc ilm tier rm local-oss AI-RSQ-PREDICTIONELRTD-WARM

# 其它分层配置
mc ilm rule add --expire-days 90 --noncurrent-expire-days 30  local-oss/mydata

mc ilm rule add --expire-delete-marker local-oss/mydata

mc ilm rule add --transition-days 30 --transition-tier "COLDTIER" local-oss/mydata

mc ilm rule add --noncurrent-transition-days 7 --noncurrent-transition-tier "COLDTIER"
```
