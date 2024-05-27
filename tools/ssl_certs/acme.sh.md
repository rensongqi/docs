# 使用acme.sh生成公网证书

## 安装acme.sh工具
```bash
git clone https://gitee.com/neilpang/acme.sh.git
cd acme.sh
./acme.sh --install -m songqi.ren@rensongqi.com
```

## 安装socat
```bash
# Centos
yum -y install socat

# Ubuntu
apt -y install socat
```

## 配置CA
```bash
./acme.sh --set-default-ca --server letsencrypt
```

## 使用dns的方式签发证书
```bash
# step1: 生成TXT
# ./acme.sh --issue -d openai.rensongqi.com --dns --yes-I-know-dns-manual-mode-enough-go-ahead-please
[Wed May 22 10:20:45 CST 2024] Using CA: https://acme-v02.api.letsencrypt.org/directory
[Wed May 22 10:20:45 CST 2024] Creating domain key
[Wed May 22 10:20:45 CST 2024] The domain key is here: /root/.acme.sh/openai.rensongqi.com/openai.rensongqi.com.key
[Wed May 22 10:20:45 CST 2024] Single domain='openai.rensongqi.com'
[Wed May 22 10:20:45 CST 2024] Getting domain auth token for each domain
[Wed May 22 10:20:49 CST 2024] Getting webroot for domain='openai.rensongqi.com'
[Wed May 22 10:20:49 CST 2024] Add the following TXT record:
[Wed May 22 10:20:49 CST 2024] Domain: '_acme-challenge.openai.rensongqi.com'
[Wed May 22 10:20:49 CST 2024] TXT value: 'axnHAeLnW8e1wxHVYS2PJqiL8lT0aEYk-G3vIN4CCbE'
[Wed May 22 10:20:49 CST 2024] Please be aware that you prepend _acme-challenge. before your domain
[Wed May 22 10:20:49 CST 2024] so the resulting subdomain will be: _acme-challenge.openai.rensongqi.com
[Wed May 22 10:20:49 CST 2024] Please add the TXT records to the domains, and re-run with --renew.
[Wed May 22 10:20:49 CST 2024] Please check log file for more details: /root/.acme.sh/acme.sh.log

# step2: Add Txt record
# 将第一步中的黄色标识的TXT添加至阿里云DNS解析中

# step3: renew cert
# ./acme.sh --renew -d openai.rensongqi.com --dns --yes-I-know-dns-manual-mode-enough-go-ahead-please
[Wed May 22 11:18:18 CST 2024] Renew: 'openai.rensongqi.com'
[Wed May 22 11:18:18 CST 2024] Renew to Le_API=https://acme-v02.api.letsencrypt.org/directory
[Wed May 22 11:18:22 CST 2024] Using CA: https://acme-v02.api.letsencrypt.org/directory
[Wed May 22 11:18:22 CST 2024] Single domain='openai.rensongqi.com'
[Wed May 22 11:18:22 CST 2024] Getting domain auth token for each domain
[Wed May 22 11:18:22 CST 2024] Verifying: openai.rensongqi.com
[Wed May 22 11:18:27 CST 2024] Pending, The CA is processing your order, please just wait. (1/30)
[Wed May 22 11:18:31 CST 2024] Success
[Wed May 22 11:18:31 CST 2024] Verify finished, start to sign.
[Wed May 22 11:18:31 CST 2024] Lets finalize the order.
[Wed May 22 11:18:31 CST 2024] Le_OrderFinalize='https://acme-v02.api.letsencrypt.org/acme/finalize/1739317072/27146xxx'
[Wed May 22 11:18:33 CST 2024] Downloading cert.
[Wed May 22 11:18:33 CST 2024] Le_LinkCert='https://acme-v02.api.letsencrypt.org/acme/cert/031f405a7e0a59xxx'
[Wed May 22 11:18:36 CST 2024] Cert success.
-----BEGIN CERTIFICATE-----
MIIE9jCCA96gAwIBAgISAx9AWn4KWRUevkzSqYAU6EfVMA0GCSqGSIb3DQEBCwUA
MDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD
EwJSMzAeFw0yNDA1MjIwMjE4MzJaFw0yNDA4MjAwMjE4MzFaMB8xHTAbBgNVBAMT
FG9wZW5haS5jb3dhcm9ib3QuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAvtYfa2Eye6mF8LFSvdLPz6QSR087Qufo6aKy7wG1MwCx4X87pZGl6Ex4
kvKF7sbJ6jBErJe6IfrRmLVRco+3XdI/2Lm4R9KVBDXL59AqlPYcEzou/o+lpg0l
c9xCPgdxng6J5g5w2IH0vHxo/pKq9a4eLKF4bY09VO/Qt4gqok0TcnW07778Vbdo
y/u8JXAx0FhaiBVzNthr9cGJXeurme9ToRc2So0lSG73KgNY5MQ9GqNyJ6HPdhwu
3BHdtKf3oYxgsNSBe5Fh7Bag67HwXiC9lJkgOYT/n11vq1scWgG/tf25TwKlLLnh
xxxxxxxxxx
SV25qhsmZIMw6t77KB1YtajG48bwwJHDtBg=
-----END CERTIFICATE-----
[Wed May 22 11:18:36 CST 2024] Your cert is in: /root/.acme.sh/openai.rensongqi.com/openai.rensongqi.com.cer
[Wed May 22 11:18:36 CST 2024] Your cert key is in: /root/.acme.sh/openai.rensongqi.com/openai.rensongqi.com.key
[Wed May 22 11:18:36 CST 2024] The intermediate CA cert is in: /root/.acme.sh/openai.rensongqi.com/ca.cer
[Wed May 22 11:18:36 CST 2024] And the full chain certs is there: /root/.acme.sh/openai.rensongqi.com/fullchain.cer
```

## 安装证书
```bash
# ./acme.sh --install-cert -d openai.rensongqi.com --key-file /tmp/nginx/certs/openai.rensongqi.com.key --fullchain-file /tmp/nginx/certs/openai.rensongqi.com.pem --reloadcmd "systemctl reload openresty"
[Wed May 22 11:24:38 CST 2024] Installing key to: /tmp/nginx/certs/openai.rensongqi.com.key
[Wed May 22 11:24:38 CST 2024] Installing full chain to: /tmp/nginx/certs/openai.rensongqi.com.pem
```

## 查看所有证书
```bash
# ./acme.sh --list
Main_Domain                    KeyLength  SAN_Domains  CA               Created               Renew
openai.rensongqi.com           "2048"     no           LetsEncrypt.org  2024-05-22T03:18:36Z  2024-07-21T03:18:36Z
```

## 删除证书
```bash
./acme.sh --remove -d openai.rensongqi.com
```

## 参考文档：
https://github.com/acmesh-official/acme.sh/wiki/dns-manual-mode