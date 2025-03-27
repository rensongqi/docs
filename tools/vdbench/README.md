

- [vdbench](#vdbench)
- [1 多节点测试](#1-多节点测试)
  - [1.1 顺序测试](#11-顺序测试)
  - [1.2 随机读写](#12-随机读写)
  - [1.3 混合读写](#13-混合读写)
- [2 单节点测试](#2-单节点测试)
  - [2.1 顺序测试](#21-顺序测试)
  - [2.2 随机读写](#22-随机读写)
  - [2.3 混合读写](#23-混合读写)

# vdbench

vdbench是一个I/O工作负载生成器，通常用于验证数据完整性和度量直接附加（或网络连接）存储性能。它可以运行在windows、linux环境，可用于测试文件系统或块设备基准性能。

参考：

- [vdbench存储性能测试工具](https://www.cnblogs.com/luxf0/p/13321077.html)
- [vdbench下载地址](http://download.oracle.com/otn/utilities_drivers/vdbench/vdbench50406.zip)
- [windows jdk](https://download.oracle.com/otn/java/jdk/8u251-b08/3d5a2bb8f8d4428bbe94aed7ec7ae784/jdk-8u251-windows-x64.exe?AuthParam=1590741958_7a95bcd255b6fd15a5a1217c7317fe14)
- [linux jdk](https://download.oracle.com/otn/java/jdk/8u251-b08/3d5a2bb8f8d4428bbe94aed7ec7ae784/jdk-8u251-linux-x64.tar.gz?AuthParam=1590742346_3243bad9a4f1a3147c170be345a4f972)

准备
1. 部署jdk
2. 配置免密
3. 配置/etc/hosts

# 1 多节点测试
脚本只需要在一台机器上存放即可，每台机器vdbench对应的目录`/home/vdbench`下需要有vdbench的可执行程序。测试机器master1~node01四台机器之间需要做互相免密认证，否则不能同时配置多台机器同时压测

## 1.1 顺序测试
`vim sequntial.file`
```
messagescan=no
hd=default,vdbench=/home/vdbench,user=root,shell=ssh
hd=hd1,system=master1
hd=hd2,system=master2
hd=hd3,system=master3
hd=hd4,system=node01

fsd=fsd1,anchor=/gpfs/MaGW01,depth=1,width=1,shared=yes
fsd=fsd2,anchor=/gpfs/MaGW02,depth=1,width=1,shared=yes
fsd=fsd3,anchor=/gpfs/MaGW03,depth=1,width=1,shared=yes
fsd=fsd4,anchor=/gpfs/MaGW04,depth=1,width=1,shared=yes

fwd=fwd1,fsd=fsd1,hd=hd1,fileio=sequntial,fileselect=sequntial
fwd=fwd2,fsd=fsd2,hd=hd2,fileio=sequntial,fileselect=sequntial
fwd=fwd3,fsd=fsd3,hd=hd3,fileio=sequntial,fileselect=sequntial
fwd=fwd4,fsd=fsd4,hd=hd4,fileio=sequntial,fileselect=sequntial

rd=rd1,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=5120,size=1g,warmup=6,elapsed=666,interval=1,operation=write,threads=8
rd=rd2,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=5120,size=1g,warmup=6,elapsed=666,interval=1,operation=read,threads=8
rd=rd3,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=5120,size=2g,warmup=6,elapsed=666,interval=1,operation=write,threads=8
rd=rd4,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=5120,size=2g,warmup=6,elapsed=666,interval=1,operation=read,threads=8
rd=rd5,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=5120,size=1g,warmup=6,elapsed=666,interval=1,operation=write,threads=8
rd=rd6,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=5120,size=1g,warmup=6,elapsed=666,interval=1,operation=read,threads=8
rd=rd7,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=5120,size=2g,warmup=6,elapsed=666,interval=1,operation=write,threads=8
rd=rd8,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=5120,size=2g,warmup=6,elapsed=666,interval=1,operation=read,threads=8
```

执行测试程序

```bash
./vdbench -f sequntial.file

# 默认执行完会在当前目录下生成output文件，查看平均信息
grep avg_7 output/summary.html
```

## 1.2 随机读写

```
messagescan=no
hd=default,vdbench=/home/vdbench,user=root,shell=ssh
hd=hd1,system=master1
hd=hd2,system=master2
hd=hd3,system=master3
hd=hd4,system=node01

fsd=fsd1,anchor=/gpfs/random/MaGW01,depth=1,width=1,shared=yes
fsd=fsd2,anchor=/gpfs/random/MaGW02,depth=1,width=1,shared=yes
fsd=fsd3,anchor=/gpfs/random/MaGW03,depth=1,width=1,shared=yes
fsd=fsd4,anchor=/gpfs/random/MaGW04,depth=1,width=1,shared=yes

fwd=fwd1,fsd=fsd1,hd=hd1,fileio=random,fileselect=random
fwd=fwd2,fsd=fsd2,hd=hd2,fileio=random,fileselect=random
fwd=fwd3,fsd=fsd3,hd=hd3,fileio=random,fileselect=random
fwd=fwd4,fsd=fsd4,hd=hd4,fileio=random,fileselect=random

rd=rd5,fwd=(fwd*),fwdrate=max,format=restart,xfersize=4k,files=512000,size=4k,warmup=6,elapsed=666,interval=1,operation=write,threads=32
rd=rd1,fwd=(fwd*),fwdrate=max,format=restart,xfersize=4k,files=512000,size=4k,warmup=6,elapsed=666,interval=1,operation=read,threads=32
rd=rd6,fwd=(fwd*),fwdrate=max,format=restart,xfersize=64k,files=512000,size=64k,warmup=6,elapsed=666,interval=1,operation=write,threads=32
rd=rd2,fwd=(fwd*),fwdrate=max,format=restart,xfersize=64k,files=512000,size=64k,warmup=6,elapsed=666,interval=1,operation=read,threads=32
rd=rd7,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=512000,size=1m,warmup=6,elapsed=666,interval=1,operation=write,threads=32
rd=rd3,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=512000,size=1m,warmup=6,elapsed=666,interval=1,operation=read,threads=32
rd=rd8,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=512000,size=2m,warmup=6,elapsed=666,interval=1,operation=write,threads=32
rd=rd4,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=512000,size=2m,warmup=6,elapsed=666,interval=1,operation=read,threads=32
```

## 1.3 混合读写

```
messagescan=no
hd=default,vdbench=/home/vdbench,user=root,shell=ssh
hd=hd1,system=master1
hd=hd2,system=master2
hd=hd3,system=master3
hd=hd4,system=node01

fsd=fsd1,anchor=/gpfs/mix/MaGW01,depth=1,width=1,shared=yes
fsd=fsd2,anchor=/gpfs/mix/MaGW02,depth=1,width=1,shared=yes
fsd=fsd3,anchor=/gpfs/mix/MaGW03,depth=1,width=1,shared=yes
fsd=fsd4,anchor=/gpfs/mix/MaGW04,depth=1,width=1,shared=yes

fwd=fwd1,fsd=fsd1,hd=hd1,fileio=random,fileselect=random
fwd=fwd2,fsd=fsd2,hd=hd2,fileio=random,fileselect=random
fwd=fwd3,fsd=fsd3,hd=hd3,fileio=random,fileselect=random
fwd=fwd4,fsd=fsd4,hd=hd4,fileio=random,fileselect=random

rd=rd1,fwd=(fwd*),fwdrate=max,format=restart,xfersize=4k,files=512000,size=4k,warmup=6,elapsed=666,interval=1,operation=read,rdpct=70,threads=32
rd=rd2,fwd=(fwd*),fwdrate=max,format=restart,xfersize=64k,files=512000,size=64k,warmup=6,elapsed=666,interval=1,operation=read,rdpct=70,threads=32
rd=rd3,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=512000,size=1m,warmup=6,elapsed=666,interval=1,operation=read,rdpct=70,threads=32
rd=rd4,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=512000,size=2m,warmup=6,elapsed=666,interval=1,operation=read,rdpct=70,threads=32
```

# 2 单节点测试
- rd中files数量 * width = 生成文件总数
- depth: 目录深度
- 直通模式：openflags=o_direct

## 2.1 顺序测试

```
messagescan=no
hd=default,vdbench=/home/vdbench,user=root
hd=hd1,system=train-4090x4-amd-1-246

fsd=fsd1,anchor=/disk/deepdata/test_io/sequntial,depth=1,width=1,shared=yes

fwd=fwd1,fsd=fsd1,hd=hd1,fileio=sequntial,fileselect=sequntial

rd=rd1,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=512,size=1g,warmup=6,elapsed=60,interval=1,operation=write,threads=8
rd=rd2,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=512,size=1g,warmup=6,elapsed=60,interval=1,operation=read,threads=8
rd=rd5,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=512,size=1g,warmup=6,elapsed=60,interval=1,operation=write,threads=8
rd=rd6,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=512,size=1g,warmup=6,elapsed=60,interval=1,operation=read,threads=8
rd=rd3,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=512,size=2g,warmup=6,elapsed=60,interval=1,operation=write,threads=8
rd=rd4,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=512,size=2g,warmup=6,elapsed=60,interval=1,operation=read,threads=8
rd=rd7,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=512,size=2g,warmup=6,elapsed=60,interval=1,operation=write,threads=8
rd=rd8,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=512,size=2g,warmup=6,elapsed=60,interval=1,operation=read,threads=8
```

## 2.2 随机读写

```
messagescan=no
hd=default,vdbench=/home/vdbench,user=root
hd=hd1,system=train-4090x4-amd-1-246

fsd=fsd1,anchor=/disk/deepdata/test_io/random,depth=1,width=1,shared=yes
fwd=fwd1,fsd=fsd1,hd=hd1,fileio=random,fileselect=random

rd=rd5,fwd=(fwd*),fwdrate=max,format=restart,xfersize=4k,files=51200,size=4k,warmup=6,elapsed=60,interval=1,operation=write,threads=32
rd=rd1,fwd=(fwd*),fwdrate=max,format=restart,xfersize=4k,files=51200,size=4k,warmup=6,elapsed=60,interval=1,operation=read,threads=32
rd=rd6,fwd=(fwd*),fwdrate=max,format=restart,xfersize=64k,files=51200,size=64k,warmup=6,elapsed=60,interval=1,operation=write,threads=32
rd=rd2,fwd=(fwd*),fwdrate=max,format=restart,xfersize=64k,files=51200,size=64k,warmup=6,elapsed=60,interval=1,operation=read,threads=32
rd=rd7,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=51200,size=1m,warmup=6,elapsed=60,interval=1,operation=write,threads=32
rd=rd3,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=51200,size=1m,warmup=6,elapsed=60,interval=1,operation=read,threads=32
rd=rd8,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=51200,size=2m,warmup=6,elapsed=60,interval=1,operation=write,threads=32
rd=rd4,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=51200,size=2m,warmup=6,elapsed=60,interval=1,operation=read,threads=32
```

## 2.3 混合读写

```
messagescan=no
hd=default,vdbench=/home/vdbench,user=root
hd=hd1,system=train-4090x4-amd-1-246

# 带缓存
# fsd=fsd1,anchor=/disk/deepdata/test_io/mix,depth=32,width=1,shared=yes
# 不使用缓存
fsd=fsd1,anchor=/disk/deepdata/test_io/mix,depth=32,width=1,shared=yes,openflags=o_direct

fwd=fwd1,fsd=fsd1,hd=hd1,fileio=random,fileselect=random

rd=rd1,fwd=(fwd*),fwdrate=max,format=restart,xfersize=4k,files=51200,size=4k,warmup=6,elapsed=60,interval=1,operation=read,rdpct=70,threads=32
rd=rd2,fwd=(fwd*),fwdrate=max,format=restart,xfersize=64k,files=51200,size=64k,warmup=6,elapsed=60,interval=1,operation=read,rdpct=70,threads=32
rd=rd3,fwd=(fwd*),fwdrate=max,format=restart,xfersize=1m,files=51200,size=1m,warmup=6,elapsed=60,interval=1,operation=read,rdpct=70,threads=32
rd=rd4,fwd=(fwd*),fwdrate=max,format=restart,xfersize=2m,files=51200,size=2m,warmup=6,elapsed=60,interval=1,operation=read,rdpct=70,threads=32
```

