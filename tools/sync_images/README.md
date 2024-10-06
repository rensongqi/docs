# sync_images

自动同步公网docker镜像，包含gcr.io、quay.io、docker.io、registry.k8s.io镜像至内网gitlab镜像仓库

依赖：[skopeo](https://github.com/containers/skopeo)

## 使用

通过在gitlab中提交issue来自动同步公网镜像至内网镜像仓库中

> 使用指定的 tag 用于同步。

```
k8s.gcr.io/pause
k8s.gcr.io/defaultbackend-amd64
```

同步规则

```
k8s.gcr.io/{image_name}  ==>  harbor.rsq.cn/gcr.io/{image_name}
```

**拉取镜像**

```bash
$ docker pull harbor.rsq.cn/gcr.io/<image_name>:[镜像版本号]
```

## 文件介绍

- `config.yaml`: 供 `generate_sync_yaml.py` 脚本使用，此文件配置了需要动态(获取`last`个最新的版本)同步的镜像列表。
- `custom_sync.yaml`: 自定义的 [`skopeo`](https://github.com/containers/skopeo) 同步源配置文件。
- `generate_sync_yaml.py`: 根据配置，动态生成 [`skopeo`](https://github.com/containers/skopeo) 同步源配置文件。定时同步时才需要用到此脚本。
- `run.py`: 用于手动或定时执行同步镜像操作。



## 手动同步镜像

将 `nvcr.io/nvidia/cuda:12.5.1-devel-ubuntu22.04` 同步到内网镜像仓库 `harbor.rsq.cn/nvcr.io/nvidia`

需要将`nvcr.io`更换为代理地址`ngc.nju.edu.cn`
```bash
skopeo --insecure-policy sync -a --keep-going --src docker --dest docker ngc.nju.edu.cn/nvidia/cuda:12.5.1-devel-ubuntu22.04 harbor.rsq.cn/nvcr.io/nvidia
```