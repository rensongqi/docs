
# 密码管理工具

推荐本地安装

下载地址：[Passbolt Install](https://www.passbolt.com/docs/hosting/install/)

centos 安装
```bash
sha512sum -c passbolt-ce-SHA512SUM.txt && sudo bash ./passbolt-repo-setup.ce.sh  || echo \"Bad checksum. Aborting\" && rm -f passbolt-repo-setup.ce.sh

sudo yum install passbolt-ce-server
```