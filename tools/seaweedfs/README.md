# 目录
- [目录](#目录)
- [1 规划](#1-规划)
- [2 配置](#2-配置)
  - [2.1 初始化配置](#21-初始化配置)
    - [2.1.1 配置Volume元数据存储](#211-配置volume元数据存储)
    - [2.1.2 配置S3权限](#212-配置s3权限)
    - [2.1.3 172.16.90.172](#213-1721690172)
    - [2.1.4 172.16.90.173](#214-1721690173)
    - [2.1.5 172.16.90.174](#215-1721690174)
  - [2.2 启动服务](#22-启动服务)
- [3 挂载](#3-挂载)

# 1 规划

计算需要创建多少volume，有两种方式，一种是自动分配卷数量，另一种是手动分配。

1、手动分配

> volumeSize * volumeNums >= TotalStorage（单节点）

以`172.16.90.172`为例，有两块磁盘，每块可用磁盘空间为212TB，volumeSize=30GB来计算，每块磁盘可以创建的卷数量为：`7236`，则创建Volume时需要指定卷最大数量为`-max=7236`
```bash
212 * 1024 / 30 ~= 7236
```

2、自动分配

Volume可以配置`-max=0`

主机列表

| 主机         | 部署模块  | 备注                               |
| ------------ | --------| ---------------------------------- |
| 172.16.90.172        | Master、Volume、Filer、S3、Redis    | 元数据存储在redis中     |
| 172.16.90.173       | Master      |                              |
| 172.16.90.174        | Master    |                    |

173、174为临时master，后续可以stop任意一个Master节点，然后其对应的metadata数据迁移至新的节点，可以达到更换Master的效果

# 2 配置
## 2.1 初始化配置
```bash
mkdir -p /data/seaweedfs/{config,metadata}
```

### 2.1.1 配置Volume元数据存储
```bash
vim /data/seaweedfs/config/filer.toml
cat >/data/seaweedfs/config/filer.toml<<EOF
[leveldb2]
# 元数据存储至本地
enabled = false
dir="/data/metadata"

[redis]
enabled = true
address  = "172.16.90.172:6379"
password = "xxxxxxxx"
EOF
```
### 2.1.2 配置S3权限
```bash
vim /data/seaweedfs/config/s3.json
cat >/data/seaweedfs/config/s3.json<<EOF
{
  "identities": [
    {
      "name": "admin",
      "credentials": [
        {
          "accessKey": "LPW1Rq7V2xxxxxxx",
          "secretKey": "LGy6grRx3V3bqqxxxxxxxxx"
        }
      ],
      "actions": ["Admin", "Read", "Write", "List"]
    },
    {
      "name": "user1",
      "credentials": [
        {
          "accessKey": "hygumFx9xxxxxxxxxx",
          "secretKey": "4waojUHo7Yw231Axxxxxxxxxxxxxxxxxxxx"
        }
      ],
      "actions": ["Read", "List"]
    }
  ]
}
EOF
```
### 2.1.3 172.16.90.172
```yaml
version: '3.8'
services:
  seaweed-master:
    image: harbor.rsq.cn/docker.io/chrislusf/seaweedfs:latest
    container_name: seaweedfs-master
    network_mode: host
    restart: always
    volumes:
      - /data/seaweedfs/metadata:/weed/metadata:rw
    entrypoint: /usr/bin/weed
    command: >
      -logtostderr=true
      master 
      -mdir=/weed/metadata
      -peers=172.16.90.172:9333,172.16.90.173:9333,172.16.90.174:9333
      -ip=172.16.90.172
      -port=9333
      -metricsPort=9999
      -volumePreallocate

  seaweed-volume:
    image: harbor.rsq.cn/docker.io/chrislusf/seaweedfs:latest
    container_name: seaweedfs-volume
    network_mode: host
    restart: always
    volumes:
      - /data/seaweedfs/metadata:/weed/metadata:rw
      - /disk/local_disk1:/disk/local_disk1:rw
      - /disk/local_disk2:/disk/local_disk2:rw
    entrypoint: /usr/bin/weed
    command: >
      -logtostderr=true
      volume 
      -mserver=172.16.90.172:9333,172.16.90.173:9333,172.16.90.174:9333
      -ip=172.16.90.172
      -port=8080
      -dir=/disk/local_disk1,/disk/local_disk2
      -dir.idx=/weed/metadata
      -max=14472
      -metricsPort=9998
      -readBufferSizeMB=64

  seaweed-filer:
    image: harbor.rsq.cn/docker.io/chrislusf/seaweedfs:latest
    container_name: seaweedfs-filer
    network_mode: host
    restart: always
    volumes:
      - /data/seaweedfs/metadata:/weed/metadata:rw
      - /data/seaweedfs/config:/etc/seaweedfs:rw
    entrypoint: /usr/bin/weed
    command: >
      -logtostderr=true
      filer
      -master=172.16.90.172:9333,172.16.90.173:9333,172.16.90.174:9333
      -ip=172.16.90.172
      -port=8888
      -metricsPort=9997

  seaweed-s3:
    image: harbor.rsq.cn/docker.io/chrislusf/seaweedfs:latest
    container_name: seaweedfs-s3
    network_mode: host
    restart: always
    entrypoint: /usr/bin/weed
    volumes:
      - /data/seaweedfs/config:/etc/seaweedfs
    command: >
      -logtostderr=true
      s3
      -filer=172.16.90.172:8888
      -ip.bind=0.0.0.0
      -metricsPort=9996
      -config=/etc/seaweedfs/s3.json
```
### 2.1.4 172.16.90.173
```yaml
version: '3.8'
services:
  seaweedfs-master:
    image: harbor.rsq.cn/docker.io/chrislusf/seaweedfs:latest
    container_name: seaweedfs-master
    network_mode: host
    restart: always
    volumes:
      - /data/seaweedfs/metadata:/weed/metadata:rw
    entrypoint: /usr/bin/weed
    command: >
      -logtostderr=true
      master 
      -mdir=/weed/metadata
      -peers=172.16.90.172:9333,172.16.90.173:9333,172.16.90.174:9333
      -ip=172.16.90.173
      -port=9333
      -metricsPort=9999
      -volumePreallocate
```
### 2.1.5 172.16.90.174
```yaml
version: '3.8'
services:
  seaweedfs-master:
    image: harbor.rsq.cn/docker.io/chrislusf/seaweedfs:latest
    container_name: seaweedfs-master
    network_mode: host
    restart: always
    volumes:
      - /data/seaweedfs/metadata:/weed/metadata:rw
    entrypoint: /usr/bin/weed
    command: >
      -logtostderr=true
      master 
      -mdir=/weed/metadata
      -peers=172.16.90.172:9333,172.16.90.173:9333,172.16.90.174:9333
      -ip=172.16.90.174
      -port=9333
      -metricsPort=9999
      -volumePreallocate
```
## 2.2 启动服务
```bash
# 先启动master
docker-compose up -d seaweedfs-master

# 然后启动volume
docker-compose up -d seaweedfs-volume

# 再启动filer
docker-compose up -d seaweedfs-filer

# 启动s3
docker-compose up -d seaweedfs-s3
```

# 3 挂载
```bash
weed -v 4 mount -filer=172.16.90.172:8888 -filer.path=/data -dir=/mnt/weed -cacheDir=/mnt -chunkSizeLimitMB=8 -volumeServerAccess=direct
```