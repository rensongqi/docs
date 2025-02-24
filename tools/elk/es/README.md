
- [三节点ES集群(带密码校验)](#三节点es集群带密码校验)
- [初始化获取证书](#初始化获取证书)
- [Configuration ES node](#configuration-es-node)
  - [172.16.10.20](#172161020)
  - [172.16.10.21](#172161021)
  - [172.16.10.22](#172161022)
- [配置密码](#配置密码)
- [kibana配置](#kibana配置)

# 三节点ES集群(带密码校验)

| IP | 角色 | 
| ---- | ---- |
| 172.16.10.20 | master,data 
| 172.16.10.21  | master,data 
| 172.16.10.22 | master,data 

# 初始化获取证书

```bash
mkdir /data/elk/es/{config,data,plugins} -p
chmod 777 /data/elk/es/
```

先部署一个不开启安全验证的容器

`elasticsearch.yml`

```yml
cluster.name: soft-es
node.name: es-node1
node.roles: [ master, data ]
network.host: 0.0.0.0
network.publish_host: 172.16.10.20
http.port: 9200
http.cors.enabled: true
http.cors.allow-origin: "*"
ingest.geoip.downloader.enabled: false
discovery.seed_hosts: ["172.16.10.20:9300", "172.16.10.21:9300", "172.16.10.22:9300"]
cluster.initial_master_nodes: ["es-node1","es-node2","es-node3"]
xpack.security.enabled: false
xpack.security.transport.ssl.enabled: false
```

`docker-compose.yml`

```yml
version: '3.7'

services:
  elasticsearch:
    image: harbor.cowarobot.cn/docker.io/elasticsearch:8.17.2
    container_name: elasticsearch
    restart: always
    environment:
      - ES_JAVA_OPTS=-Xms6000m -Xmx6000m
    ports:
      - "9200:9200"
      - "9300:9300"
    volumes:
      - /data/elk/es/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml
      - /data/elk/es/data:/usr/share/elasticsearch/data
      - /data/elk/es/plugins:/usr/share/elasticsearch/plugins
    networks:
      - elk
```

进入容器，生成证书

```bash
docker exec -it elasticsearch bash

# 创建CA，输入密码时直接回车，无需密码即可
./bin/elasticsearch-certutil ca --days 3650

# 生成证书，该证书可应用到集群中所有节点，仅需生成一次即可，输入密码时同样直接回车即可
./bin/elasticsearch-certutil cert --ca elastic-stack-ca.p12 --days 3650

# 会生成两个文件，CA证书elastic-stack-ca.p12和自签证书elastic-certificates.p12，将其拷贝出容器外
# 仅将elastic-certificates.p12文件拷贝至其它node节点的 /data/elk/es/config/ 目录中即可
docker cp elasticsearch:/usr/share/elasticsearch/elastic-certificates.p12 /data/elk/es/config/
docker cp elasticsearch:/usr/share/elasticsearch/elastic-stack-ca.p12 /data/elk/es/config/   
```

# Configuration ES node

仅 `172.16.10.20` 配置kibana即可

`docker-compose.yml`

```yaml
version: '3.7'

services:
  elasticsearch:
    image: harbor.cowarobot.cn/docker.io/elasticsearch:8.17.2
    container_name: elasticsearch
    restart: always
    environment:
      - ES_JAVA_OPTS=-Xms6000m -Xmx6000m
    ports:
      - "9200:9200"
      - "9300:9300"
    volumes:
      - /data/elk/es/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml
      - /data/elk/es/data:/usr/share/elasticsearch/data
      - /data/elk/es/config/elastic-certificates.p12:/usr/share/elasticsearch/config/elastic-certificates.p12
      - /data/elk/es/config/elastic-stack-ca.p12:/usr/share/elasticsearch/config/elastic-stack-ca.p12
      - /data/elk/es/plugins:/usr/share/elasticsearch/plugins
    networks:
      - elk

  kibana:
    image: harbor.cowarobot.cn/docker.io/kibana:8.17.2
    hostname: kibana
    container_name: kibana
    volumes: 
      - ./es/config/kibana.yml:/usr/share/kibana/config/kibana.yml
    ports:
      - "5601:5601"
    restart: always
    networks:
      - elk

networks:
  elk:
    driver: bridge
```


## 172.16.10.20
> xpack.security.transport.ssl.keystore.path 配置项默认会在目录 `/usr/share/elasticsearch/` 中找证书文件

```yml
cluster.name: soft-es
node.name: es-node1
node.roles: [ master, data ]
network.host: 0.0.0.0
network.publish_host: 172.16.10.20
http.port: 9200
http.cors.enabled: true
http.cors.allow-origin: "*"
ingest.geoip.downloader.enabled: false
discovery.seed_hosts: ["172.16.10.20:9300", "172.16.10.21:9300", "172.16.10.22:9300"]
cluster.initial_master_nodes: ["es-node1","es-node2","es-node3"]
xpack.security.enabled: true
xpack.security.transport.ssl.enabled: true
xpack.security.transport.ssl.verification_mode: certificate
xpack.security.transport.ssl.client_authentication: required
xpack.security.transport.ssl.keystore.path: elastic-certificates.p12
xpack.security.transport.ssl.truststore.path: elastic-certificates.p12
```

## 172.16.10.21

```yml
cluster.name: soft-es
node.name: es-node2
node.roles: [ master, data ]
network.host: 0.0.0.0
network.publish_host: 172.16.10.21
http.port: 9200
http.cors.enabled: true
http.cors.allow-origin: "*"
ingest.geoip.downloader.enabled: false
discovery.seed_hosts: ["172.16.10.20:9300", "172.16.10.21:9300", "172.16.10.22:9300"]
cluster.initial_master_nodes: ["es-node1","es-node2","es-node3"]
xpack.security.enabled: true
xpack.security.transport.ssl.enabled: true
xpack.security.transport.ssl.verification_mode: certificate
xpack.security.transport.ssl.client_authentication: required
xpack.security.transport.ssl.keystore.path: elastic-certificates.p12
xpack.security.transport.ssl.truststore.path: elastic-certificates.p12
```

## 172.16.10.22

```yml
cluster.name: soft-es
node.name: es-node3
node.roles: [ master, data ]
network.host: 0.0.0.0
network.publish_host: 172.16.10.22
http.port: 9200
http.cors.enabled: true
http.cors.allow-origin: "*"
ingest.geoip.downloader.enabled: false
discovery.seed_hosts: ["172.16.10.20:9300", "172.16.10.21:9300", "172.16.10.22:9300"]
cluster.initial_master_nodes: ["es-node1","es-node2","es-node3"]
xpack.security.enabled: true
xpack.security.transport.ssl.enabled: true
xpack.security.transport.ssl.verification_mode: certificate
xpack.security.transport.ssl.client_authentication: required
xpack.security.transport.ssl.keystore.path: elastic-certificates.p12
xpack.security.transport.ssl.truststore.path: elastic-certificates.p12
```

# 配置密码

> 进入 172.16.10.20 的es容器中生成密码接口，仅需在一台节点中生成密码，生成密码后重启所有es node

```bash
docker exec -it elasticsearch bash
./bin/elasticsearch-setup-passwords interactive  # 手动输入密码

# 重启所有节点es node
docker-compose restart elasticsearch
```

# kibana配置

```yml
i18n.locale: zh-CN
server.host: "0.0.0.0"
server.shutdownTimeout: "5s"
elasticsearch.hosts: [ "http://172.16.10.20:9200" ]
monitoring.ui.container.elasticsearch.enabled: true
elasticsearch.username: "kibana"
elasticsearch.password: "xxxxxxxxxx"
```
