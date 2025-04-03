
- [172.16.10.81](#172161081)
  - [启动MySQL](#启动mysql)
  - [启动Nacos](#启动nacos)
- [172.16.10.82](#172161082)
- [172.16.10.83](#172161083)
- [添加监控](#添加监控)

|IP   | 服务 | 部署 |
| --- | ---- | --- |
| 172.16.10.81 | mysql、nacos |
| 172.16.10.82 | nacos |
| 172.16.10.83 | nacos | 

# 172.16.10.81

```
version: "3.8"
services:
  nacos:
    hostname: nacos
    container_name: nacos
    network_mode: host
    image: nacos/nacos-server:v2.3.2
    volumes:
      - /data/nacos/data:/home/nacos/logs
    ports:
      - "7848:7848"
      - "8848:8848"
      - "9848:9848"
      - "9849:9849"
    restart: always
    environment:
      - NACOS_SERVER_IP=172.16.10.81
      - NACOS_SERVERS=172.16.10.81:8848 172.16.10.82:8848 172.16.10.83:8848
      - SPRING_DATASOURCE_PLATFORM=mysql
      - MYSQL_SERVICE_HOST=172.16.10.81
      - MYSQL_SERVICE_DB_NAME=nacos
      - MYSQL_SERVICE_PORT=3306
      - MYSQL_SERVICE_USER=nacos
      - MYSQL_SERVICE_PASSWORD=nacos  # mysql pass
      - MYSQL_SERVICE_DB_PARAM=characterEncoding=utf8&connectTimeout=1000&socketTimeout=3000&autoReconnect=true&useSSL=false&allowPublicKeyRetrieval=true
      - NACOS_AUTH_ENABLE=true
      - NACOS_AUTH_TOKEN=Mkc2bnRlcFlDcUVXTlNqenlaa3Zxxxxxxxxxx=
      - NACOS_AUTH_IDENTITY_KEY=nacos
      - NACOS_AUTH_IDENTITY_VALUE=xxxxxx

  mysql:
    container_name: nacos_mysql
    image: nacos/nacos-mysql:8.0.31
    env_file:
      - ./env
    volumes:
      - /data/nacos/mysql/conf/my.cnf:/etc/my.cnf
      - /data/nacos/mysql/datadir:/var/lib/mysql
      - /etc/localtime:/etc/localtime
      - /usr/share/zoneinfo/Asia/Shanghai:/etc/timezone
    ports:
      - "3306:3306"
    healthcheck:
      test: [ "CMD", "mysqladmin", "ping", "-h", "localhost" ]
      interval: 5s
      timeout: 10s
      retries: 10
```

/data/nacos/env

```
MYSQL_ROOT_PASSWORD=xxxxxxx
MYSQL_DATABASE=nacos
MYSQL_USER=nacos
MYSQL_PASSWORD=nacos
```

## 启动MySQL
> 使用nacos/nacos-mysql:8.0.31镜像会自动初始化好nacos所需要的库和表，但是当第一次启动完Nacos之后会出现账号密码不对的情况，需要手动初始化一个账号密码
>

进入mysql容器
```bash
# 启动mysql容器
cd /data/nacos
docker-compose up -d mysql

# 进入mysql容器
docker exec -it nacos_mysql bash
mysql -u root -p xxxxx

# 以下是SQL命令
show databases;
use nacos;

# users表新增一条nacos用户和密码记录
# XYLb8fvPnnpZRC4N : $2a$10$2k0x0hKgIMqhTTuwHFnLAOxGwkKFKgUkRN2c1l62qfnzwkY2oJ1EW
INSERT INTO users (username, password, enabled) VALUES ('nacos', '$2a$10$2k0x0hKgIMqhTTuwHFnLAOxGwkKFKgUkRN2c1l62qfnzwkY2oJ1EW', TRUE);

INSERT INTO roles (username, role) VALUES ('admin', 'ROLE_ADMIN');
```

## 启动Nacos
```
cd /data/nacos
docker-compose up -d nacos
```

# 172.16.10.82

```
version: "3.8"
services:
  nacos:
    hostname: nacos
    container_name: nacos
    network_mode: host
    image: nacos/nacos-server:v2.3.2
    volumes:
      - /data/nacos/data:/home/nacos/logs
    restart: always
    extra_hosts:
      - "host.docker.internal:host-gateway"
    environment:
      - NACOS_SERVER_IP=172.16.10.82
      - NACOS_SERVERS=172.16.10.81:8848 172.16.10.82:8848 172.16.10.83:8848
      - SPRING_DATASOURCE_PLATFORM=mysql
      - MYSQL_SERVICE_HOST=172.16.10.81
      - MYSQL_SERVICE_DB_NAME=nacos
      - MYSQL_SERVICE_PORT=3306
      - MYSQL_SERVICE_USER=nacos
      - MYSQL_SERVICE_PASSWORD=nacos
      - MYSQL_SERVICE_DB_PARAM=characterEncoding=utf8&connectTimeout=1000&socketTimeout=3000&autoReconnect=true&useSSL=false&allowPublicKeyRetrieval=true
      - NACOS_AUTH_ENABLE=true
      - NACOS_AUTH_TOKEN=Mkc2bnRlcFlDcUVXTlNqenlaa3Zxxxxxxxxxx
      - NACOS_AUTH_IDENTITY_KEY=nacos
      - NACOS_AUTH_IDENTITY_VALUE=xxxxxx
```

# 172.16.10.83

```
version: "3.8"
services:
  nacos:
    hostname: nacos
    container_name: nacos
    network_mode: host
    image: nacos/nacos-server:v2.3.2
    volumes:
      - /data/nacos/data:/home/nacos/logs
    restart: always
    extra_hosts:
      - "host.docker.internal:host-gateway"
    environment:
      - NACOS_SERVER_IP=172.16.10.83
      - NACOS_SERVERS=172.16.10.81:8848 172.16.10.82:8848 172.16.10.83:8848
      - SPRING_DATASOURCE_PLATFORM=mysql
      - MYSQL_SERVICE_HOST=172.16.10.81
      - MYSQL_SERVICE_DB_NAME=nacos
      - MYSQL_SERVICE_PORT=3306
      - MYSQL_SERVICE_USER=nacos
      - MYSQL_SERVICE_PASSWORD=nacos
      - MYSQL_SERVICE_DB_PARAM=characterEncoding=utf8&connectTimeout=1000&socketTimeout=3000&autoReconnect=true&useSSL=false&allowPublicKeyRetrieval=true
      - NACOS_AUTH_ENABLE=true
      - NACOS_AUTH_TOKEN=Mkc2bnRlcFlDcUVXTlNqenlaa3Zxxxxxxxxxx
      - NACOS_AUTH_IDENTITY_KEY=nacos
      - NACOS_AUTH_IDENTITY_VALUE=xxxxxx
```

# 添加监控

> [Nacos 监控手册](https://nacos.io/docs/latest/guide/admin/monitor-guide/)
> 
> [nacos-grafana](https://raw.githubusercontent.com/nacos-group/nacos-template/refs/heads/master/nacos-grafana.json)

进容器，修改 `conf/application.properties` 暴露metrics数据

```
management.endpoints.web.exposure.include=*
```

访问: `x.x.x.x:8848/nacos/actuator/prometheus` 获取监控指标