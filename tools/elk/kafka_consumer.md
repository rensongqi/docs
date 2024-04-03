
# 消费者

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