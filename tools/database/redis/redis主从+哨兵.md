
# 1 主从
docker-compose.yml
```yml
version: '3'
services:
  master:
    image: harbor.rsq.cn/library/redis/redis:6.2.12-alpine3.18
    container_name: redis-master
    restart: always
    command: redis-server --port 6379 --requirepass cowa2022  --appendonly yes
    ports:
      - 6379:6379
    volumes:
      - ./data:/data

  slave1:
    image: harbor.rsq.cn/library/redis/redis:6.2.12-alpine3.18
    container_name: redis-slave-1
    restart: always
    command: redis-server --slaveof 172.16.104.109 6379 --port 6380  --requirepass cowa2022 --masterauth cowa2022  --appendonly yes
    ports:
      - 6379:6379
    volumes:
      - ./data:/data


  slave2:
    image: harbor.rsq.cn/library/redis/redis:6.2.12-alpine3.18
    container_name: redis-slave-2
    restart: always
    command: redis-server --slaveof 172.16.104.109 --port 6381  --requirepass cowa2022 --masterauth cowa2022  --appendonly yes
    ports:
      - 6379:6379
    volumes:
      - ./data:/data
```

# 2 哨兵

docker-compose.yml

```yml
version: '3.4'
services:
  sentinel1:
    image: harbor.rsq.cn/library/redis/redis:6.2.12-alpine3.18
    container_name: redis-sentinel-1
    ports:
      - 26379:26379
    command: redis-sentinel /data/sentinel.conf
    restart: always
    volumes:
      - ./sentinel.conf:/data/sentinel.conf
  sentinel2:
    image: harbor.rsq.cn/library/redis/redis:6.2.12-alpine3.18
    container_name: redis-sentinel-2
    ports:
      - 26379:26379
    command: redis-sentinel /data/sentinel.conf
    restart: always
    volumes:
      - ./sentinel.conf:/data/sentinel.conf
  sentinel3:
    image: harbor.rsq.cn/library/redis/redis:6.2.12-alpine3.18
    container_name: redis-sentinel-3
    ports:
      - 26379:26379
    command: redis-sentinel /data/sentinel.conf
    restart: always
    volumes:
      - ./sentinel.conf:/data/sentinel.conf
```

sentinel.conf
```bash
port 26379
dir "/tmp"
sentinel deny-scripts-reconfig yes
sentinel monitor mymaster 172.16.104.111 6379 2
sentinel auth-pass mymaster cowa2022
```

# 3 单节点实例

```bash
version: '3'
services:
  master:
    image: harbor.rsq.cn/library/redis/redis:6.2.12-alpine3.18
    container_name: redis-master
    restart: always
    command: redis-server --port 6379 --requirepass cowa2022  --appendonly yes
    ports:
      - 6379:6379
    volumes:
      - ./data:/data
```