
# 常用命令

```bash
# 创建topic
kafka-topics.sh --create --bootstrap-server localhost:9092 --topic filebeat

# 删除topic
kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic filebeat

# 列出当前所有topic
kafka-topics.sh --list --bootstrap-server localhost:9092

# 查看指定topic
kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic filebeat

# 查询指定group消费情况
kafka-consumer-groups.sh --bootstrap-server localhost:9092 --describe --group logstash

# 创建指定分区数量的topic
kafka-topics.sh --create --topic rensongqi --bootstrap-server localhost:9092 --partitions 3 --replication-factor 3

# 给指定topic增加分区，分区只能增加，不能缩小
kafka-topics.sh --bootstrap-server 172.16.100.67:9092 --alter --topic rensongqi --partitions 3
```