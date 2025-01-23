
- [github blackbox_exporter](https://github.com/prometheus/blackbox_exporter)
- [dashboard](https://grafana.com/grafana/dashboards/7587-prometheus-blackbox-exporter/)

docker-compose.yml

```yml
version: "3"

services:
  blackbox_exporter:
    image: blackbox-exporter:latest
    command:
     - '--config.file=/config/blackbox.yml'
    restart: always
    volumes:
      - /data/blackbox_exporter/config:/config
    ports:
    - "9115:9115"
```

config.yml
- [example.yml](https://github.com/prometheus/blackbox_exporter/blob/master/example.yml)
```
modules:
  http_2xx:
    prober: http
    timeout: 3s
    http:
      method: GET
```

prometheus config

```
- job_name: 'http_service_probe'
  metrics_path: /probe
  params:
    module: [http_2xx]
  static_configs:
  - targets:
    - http://172.16.104.21:8893/healthy
  relabel_configs:
  - source_labels: [__address__]
    target_label: __param_target
  - source_labels: [__param_target]
    target_label: instance
  - target_label: __address__
    replacement: 172.16.104.21:9115
```