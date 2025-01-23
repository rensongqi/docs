# k8s deploy
```bash
kubectl create namespace monitoring-system
kubectl apply node_exporter.yaml
```

# 监控IB网卡流量

> https://github.com/prometheus/node_exporter/blob/master/collector/infiniband_linux.go#L73

```
# 指标
node_infiniband_port_data_transmitted_bytes_total
node_infiniband_port_data_received_bytes_total
node_infiniband_port_packets_transmitted_total
node_infiniband_port_packets_received_total
node_infiniband_unicast_packets_transmitted_total
node_infiniband_unicast_packets_received_total
node_infiniband_rate_bytes_per_second

# promtql 语法
rate(node_infiniband_port_data_transmitted_bytes_total{instance=~"$instance", device=~"$mlx"}[$__rate_interval])
node_infiniband_port_data_transmitted_bytes_total{instance=~"$instance", device=~"$mlx"}
{{instance}}/{{device}}/{{port}}
```