
# 1 Redis主从 + 哨兵配置

初始化

```bash
mkdir -p /data/redis/sentinel{1..3}
mkdir -p /data/redis/{master_data,slave1_data,slave2_data}
chmod 777 /data/redis/sentinel
```

> 注意需要手动修改 xxxxxxxx 为自定义的密码

`/data/redis/docker-compose.yml`
```yml
version: '3.8'

networks:
  redis-network:
    driver: bridge
    ipam:
      driver: default
      config:
        - subnet: 172.30.1.0/24

services:
  redis-master:
    container_name: redis-master
    image: redis:7.4.0-alpine
    volumes:
      - ./master_data:/data
    ports:
      - "7001:6379"
    command: redis-server --port 6379 --requirepass xxxxxxxx --appendonly yes --protected-mode no
    networks:
      redis-network:
        ipv4_address: 172.30.1.2

  redis-replica1:
    container_name: redis-replica1
    image: redis:7.4.0-alpine
    volumes:
      - ./slave1_data:/data
    ports:
      - "7002:6379"
    command: redis-server --slaveof 172.30.1.2 6379 --port 6379 --requirepass xxxxxxxx --masterauth xxxxxxxx --appendonly yes --protected-mode no
    depends_on:
      - redis-master
    networks:
      redis-network:
        ipv4_address: 172.30.1.3

  redis-replica2:
    container_name: redis-replica2
    image: redis:7.4.0-alpine
    volumes:
      - ./slave1_data:/data
    ports:
      - "7003:6379"
    command: redis-server --slaveof 172.30.1.2 6379 --port 6379 --requirepass xxxxxxxx --masterauth xxxxxxxx --appendonly yes --protected-mode no
    depends_on:
      - redis-master
    networks:
      redis-network:
        ipv4_address: 172.30.1.4

  redis-sentinel1:
    container_name: redis-sentinel1
    image: redis:7.4.0-alpine
    volumes:
      - ./sentinel1:/data/sentinel
    ports:
      - "27001:26379"
    command: ["redis-sentinel", "/data/sentinel/sentinel.conf"]
    depends_on:
      - redis-master
    networks:
      redis-network:
        ipv4_address: 172.30.1.11

  redis-sentinel2:
    container_name: redis-sentinel2
    image: redis:7.4.0-alpine
    volumes:
      - ./sentinel2:/data/sentinel
    ports:
      - "27002:26379"
    command: [ "redis-sentinel", "/data/sentinel/sentinel.conf" ]
    depends_on:
      - redis-master
    networks:
      redis-network:
        ipv4_address: 172.30.1.12

  redis-sentinel3:
    container_name: redis-sentinel3
    image: redis:7.4.0-alpine
    volumes:
      - ./sentinel3:/data/sentinel
    ports:
      - "27003:26379"
    command: [ "redis-sentinel", "/data/sentinel/sentinel.conf" ]
    depends_on:
      - redis-master
    networks:
      redis-network:
        ipv4_address: 172.30.1.13
```

在每一个sentinel.conf的配置文件中添加如下相同的内容

`/data/redis/sentinel{1..3}/sentinel.conf`

```bash
dir "/data/sentinel"
sentinel monitor mymaster 172.30.1.2 6379 2
sentinel down-after-milliseconds mymaster 5000
sentinel failover-timeout mymaster 60000
sentinel auth-pass mymaster xxxxxxxx
```

启动服务

```bash
cd /data/redis
docker-compose up -d redis-master
docker-compose up -d redis-replica1
docker-compose up -d redis-replica2
docker-compose up -d redis-sentinel1
docker-compose up -d redis-sentinel2
docker-compose up -d redis-sentinel3
```

参考文章

- [Redis 7.x 哨兵配置](https://juejin.cn/post/7417635848987164687)
- [Docker-Compose部署Redis(v7.2)哨兵模式](https://blog.csdn.net/m0_51390969/article/details/135413933)

# 2 单节点实例

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