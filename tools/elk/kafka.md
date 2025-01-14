- [1 单节点部署](#1-单节点部署)
- [2 多节点部署](#2-多节点部署)
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