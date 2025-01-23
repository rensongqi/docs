
docker-compose.yml

```yml
version: "3"

services:
  kafka_exporter:
    image: kafka-exporter:latest
    command:
     - '--kafka.server=172.16.100.107:9092'
     - '--kafka.server=172.16.100.108:9092'
     - '--kafka.server=172.16.100.109:9092'
    restart: always
    ports:
    - "9308:9308"
```

prometheus config
```
- job_name: "elk-kafka"
  static_configs:
  - targets: ["172.16.100.107:9308"]
```