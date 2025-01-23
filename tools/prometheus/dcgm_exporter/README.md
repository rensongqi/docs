
- [Install nvidia-docker](#install-nvidia-docker)
  - [Centos](#centos)
  - [Ubuntu](#ubuntu)
- [Run dcgm\_exporter](#run-dcgm_exporter)

# Install nvidia-docker

## Centos
```bash
# 配置repo
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.repo | \
  sudo tee /etc/yum.repos.d/nvidia-docker.repo

# import repo key
DIST=$(sed -n 's/releasever=//p' /etc/yum.conf)
DIST=${DIST:-$(. /etc/os-release; echo $VERSION_ID)}
sudo yum makecache

# 安装nvidia-docker2
sudo yum install nvidia-docker2

# 修改daemon.json ，把cgroupdriver改为systemd
sed -i 's/cgroupfs/systemd/g' /etc/docker/daemon.json

# daemon reload && 重启docker
systemctl daemon-reload
systemctl restart docker

# 进行测试
sudo docker run --rm --gpus all nvidia/cuda:11.0.3-base-ubuntu18.04 nvidia-smi
```

## Ubuntu

```bash
# Ubuntu 配置repo(没有key会报)
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | \sudo tee /etc/apt/sources.list.d/nvidia-docker.list

# Ubuntu import repo key
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | \sudo apt-key add -
sudo apt update

# Ubuntu 安装sudo apt-get install nvidia-docker2
sudo apt install nvidia-docker2

# 修改daemon.json ，把cgroupdriver改为systemd
sed -i 's/cgroupfs/systemd/g' /etc/docker/daemon.json

# daemon reload && 重启docker
systemctl daemon-reload
systemctl restart docker

# 进行测试
sudo docker run --rm --gpus all nvidia/cuda:11.0.3-base-ubuntu18.04 nvidia-smi

```


# Run dcgm_exporter
```bash
nvidia-docker run --cap-add SYS_ADMIN --name 249dcgm-exporter --restart=unless-stopped -d -p 9400:9400 nvidia/dcgm-exporter:2.2.9-2.4.0-ubuntu18.04
```