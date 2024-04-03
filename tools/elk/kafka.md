
# 1 单节点部署

```yaml
version: '3.5'
services:
  zookeeper:
    image: wurstmeister/zookeeper
    container_name: zookeeper
    ports:
      - "2181:2181"
  kafka:
    image: wurstmeister/kafka
    container_name: kafka
    volumes: 
        - /var/run/docker.sock:/var/run/docker.sock
    ports:
      - "9092:9092"
    environment:
      KAFKA_ADVERTISED_HOST_NAME: 172.16.108.89
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_PORT: 9092
      KAFKA_LOG_RETENTION_HOURS: 120
      KAFKA_MESSAGE_MAX_BYTES: 10000000
      KAFKA_REPLICA_FETCH_MAX_BYTES: 10000000
      KAFKA_GROUP_MAX_SESSION_TIMEOUT_MS: 60000
      KAFKA_NUM_PARTITIONS: 3
      KAFKA_DELETE_RETENTION_MS: 1000
      KAFKA_LISTENERS: PLAINTEXT://:9092
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://172.18.108.89:9092
      KAFKA_BROKER_ID: 1
  kafka-manager:
    image: sheepkiller/kafka-manager
    container_name: kafka-manager
    environment:
        ZK_HOSTS: 172.18.108.89
    ports:  
      - "9009:9000"
```

# 2 多节点部署

```yaml
# 172.18.11.109 节点
version: '3'
services:
  zookeeper:
    image: wurstmeister/zookeeper
    ports:
      - 2182:2181
    restart: always
  
  kafka1:
    image: wurstmeister/kafka
    container_name: kafka1
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 0
      KAFKA_NUM_PARTITIONS: 3
      KAFKA_DEFAULT_REPLICATION_FACTOR: 2
      KAFKA_ZOOKEEPER_CONNECT: 172.18.11.109:2182
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://172.18.11.109:9092
    volumes:
      - ./logs:/opt/kafka/logs
      - /var/run/docker.sock:/var/run/docker.sock
    restart: always
        
# 172.18.11.110 节点
version: '3'
services:
  zookeeper:
    image: wurstmeister/zookeeper
    ports:
      - 2182:2181
    restart: always

  kafka2:
    image: wurstmeister/kafka
    container_name: kafka2
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_NUM_PARTITIONS: 3
      KAFKA_DEFAULT_REPLICATION_FACTOR: 2
      KAFKA_ZOOKEEPER_CONNECT: 172.18.11.109:2182
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://172.18.11.110:9092
    volumes:
      - ./logs:/opt/kafka/logs
      - /var/run/docker.sock:/var/run/docker.sock
    restart: always
        
# 172.18.11.111 节点
version: '3'
services:
  zookeeper:
    image: wurstmeister/zookeeper
    ports:
      - 2182:2181
    restart: always

  kafka3:
    image: wurstmeister/kafka
    container_name: kafka3
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 2
      KAFKA_NUM_PARTITIONS: 3
      KAFKA_DEFAULT_REPLICATION_FACTOR: 2
      KAFKA_ZOOKEEPER_CONNECT: 172.18.11.109:2182
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://172.18.11.111:9092
    volumes:
      - ./logs:/opt/kafka/logs
      - /var/run/docker.sock:/var/run/docker.sock
    restart: always
```

# 3 公网访问kafka

参考文章：[Kafka内外网访问的设置](https://www.cnblogs.com/gentlescholar/p/15179258.html)

如果内网的kafka服务要暴露到公网，需要配置如下

内网IP：`172.18.100.67`, 公网IP：`123.23.4.5`

```
listener.security.protocol.map=INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT
listeners=INTERNAL://172.18.100.67:9092,EXTERNAL://172.18.100.67:19092
advertised.listeners=INTERNAL://172.18.100.67:9092,EXTERNAL://123.23.4.5:19092
inter.broker.listener.name=INTERNAL
```