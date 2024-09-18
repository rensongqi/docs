
# 制作镜像
```bash
wget https://github.com/tuna/tunasync/releases/download/v0.8.0/tunasync-linux-amd64-bin.tar.gz
tar -xf tunasync-linux-amd64-bin.tar.gz

# 解压之后会生成两个可执行文件
tunasync
tunasynctl

# build
docker build -t harbor.rsq.cn/library/tunasync:v3 --no-cache .
```

# 运行容器
```
docker-compose up -d
```