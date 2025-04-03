- [1 单节点部署](#1-单节点部署)
- [2 多节点部署](#2-多节点部署)
  - [2.1 基于Zookeeper](#21-基于zookeeper)
  - [2.2 基于Kraft](#22-基于kraft)
- [3 公网访问kafka](#3-公网访问kafka)
- [4 生产者](#4-生产者)
- [5 消费者](#5-消费者)

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

## 2.1 基于Zookeeper

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

监控：

```yml
version: '3'
services:
  kafka_exporter:
    image: kafka-exporter:latest
    command:
     - '--kafka.server=172.18.11.109:9092'
     - '--kafka.server=172.18.11.110:9092'
     - '--kafka.server=172.18.11.111:9092'
    restart: always
    ports:
    - "9308:9308"
```

Grafana看板ID
- [7589-kafka-exporter-overview](https://grafana.com/grafana/dashboards/7589-kafka-exporter-overview/)

## 2.2 基于Kraft

> 需要注意的是，基于Kraft的Kafka创建进入容器之后不能像使用zookeeper那样创建topic，具体原因可参考文章：https://github.com/wurstmeister/kafka-docker/issues/390
> 
> 解决办法，进入容器之后需要执行命令 `unset KAFKA_OPTS` ，然后才能执行如下命令

初始化
```
mkdir /data/kafka/kraft -p

# 下载jmx_prometheus_javaagent-0.20.0.jar
cd /data/kafka/
wget https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/0.20.0/jmx_prometheus_javaagent-0.20.0.jar
wget https://raw.githubusercontent.com/prometheus/jmx_exporter/refs/heads/main/examples/kafka-kraft-3_0_0.yml
```

`/data/kafka/docker-compose.yml`
```yaml
# 172.16.10.85
version: "3"
services:
  kafka:
    image: bitnami/kafka:3.7.0
    network_mode: host
    container_name: kafka
    user: root
    ports:
      - 9092:9092
      - 9093:9093
    environment:
      - TZ=Asia/Shanghai
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_LOG_RETENTION_HOURS=72
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_KRAFT_CLUSTER_ID=5L6g3nShT-eMCtK--X86sw
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@172.16.10.85:9093,2@172.16.10.86:9093,3@172.16.10.87:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_HEAP_OPTS=-Xmx10G -Xms10G
      - KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=false
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://172.16.10.85:9092
      - KAFKA_CFG_NODE_ID=1
      - KAFKA_OPTS=-javaagent:/opt/jmx_prometheus_javaagent-0.20.0.jar=9999:/opt/kafka-kraft-3_0_0.yml
    volumes:
      - ./kraft:/bitnami/kafka/data:rw
      - ./jmx_prometheus_javaagent-0.20.0.jar:/opt/jmx_prometheus_javaagent-0.20.0.jar
      - ./kafka-kraft-3_0_0.yml:/opt/kafka-kraft-3_0_0.yml

# 172.16.10.86
version: "3"
services:
  kafka:
    image: bitnami/kafka:3.7.0
    network_mode: host
    container_name: kafka
    user: root
    ports:
      - 9092:9092
      - 9093:9093
    environment:
      - TZ=Asia/Shanghai
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_LOG_RETENTION_HOURS=72
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_KRAFT_CLUSTER_ID=5L6g3nShT-eMCtK--X86sw
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@172.16.10.85:9093,2@172.16.10.86:9093,3@172.16.10.87:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_HEAP_OPTS=-Xmx10G -Xms10G
      - KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=false
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://172.16.10.86:9092
      - KAFKA_CFG_NODE_ID=2
      - KAFKA_OPTS=-javaagent:/opt/jmx_prometheus_javaagent-0.20.0.jar=9999:/opt/kafka-kraft-3_0_0.yml
    volumes:
      - ./kraft:/bitnami/kafka/data:rw
      - ./jmx_prometheus_javaagent-0.20.0.jar:/opt/jmx_prometheus_javaagent-0.20.0.jar
      - ./kafka-kraft-3_0_0.yml:/opt/kafka-kraft-3_0_0.yml

# 172.16.10.87
version: "3"
services:
  kafka:
    image: bitnami/kafka:3.7.0
    network_mode: host
    container_name: kafka
    user: root
    ports:
      - 9092:9092
      - 9093:9093
    environment:
      - TZ=Asia/Shanghai
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_LOG_RETENTION_HOURS=72
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_KRAFT_CLUSTER_ID=5L6g3nShT-eMCtK--X86sw
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@172.16.10.85:9093,2@172.16.10.86:9093,3@172.16.10.87:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_HEAP_OPTS=-Xmx10G -Xms10G
      - KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=false
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://172.16.10.87:9092
      - KAFKA_CFG_NODE_ID=3
      - KAFKA_OPTS=-javaagent:/opt/jmx_prometheus_javaagent-0.20.0.jar=9999:/opt/kafka-kraft-3_0_0.yml
    volumes:
      - ./kraft:/bitnami/kafka/data:rw
      - ./jmx_prometheus_javaagent-0.20.0.jar:/opt/jmx_prometheus_javaagent-0.20.0.jar
      - ./kafka-kraft-3_0_0.yml:/opt/kafka-kraft-3_0_0.yml
```

监控：
> 需要基于jmx_exporter实现kraft的监控指标数据的获取
- [prometheus监控Kafka (kafka_exporter和 jmx_exporter)](https://blog.csdn.net/u010533742/article/details/119992040)

Grafana看板ID 
> 下载对应的json文件之后，需要将该文件中所有的 `kafka_controller_kafkacontroller_controllerstate` 改为 `kafka_controller_kafkacontroller_activecontrollercount` 后看板才可使用（实际上对应的就是获取看板变量的label），修改的值要保证在对应的exporter metrics中搜索到，不一定非要修改上上述值，了解即可
- [11962-kafka-metrics](https://grafana.com/grafana/dashboards/11962-kafka-metrics/)

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

# 4 生产者
```go
package main
import (
   "fmt"
   "github.com/IBM/sarama"
   "os"
)
func main() {
   // Kafka broker 地址
   //brokers := []string{"172.16.100.107:9092"}
   brokers := []string{"172.18.100.67:9092"}
   // 新建 sarama 配置实例
   config := sarama.NewConfig()
   config.Producer.RequiredAcks = sarama.WaitForAll          // 等待所有副本都保存成功
   config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机的分区方式
   config.Producer.Return.Successes = true                   // 成功交付的消息将在 success channel 返回
   // 使用配置,新建一个异步生产者
   producer, err := sarama.NewAsyncProducer(brokers, config)
   if err != nil {
      fmt.Println("Failed to start producer:", err)
      os.Exit(1)
   }
   // 构建发送的消息，
   msg := &sarama.ProducerMessage{
      Topic: "ultron",
      Value: sarama.StringEncoder("测试消息"),
   }
   // 发送消息
   producer.Input() <- msg
   // 等待消息发送完成
   select {
   case suc := <-producer.Successes():
      fmt.Printf("offset: %d,  timestamp: %s\n", suc.Offset, suc.Timestamp.String())
   case fail := <-producer.Errors():
      fmt.Printf("err: %s\n", fail.Err.Error())
   }
}
```

# 5 消费者
```go
package main
import (
   "fmt"
   "github.com/IBM/sarama"
   "os"
   "os/signal"
   "syscall"
)
func saramaKafka() {
   config := sarama.NewConfig()
   config.Consumer.Return.Errors = true
   brokers := []string{"172.18.100.67:9092"}
   topic := "ultron"
   // 创建消费者
   master, err := sarama.NewConsumer(brokers, config)
   if err != nil {
      panic(err)
   }
   defer func() {
      if err := master.Close(); err != nil {
         panic(err)
      }
   }()
   consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
   if err != nil {
      panic(err)
   }
   signals := make(chan os.Signal, 1)
   signal.Notify(signals, syscall.SIGTERM)
   consumed := 0
ConsumerLoop:
   for {
      select {
      case msg := <-consumer.Messages():
         fmt.Printf("Consumed message offset %d\n", string(msg.Value))
         consumed++
      case <-signals:
         break ConsumerLoop
      }
   }
   fmt.Printf("Consumed: %d\n", consumed)
}
func main() {
   saramaKafka()
}
```