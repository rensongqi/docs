
# 生产者
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