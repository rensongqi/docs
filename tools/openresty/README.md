- [Nginx 模板](#nginx-模板)
- [Nginx 性能优化](#nginx-性能优化)
  - [调整 Work 进程数](#调整-work-进程数)
  - [绑核](#绑核)
  - [单个进程允许的客户端最大连接数](#单个进程允许的客户端最大连接数)
    - [配置获取更多连接数](#配置获取更多连接数)
  - [Nginx 事件处理模型优化](#nginx-事件处理模型优化)
  - [TCP优化](#tcp优化)
  - [前端性能优化](#前端性能优化)
- [Nginx 安全](#nginx-安全)
  - [隐藏 Nginx 版本号](#隐藏-nginx-版本号)
  - [Nginx 防止 DDos 和 cc 等流量攻击](#nginx-防止-ddos-和-cc-等流量攻击)
      - [限制同一时间段ip访问次数](#限制同一时间段ip访问次数)
      - [禁止ip或ip网段](#禁止ip或ip网段)
      - [屏蔽IP的方法](#屏蔽ip的方法)
- [Lua](#lua)
- [手动容器编译openresty支持slice模块](#手动容器编译openresty支持slice模块)

# Nginx 模板

[layer7.conf](./layer7_tmpl.conf)

# Nginx 性能优化

- Nginx 性能优化的关键点在于减少磁盘IO和网络IO。
- `http` 链接的尽快释放，减少请求的堆积。



## 调整 Work 进程数

Nginx运行工作进程个数一般设置CPU的核心或者核心数x2。

```bash
worker_processes 16; # 指定 Nginx 要开启的进程数，结尾的数字就是进程的个数，可以为 auto
```

如果想省麻烦也可以配置为`worker_processes auto;`，将由 Nginx 自行决定 worker 数量。当访问量快速增加时，Nginx 就会临时 fork 新进程来缩短系统的瞬时开销和降低服务的时间。



## 绑核

因为 Nginx 是多进程模型，每个进程中是单线程的，所以为了避免资源使用不均可以进行绑核操作。

```
worker_processes 4;
worker_cpu_affinity 0001 0010 0100 1000;
```

其中 `worker_cpu_affinity` 就是配置 Nginx 进程与 CPU 亲和力的参数，即把不同的进程分给不同的 CPU 核处理。这里的`0001 0010 0100 1000`是掩码，分别代表第1、2、3、4核CPU。上述配置会为每个进程分配一核CPU处理。

当然，如果想省麻烦也可以配置`worker_cpu_affinity auto;`，将由 Nginx 按需自动分配。



## 单个进程允许的客户端最大连接数

通过调整控制连接数的参数来调整 Nginx 单个进程允许的客户端最大连接数。

```
events {
  worker_connections 20480;
}
```

`worker_connections` 也是个事件模块指令，用于定义 Nginx 每个进程的最大连接数，默认是 1024。

最大连接数的计算公式是：`max_clients = worker_processes * worker_connections`

另外，**进程的最大连接数受 Linux 系统进程的最大打开文件数限制**，在执行操作系统命令 `ulimit -HSn 65535`或配置相应文件后， `worker_connections` 的设置才能生效。

### 配置获取更多连接数

默认情况下，Nginx 进程只会在一个时刻接收一个新的连接，我们可以配置`multi_accept` 为 `on`，实现在一个时刻内可以接收多个新的连接，提高处理效率。该参数默认是 `off`，建议开启。

```
events {
  multi_accept on;
}
```

参考：http://nginx.org/en/docs/ngx_core_module.html#multi_accept



## Nginx 事件处理模型优化

Nginx 的连接处理机制在不同的操作系统中会采用不同的 I/O 模型，在 linux 下，Nginx 使用 epoll 的 I/O 多路复用模型，在 Freebsd 中使用 kqueue 的 I/O 多路复用模型，在 Solaris 中使用 /dev/poll 方式的 I/O 多路复用模型，在 Windows 中使用 icop，等等。

配置如下：

```
events {
  use epoll;
}
```

`events` 指令是设定 Nginx 的工作模式及连接数上限。`use`指令用来指定 Nginx 的工作模式。Nginx 支持的工作模式有 select、 poll、 kqueue、 epoll 、 rtsig 和/ dev/poll。当然，也可以不指定事件处理模型，Nginx 会自动选择最佳的事件处理模型。


## TCP优化
```
http {
  sendfile on;
  tcp_nopush on;

  keepalive_timeout 120;
  tcp_nodelay on;
}
```
第一行的 `sendfile` 配置可以提高 Nginx 静态资源托管效率。sendfile 是一个系统调用，直接在内核空间完成文件发送，不需要先 read 再 write，没有上下文切换开销。

TCP_NOPUSH 是 FreeBSD 的一个 socket 选项，对应 Linux 的 TCP_CORK，Nginx 里统一用 `tcp_nopush` 来控制它，并且只有在启用了 `sendfile` 之后才生效。启用它之后，数据包会累计到一定大小之后才会发送，减小了额外开销，提高网络效率。

TCP_NODELAY 也是一个 socket 选项，启用后会禁用 Nagle 算法，尽快发送数据，某些情况下可以节约 200ms（Nagle 算法原理是：在发出去的数据还未被确认之前，新生成的小数据先存起来，凑满一个 MSS 或者等到收到确认后再发送）。Nginx 只会针对处于 keep-alive 状态的 TCP 连接才会启用 `tcp_nodelay`。

## 前端性能优化
Nginx 开启 gzip 压缩之后可以降低网络IO，但是对前端来说，可以提前进行 gzip 压缩，这样请求的时候就不用再压缩了，减少对 cpu 的损耗。

Nginx 给你返回静态文件的时候，会判断是否开启gzip，然后压缩后再还给浏览器。

但是nginx其实会先判断是否有.gz后缀的相同文件，有的话直接返回，nginx自己不再进行压缩处理。

压缩是要时间的！不同级别的压缩率花的时间也不一样。所以提前准备gz文件，可以更加优化。而且你可以把压缩率提高点，这样带宽消耗会更小！！！

```
server {
    gzip on;
    gzip_vary on;
    gzip_min_length 1k;
    gzip_buffers 4 16k;
    gzip_comp_level 6;
    gzip_disable "MSIE [1-6]\."; #配置禁用gzip条件，支持正则。此处表示ie6及以下不启用gzip（因为ie低版本不支持）
    gzip_types  text/plain application/json application/javascript application/x-javascript text/css application/xml text/javascript application/x-httpd-php image/gif image/x-icon;
}
```

> 不适合压缩的数据：
> 
> 二进制资源：例如图片/mp3这样的二进制文件,不必压缩；因为压缩率比较小, 比如100->80字节,而且压缩也是耗费CPU资源的.


# Nginx 安全



## 隐藏 Nginx 版本号

添加配置：

```
server_tokens off;
```



## Nginx 防止 DDos 和 cc 等流量攻击

#### 限制同一时间段ip访问次数

nginx可以通过`ngx_http_limit_conn_module`和`ngx_http_limit_req_module`配置来限制ip在同一时间段的访问次数.

**ngx_http_limit_conn_module**：该模块用于限制每个定义的密钥的连接数，特别是单个IP地址的连接数．使用limit_conn_zone和limit_conn指令．

**ngx_http_limit_req_module**：用于限制每一个定义的密钥的请求的处理速率，特别是从一个单一的IP地址的请求的处理速率。使用“泄漏桶”方法进行限制．指令：limit_req_zone和limit_req．



ngx_http_limit_conn_module：限制单个IP的连接数示例：

```nginx
http { 
  limit_conn_zone $binary_remote_addr zone=addr：10m; 
　　 #定义一个名为addr的limit_req_zone用来存储session，大小是10M内存，
  #以$binary_remote_addr 为key,
  #nginx 1.18以后用limit_conn_zone替换了limit_conn,
  #且只能放在http{}代码段．
  ... 
  server { 
    ... 
    location /download/ { 
      limit_conn addr 1; 　　#连接数限制
      #设置给定键值的共享内存区域和允许的最大连接数。超出此限制时，服务器将返回503（服务临时不可用）错误.
　　　　　　　＃如果区域存储空间不足，服务器将返回503（服务临时不可用）错误
    }

```

可能有几个limit_conn指令,以下配置将限制每个客户端IP与服务器的连接数，同时限制与虚拟服务器的总连接数：

```nginx
http { 
  limit_conn_zone $binary_remote_addr zone=perip：10m; 
  limit_conn_zone $server_name zone=perserver：10m 
  ... 
  server { 
    ... 
    limit_conn perip 10; 　　　　 #单个客户端ip与服务器的连接数．
    limit_conn perserver 100;　　＃限制与服务器的总连接数
    }
```



**ngx_http_limit_req_module：限制某一时间内，单一IP的请求数**．

```nginx
http {
  limit_req_zone $binary_remote_addr zone=one:10m rate=1r/s;
  ...
　　#定义一个名为one的limit_req_zone用来存储session，大小是10M内存，　　
　　#以$binary_remote_addr 为key,限制平均每秒的请求为1个，
　　#1M能存储16000个状态，rete的值必须为整数，
　　
  server {
    ...
    location /search/ {
      limit_req zone=one burst=5;
　　　　　　　　
　　　　　　　　#限制每ip每秒不超过1个请求，漏桶数burst为5,也就是队列．
　　　　　　　　#nodelay，如果不设置该选项，严格使用平均速率限制请求数，超过的请求被延时处理．
　　　　　　　　#举个栗子：
　　　　　　　　＃设置rate=20r/s每秒请求数为２０个，漏桶数burst为5个，
　　　　　　　　#brust的意思就是，如果第1秒、2,3,4秒请求为19个，第5秒的请求为25个是被允许的，可以理解为20+5
　　　　　　　　#但是如果你第1秒就25个请求，第2秒超过20的请求返回503错误．
　　　　　　　　＃如果区域存储空间不足，服务器将返回503（服务临时不可用）错误　
　　　　　　　　＃速率在每秒请求中指定（r/s）。如果需要每秒少于一个请求的速率，则以每分钟的请求（r/m）指定。　
　　　　　　　　
    }
```

还可以限制来自单个IP地址的请求的处理速率，同时限制虚拟服务器的请求处理速率：

```nginx
http {
  limit_req_zone $binary_remote_addr zone=perip:10m rate=1r/s;
  limit_req_zone $server_name zone=perserver:10m rate=10r/s;
  ...
  server {
    ...
      limit_req zone=perip burst=5 nodelay;　　#漏桶数为５个．也就是队列数．nodelay:不启用延迟．
      limit_req zone=perserver burst=10;　　　　#限制nginx的处理速率为每秒10个
    }
```



#### 禁止ip或ip网段

查找服务器所有访问者ip方法:

```bash
$ awk '{print $1}' nginx_access.log |sort |uniq -c|sort -n
```

nginx.access.log 为nginx访问日志文件所在路径

会到如下结果，前面是ip的访问次数，后面是ip，很明显我们需要把访问次数多的ip并且不是蜘蛛的ip屏蔽掉，如下面结果， 
若 66.249.79.84 不为蜘蛛则需要屏蔽：

```
     89 106.75.133.167
     90 118.123.114.57
     91 101.78.0.210
     92 116.113.124.59
     92 119.90.24.73
     92 124.119.87.204
    119 173.242.117.145
   4320 66.249.79.84
```



#### 屏蔽IP的方法

在nginx的安装目录下面,新建屏蔽ip文件，命名为guolv_ip.conf，以后新增加屏蔽ip只需编辑这个文件即可。 加入如下内容并保存：

```
deny 66.249.79.84 ; 
```

在nginx的配置文件nginx.conf中加入如下配置，可以放到http, server, location, limit_except语句块，需要注意相对路径，本例当中nginx.conf，guolv_ip.conf在同一个目录中。

```
include guolv_ip.conf; 
```

屏蔽ip的配置文件既可以屏蔽单个ip，也可以屏蔽ip段，或者只允许某个ip或者某个ip段访问。

```c
//屏蔽单个ip访问

deny IP; 

//允许单个ip访问

allow IP; 

//屏蔽所有ip访问

deny all; 

//允许所有ip访问

allow all; 

//屏蔽整个段即从123.0.0.1到123.255.255.254访问的命令

deny 123.0.0.0/8

//屏蔽IP段即从123.45.0.1到123.45.255.254访问的命令

deny 124.45.0.0/16

//屏蔽IP段即从123.45.6.1到123.45.6.254访问的命令

deny 123.45.6.0/24

//如果你想实现这样的应用，除了几个IP外，其他全部拒绝，
//那需要你在guolv_ip.conf中这样写

allow 1.1.1.1; 
allow 1.1.1.2;
deny all; 
```


# Lua

- [lua-resty-http](https://github.com/ledgetech/lua-resty-http)

# 手动容器编译openresty支持slice模块

[Dockerfile](./Dockerfile)

参考文章：
- [Dockerfile 编写openresty 支持slice模块切片回源](http://www.xixicool.com/884.html)
- [Centos7编译openssl 3.0.x](https://blog.csdn.net/qq_42020376/article/details/143949796)

