# Ubuntu

```bash
cat >>/etc/security/limits.conf<<EOF
* soft nofile 1055350
* hard nofile 1055350
* soft nproc 1055350
* hard nproc 1055350
* soft memlock unlimited
* hard memlock unlimited
EOF
# 修改文件最大打开数量，root & other users
echo 'DefaultLimitNOFILE=1055350' >> /etc/systemd/system.conf
echo 'DefaultLimitNOFILE=1055350' >> /etc/systemd/user.conf
echo "* - nofile 1055350" > /etc/security/limits.d/nofile.conf
echo "root - nofile 1055350" >> /etc/security/limits.d/nofile.conf

cat >/etc/sysctl.conf<<EOF
fs.file-max = 4194303
vm.max_map_count=262144
vm.swappiness = 1
vm.vfs_cache_pressure = 50
vm.min_free_kbytes = 1000000
net.ipv4.tcp_timestamps = 0
net.ipv4.tcp_sack = 1
net.core.netdev_max_backlog = 250000
net.core.rmem_max = 4194304
net.core.wmem_max = 4194304
net.core.rmem_default = 4194304
net.core.wmem_default = 4194304
net.core.optmem_max = 4194304
net.ipv4.tcp_low_latency = 1
net.ipv4.tcp_adv_win_scale = 1
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 10000
net.ipv4.tcp_max_syn_backlog = 4096
net.ipv4.tcp_fin_timeout = 15
net.ipv4.conf.all.send_redirects = 0
net.ipv4.conf.all.accept_redirects = 0
net.ipv4.conf.all.accept_source_route = 0
net.ipv4.tcp_mtu_probing = 1
fs.inotify.max_user_watches=524288
EOF
sysctl -p
```

# Centos7

如果出现节点仅收到syn的请求，但是并没有对这些请求做ack响应，则修改如下配置

`/etc/sysctl.conf`

```bash
net.ipv4.tcp_fin_timeout = 2
net.ipv4.tcp_tw_reuse = 1
#net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_tw_recycle = 0
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_keepalive_time = 600
net.ipv4.ip_local_port_range = 4000 65000
net.ipv4.tcp_max_syn_backlog = 16384
net.ipv4.tcp_max_tw_buckets = 36000
net.ipv4.route.gc_timeout = 100
net.ipv4.tcp_syn_retries = 1
net.ipv4.tcp_synack_retries = 1
net.core.somaxconn = 16384
net.core.netdev_max_backlog = 16384
net.ipv4.tcp_max_orphans = 16384
vm.max_map_count=262144
fs.file-max=104857600
net.ipv4.ip_forward=1                      # 1表示开启 0表示关闭
net.ipv4.conf.default.rp_filter=0     
net.ipv4.conf.all.rp_filter=0              #控制系统是否开启对数据包源地址的校验
```