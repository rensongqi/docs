
# nfpm是一款构建deb或rpm安装包工具

# 使用nfpm打包二进制程序为deb包
> 官方文档: https://nfpm.goreleaser.com/usage/

# 下载nfpm

1. go install
```bash
go install github.com/goreleaser/nfpm/v2/cmd/nfpm@latest
```

2. apt install
```bash
echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | sudo tee /etc/apt/sources.list.d/goreleaser.list
sudo apt update
sudo apt install nfpm
```

3. manually install
```bash
https://github.com/goreleaser/nfpm/releases
```

# 详细配置
```bash
mkdir /data/custom_process
cd /data/custom_process 

# 可以使用如下命令初始化一个打包deb的配置清单
nfpm init
```

nfpm.yaml
> 如下配置文件会将指定的二进制程序  custom_process 拷贝至 /usr/local/bin/ 且文件权限会被改成755

```yaml
# nfpm example configuration file
#
# check https://nfpm.goreleaser.com/configuration for detailed usage
#
name: "custom_process"
arch: "amd64"
platform: "linux"
version: "1.0.0"
section: "default"
priority: "extra"
maintainer: "Rensongqi <rensongqi@rsq.com>"
description: |
  custom_process linux and64 binary
contents:
- src: custom_process
  dst: /usr/local/bin/custom_process
  file_info:
    mode: 0755
```

# 构建deb包
```bash
# 在目录/data/custom_process下执行如下命令， build deb包，将生成的deb包放到/tmp/目录下
nfpm pkg --packager deb --target /tmp/

# 如果是rpm包，则执行如下命令
nfpm pkg --packager rpm --target /tmp/
```

# 使用自定义源

参考此前的文章 [ubuntu private repo](../linux/ubuntu_private_repo.md)