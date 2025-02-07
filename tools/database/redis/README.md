
# 1 主从配置
docker-compose.yml
```yml
version: '3'
services:
  master:
    image: docker.io/redis:7.4.0-alpine
    container_name: redis-master
    restart: always
    command: redis-server --port 6379 --requirepass xxxxxxxx --appendonly yes
    ports:
      - 6379:6379
    volumes:
      - ./master_data:/data

  slave1:
    image: docker.io/redis:7.4.0-alpine
    container_name: redis-slave-1
    restart: always
    command: redis-server --slaveof 172.16.90.178 6379 --port 6380 --requirepass xxxxxxxx --masterauth xxxxxxxx --appendonly yes
    ports:
      - 6380:6380
    volumes:
      - ./slave1_data:/data

  slave2:
    image: docker.io/redis:7.4.0-alpine
    container_name: redis-slave-2
    restart: always
    command: redis-server --slaveof 172.16.90.178 6379 --port 6381 --requirepass xxxxxxxx --masterauth xxxxxxxx --appendonly yes
    ports:
      - 6381:6381
    volumes:
      - ./slave2_data:/data
```

# 2 哨兵

初始化
```bash
mkdir -p /data/redis/sentinel{1..3}
chmod 777 /data/redis/sentinel
```

`/data/redis/docker-compose.yml`

```yml
version: '3.4'
services:
  sentinel1:
    image: docker.io/redis:7.4.0-alpine
    container_name: redis-sentinel-1
    #user: "0:0"
    ports:
      - 26379:26379
    command: redis-sentinel /data/sentinel/sentinel.conf --sentinel --loglevel verbose
    restart: always
    volumes:
      - ./sentinel1:/data/sentinel

  sentinel2:
    image: docker.io/redis:7.4.0-alpine
    container_name: redis-sentinel-2
    ports:
      - 26380:26379
    command: redis-sentinel /data/sentinel/sentinel.conf --sentinel
    restart: always
    volumes:
      - ./sentinel2:/data/sentinel

  sentinel3:
    image: docker.io/redis:7.4.0-alpine
    container_name: redis-sentinel-3
    ports:
      - 26381:26379
    command: redis-sentinel /data/sentinel/sentinel.conf --sentinel
    restart: always
    volumes:
      - ./sentinel3:/data/sentinel
```

`/data/redis/sentinel{1..3}/sentinel.conf`

```bash
bind 0.0.0.0
daemonize yes
protected-mode no
port 26379
dir "/data/sentinel"
pidfile "/var/run/redis-sentinel.pid"
syslog-enabled no
sentinel deny-scripts-reconfig yes
sentinel monitor mymaster x.x.x.x 6379 2
sentinel auth-pass mymaster xxxxxxxx
```

[Redis 7.x 哨兵配置](https://juejin.cn/post/7417635848987164687)

# 3 单节点实例

```bash
version: '3'
services:
  master:
    image: docker.io/redis:7.4.0-alpine
    container_name: redis-master
    restart: always
    command: redis-server --port 6379 --requirepass xxxxxxxx  --appendonly yes
    ports:
      - 6379:6379
    volumes:
      - ./data:/data
```