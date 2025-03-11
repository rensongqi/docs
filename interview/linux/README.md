# 遇到过的印象比较深刻的问题


## 运维侧

Centos dmz网络问题，首先背景是我们的IDC机房有自己的出口防火墙和固定的公网IP，内部的服务映射至公网之后外部访问公网服务没有问题，内部访问映射出去的公网服务就出现了问题，原因在于防火墙做了地址转换后，内网ng有一项内核参数没有调整好，就是禁用反向路径过滤，反向路径过滤是一种安全机制，用于验证接收到的数据包的源地址是否可通过接收接口到达。net.ipv4.conf.default.rp_filter=0 关闭之后解决了该问题

`vim /etc/sysctl.conf`
```bash
net.ipv4.ip_forward=1
net.ipv4.conf.default.rp_filter=0     
net.ipv4.conf.all.rp_filter=0
```

## 开发侧

1. 内存泄漏问题，在读写数据中有些内存未释放

2. docker ansible ssh执行playbook完会后会遗留下来一批僵尸进程不会被回收，宿主机运行playboos后进程则会被init进程回收
    > https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/

    解决办法：docker-compose运行容器时指定由init为true或使用tini进程接手编译好的二进制程序

