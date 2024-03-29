搭建ubuntu私有repo，使用gpg对私有包进行自签认证，通过apt安装签名后的包

# 1 GPG配置

```bash
# 1 生成gpg密钥，这里没有用密码认证
root@RSQ2204:~# gpg --full-generate-key
gpg (GnuPG) 2.2.27; Copyright (C) 2021 Free Software Foundation, Inc.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Please select what kind of key you want:
   (1) RSA and RSA (default)
   (2) DSA and Elgamal
   (3) DSA (sign only)
   (4) RSA (sign only)
  (14) Existing key from card
Your selection? 1
RSA keys may be between 1024 and 4096 bits long.
What keysize do you want? (3072) 4096
Requested keysize is 4096 bits
Please specify how long the key should be valid.
         0 = key does not expire
      <n>  = key expires in n days
      <n>w = key expires in n weeks
      <n>m = key expires in n months
      <n>y = key expires in n years
Key is valid for? (0) 0
Key does not expire at all
Is this correct? (y/N) y

GnuPG needs to construct a user ID to identify your key.

Real name: rensongqi
Email address: admin@rsq.com
Comment: rensongqi
You selected this USER-ID:
    "rensongqi (rensongqi) <admin@rsq.com>"

Change (N)ame, (C)omment, (E)mail or (O)kay/(Q)uit? O
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.
gpg: /root/.gnupg/trustdb.gpg: trustdb created
gpg: key 29BDA33B3FF40E51 marked as ultimately trusted
gpg: directory '/root/.gnupg/openpgp-revocs.d' created
gpg: revocation certificate stored as '/root/.gnupg/openpgp-revocs.d/3EFC047A4FB62F5DD1324E3629BDA33B3FF40E51.rev'
public and secret key created and signed.

pub   rsa4096 2024-03-29 [SC]
      3EFC047A4FB62F5DD1324E3629BDA33B3FF40E51
uid                      rensongqi (rensongqi) <admin@rsq.com>
sub   rsa4096 2024-03-29 [E]

# 2 先获取当前key有哪些，使用指定的key生成Release.gpg文件
root@RSQ2204:~# gpg --list-secret-keys --keyid-format=long
gpg: checking the trustdb
gpg: marginals needed: 3  completes needed: 1  trust model: pgp
gpg: depth: 0  valid:   1  signed:   0  trust: 0-, 0q, 0n, 0m, 0f, 1u
/root/.gnupg/pubring.kbx
------------------------
sec   rsa4096/29BDA33B3FF40E51 2024-03-29 [SC]
      3EFC047A4FB62F5DD1324E3629BDA33B3FF40E51
uid                 [ultimate] rensongqi (rensongqi) <admin@rsq.com>
ssb   rsa4096/DF02B8F8B21963C9 2024-03-29 [E]

# 3 生成gpgkey文件（只需要生成一次就行）
gpg --export 29BDA33B3FF40E51 > gpgkey
```

配置内部环境

```bash
# 1 创建相对应的目录
cd /disk/storage/mirrors/custom
mkdir -p dists/{jammy,bionic,focal}/pool/stable/amd64/
mkdir -p dists/{jammy,bionic,focal}/stable/binary-amd64/

# 2 生成binary-amd64目录下的索引文件
# 生成Packages文件
cd /disk/storage/mirrors/custom/
apt-ftparchive packages dists/jammy > dists/jammy/stable/binary-amd64/Packages

# 3 生成jammy当前目录的签名文件
# patch.conf
APT::FTPArchive::Release {
  Origin "custom";
  Label "custom";
  Suite "jammy";
  Codename "jammy";
  Architectures "amd64";
  Components "stable";
  Description "jammy";
};

# 指定 patch.conf 生成Release文件
apt-ftparchive release -c=./patch.conf dists/jammy > dists/jammy/Release
# 生成InRelease文件
gpg --clearsign -o dists/jammy/InRelease dists/jammy/Release
# 生成gpg签名文件
gpg -abs -u 3EFC047A4FB62F5DD1324E3629BDA33B3FF40E51 -o dists/jammy/Release.gpg dists/jammy/Release
```


# 2 客户端使用

ubuntu20.04
```bash
echo -e "deb [arch=amd64] https://mirrors.rsq.cn/custom/ focal stable" >> /etc/apt/sources.list
curl -fsSL https://mirrors.rsq.cn/custom/gpgkey | sudo tee /etc/apt/trusted.gpg.d/custom.gpg >>/dev/null 2>&1
sudo chmod 644 /etc/apt/trusted.gpg.d/custom.gpg

# 更新源
apt update -y
```

ubuntu22.04
```bash
echo -e "deb [arch=amd64] https://mirrors.rsq.cn/custom/ jammy stable" >> /etc/apt/sources.list
curl -fsSL https://mirrors.rsq.cn/custom/gpgkey | sudo tee /etc/apt/trusted.gpg.d/custom.gpg >>/dev/null 2>&1
sudo chmod 644 /etc/apt/trusted.gpg.d/custom.gpg

# 更新源
apt update -y
```

# 3 刷新索引和签名
```bash
#!/bin/bash

# 生成jammy的索引和签名
# 1 生成Packages索引文件
cd /disk/storage/mirrors/custom/
apt-ftparchive packages dists/jammy > dists/jammy/stable/binary-amd64/Packages

# 2 生成Release文件
apt-ftparchive release -c=./patch.conf dists/jammy > dists/jammy/Release

# 3 生成InRelease文件
rm -f -r dists/jammy/InRelease
gpg --clearsign -o dists/jammy/InRelease dists/jammy/Release

# 生成gpg签名文件
rm -f -r dists/jammy/Release.gpg
gpg -abs -u 3EFC047A4FB62F5DD1324E3629BDA33B3FF40E51 -o dists/jammy/Release.gpg dists/jammy/Release


# 生成focal的索引和签名
# 1 生成Packages索引文件
cd /disk/storage/mirrors/custom/
apt-ftparchive packages dists/focal > dists/focal/stable/binary-amd64/Packages

# 2 生成Release文件
apt-ftparchive release -c=./patch.conf dists/focal > dists/focal/Release

# 3 生成InRelease文件
rm -f -r dists/focal/InRelease
gpg --clearsign -o dists/focal/InRelease dists/focal/Release

# 生成gpg签名文件
rm -f -r dists/focal/Release.gpg
gpg -abs -u 3EFC047A4FB62F5DD1324E3629BDA33B3FF40E51 -o dists/focal/Release.gpg dists/focal/Release
```
