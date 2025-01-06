# Test rdma


## Getting started

Installation dependency

```bash
sudo apt-get update
sudo apt-get install -y libibverbs-dev librdmacm-dev perftest rdma-core
```

Test

`rdma.c`
```bash
# 修改run_server()和run_client()方法中 local_info.lid 为正确的lid值

# build
gcc -o rdma  rdma.c  -libverbs -o rdma

# server 172.17.1.237
./rdma

# client 172.17.1.246
./rdma 172.17.1.237
```

`rdma_copy_file.c`

```bash
# 修改run_server()和run_client()方法中 local_info.lid 为正确的lid值

# build
gcc -o rdma_file file.c -libverbs -o rdma_file

# server
./rdma_file

# client
./rdma_file 172.17.1.237 record.video.00028

# 测试拷贝时间
# 本地文件走TCP/IP协议
time rsync -avP record.video.00028 root@172.17.1.237:/tmp/

# 使用rdma协议
time ./rdma_file 172.17.1.237 record.video.00028

```

测试效果，对于一个2GB的大文件，可以发现走rdma协议比走TCP协议快了接近2s。

![test](../../img/test.jpg)

同时也可以发现，走rdma协议基本没有消耗sys的时间，消耗的一点sys时间基本为read文件的时间，通过strace可以跟踪程序执行期间调用的系统调用，分析哪些内核操作消耗了时间。

![strace_rdma](../../img/strace_rdma.jpg)

![strace_rsync](../../img/strace_rsync.jpg)

注意 `#define MSG_SIZE    1024` MSG_SIZE的大小是决定文件拷贝性能的因素

## CMD

```bash
# 查看ib状态
ibstat

# 查看ib连接信息
iblinkinfio

# 查看ib卡详细信息
ibv_devinfo
ibv_devices

# 延迟测试
# server
ibping -S

# client
ibping <server lid port>

```