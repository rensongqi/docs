
# 参考
- [github elasticsearch_exporter](https://github.com/prometheus-community/elasticsearch_exporter)
- [dashboard](https://grafana.com/grafana/dashboards/9746-elasticsearch-example/)

docker-compose.yml

```yml
version: "3"

services:
  elasticsearch_exporter:
    image: elasticsearch-exporter:latest
    command:
     - '--es.uri=http://172.16.100.107:9200'
    restart: always
    ports:
    - "9114:9114"
```

prometheus config
```
- job_name: "elasticsearch1"
  static_configs:
  - targets: ["172.16.100.107:9114"]
- job_name: "elasticsearch2"
  static_configs:
  - targets: ["172.16.100.108:9114"]
- job_name: "elasticsearch3"
  static_configs:
  - targets: ["172.16.100.109:9114"]
```