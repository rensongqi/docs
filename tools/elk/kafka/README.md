
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
```