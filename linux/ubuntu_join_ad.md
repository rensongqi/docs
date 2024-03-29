参考如下脚本

```bash
#!/bin/bash

# 获取并验证质询密码
echo -e "$(date) 质询密码为16位字符串，只能使用一次，且将在 60 分钟内过期，获取地址：http://172.16.110.13/CertSrv/mscep_admin/" | tee -a /var/log/init_ad.log
challenge_password=""
if [ $# -eq 1 ]; then
    challenge_password="$1"
else
    echo "Please input challenge_password: "
    read -s challenge_password
fi

if [ ${#challenge_password} -ne 16 ]; then
    echo -e "$(date) 格式错误，质询密码为16位字符串！exit 1" | tee -a /var/log/init_ad.log
    exit 1
fi

# 配置质询密码
echo "${challenge_password}" > /tmp/challenge_password

#Change hostname
sn=$(cat /sys/class/dmi/id/product_serial)
#sn=$(dmidecode -s system-serial-number)
hostname="SHXHU-"$sn.ROBOT.CN
hostnamectl set-hostname ${hostname}
echo -e $(date)" new computer name is "$hostname >>/var/log/init_ad.log
sed -i "1c $hostname" /etc/hostname
sed -i "2c 127.0.0.1 $hostname" /etc/hosts
#hostname -f

#set user01 system account
if [ `grep -c "user01" /etc/passwd` -eq '0' ];then
    groupadd --gid 1001 user01
    useradd -m --uid 1001 --gid 1001 --shell /bin/bash user01
    # 兼容uid被占用的情况
    if [ $? -ne 0 ]; then
        useradd -m --shell /bin/bash user01
    fi
    # openssl passwd -1
    echo "user01:\$6\$qRP2.qc.74DO30zk\$tK.e7aoGKQJBNM2EUQL/5Gln0iNv9aaYz1/pEo0omwrKkuRvL3FM5t3zy8gzyQmEcrIgingxa8rFcW.77oitO1" |chpasswd -e
fi
#sed -i "s/false/true/g" /var/lib/AccountsService/users/user01
cat << EOF > /var/lib/AccountsService/users/user01
[User]
Session=
XSession=
Icon=/home/user01/.face
SystemAccount=true

[InputSource0]
xkb=us

[InputSource1]
ibus=libpinyin
EOF

#set GUI login interface
sed -i '$a greeter-show-manual-login=true \ngreeter-hide-users=true' /usr/share/lightdm/lightdm.conf.d/50-ubuntu.conf

#provide root privilege for domain users
if [ `grep -c "user01" /etc/sudoers` -eq '0' ];then
    echo 'user01               ALL=(ALL:ALL) ALL'>>/etc/sudoers
    echo '%Domain\ Users               ALL=(ALL:ALL) NOPASSWD: ALL'>>/etc/sudoers
fi
sed -i 's/admin/Domain Users/' /etc/polkit-1/localauthority.conf.d/51-ubuntu-admin.conf

##Wake on lan settings
apt-get install ethtool -y
#ethtool -s enp3s0 wol g
#echo -e '#! /bin/sh\nethtool -s enp3s0 wol g'>>/etc/rc.local
#echo 'NETDOWN = no'>>/etc/init.d/halt
cat << EOF > /etc/systemd/system/wol@.service
[Unit]
Description=Wake-on-LAN for %i
Requires=network.target
After=network.target

[Service]
ExecStart=/sbin/ethtool -s %i wol g
Type=oneshot

[Install]
WantedBy=multi-user.target
EOF
#systemctl enable wol@enp3s0
for i in $(ip -br l | awk '$1 !~ "lo|vir|wl" { print $1}');
do
	s_name="wol@"$i
	systemctl enable $s_name
done
##Join domain
apt-get install openssh-server xrdp sssd-ad sssd-tools realmd adcli -y

if [ `grep -c "172.16.100.11" /etc/systemd/timesyncd.conf` -eq '0' ];then
    echo -e "NTP=172.16.100.11\nFallbackNTP=ntp.ubuntu.com" >> /etc/systemd/timesyncd.conf
fi

systemctl restart systemd-timesyncd
cat << 'EOF' > /usr/bin/JoinAD.sh
#!/bin/bash
RETVAL=0
prog="JoinAD"
LOCKFILE=/var/lock/subsys/$prog
OPTION=$1
# Declare variables for script
DOMAIN=ROBOT.CN
AWSCONNECTORUSER="<aduser>"
AWSCONNECTORPASS=<adPass>
# Exit script if HOSTNAME is not set
if [ "$HOSTNAME" = "localhost" ]; then
    sn=$(dmidecode -s system-serial-number)
    hostname="SHXHU-"$sn
    echo -e $(date)" new computer name is "$hostname >>/var/log/init_ad.log
    sed -i "1c $hostname" /etc/hostname
    sed -i "2c 127.0.0.1 $hostname" /etc/hosts
    hostname -f
    #echo "System HOSTNAME is not set, cannot joing domain..exiting !"
    #exit 1
else
    echo "System Hostname: $HOSTNAME"
fi
start() {
        echo "Initiating $prog: "
        echo "Checking if already joined to domain..."
        ISJOINED=`realm list | grep -i $DOMAIN`
        RETVAL=$?
        if [ $RETVAL -gt 0 ] ;then
            echo  "Joining Domain $DOMAIN"
            echo $AWSCONNECTORPASS | realm join -U $AWSCONNECTORUSER  $DOMAIN
            ISJOINED=`realm list | grep -i $DOMAIN`
            RETVAL=$?
            if [ $RETVAL -eq 0 ]; then
                echo "Joined Domain: $DOMAIN"
                echo "Running Authconfig"
                pam-auth-update --enable mkhomedir
                if [ $? -ne 0 ]; then
                  echo "Authconfig failed to run sucessfully ..."
                fi
                # Fix SSSD.conf
                echo "Fixing /etc/sssd/sssd.conf"
                # Fix the realm Format
                sed -i 's/#\?\(use_fully_qualified_names =\s*\).*$/\1 False/' /etc/sssd/sssd.conf
                sed -i 's/#\?\(fallback_homedir =\s*\).*$/\1 \/home\/\%u/' /etc/sssd/sssd.conf
                echo "Restarting SSSD service"
                # Restart sssd daemon
                service sssd restart
                echo "Sucessfully Joined domain: $DOMAIN"
                echo "Use: realm list - to check status"
            else
                echo "Failed to Join Domain $DOMAIN"
            fi
        else
            echo "Already joined to domain $DOMAIN or "
            echo  "Skipping ..... "
        fi
        echo
        return $RETVAL
}
stop() {
        echo "Doing nothing .... if you want to leave domain"
        echo "Use: realm leave"
        RETVAL=$?
        echo
        return $RETVAL
}
status() {
        echo -n "Checking $prog status: "
        realm list
        echo
        RETVAL=$?
        return $RETVAL
}
case "$OPTION" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    restart)
        stop
        start
        ;;
    *)
        echo "Usage: $prog {start|stop|status|restart}"
        exit 1
        ;;
esac
exit $RETVAL
EOF

chmod +x /usr/bin/JoinAD.sh
cat << EOF > /etc/systemd/system/JoinAD.service
[Unit]
Description=Script to Join Active Directory
After=network-online.target
[Service]
Type=simple
ExecStart=/usr/bin/JoinAD.sh start
ExecStop=/usr/bin/JoinAD.sh stop
TimeoutStartSec=30
[Install]
WantedBy=default.target
EOF

# Enable ADJoin Service
systemctl daemon-reload
systemctl enable JoinAD
systemctl start JoinAD

##Disable ubuntu 20.04.1 software updates
if [ `grep -c "Hidden=true" /etc/xdg/autostart/update-notifier.desktop` -eq '0' ];then
  echo "Hidden=true" >> /etc/xdg/autostart/update-notifier.desktop
fi

if [ `grep -c "root        soft        nofile  655350" /etc/security/limits.conf` -eq '0' ];then
    cat >> /etc/security/limits.conf << EOF
*       soft        nofile  655350
*       hard        nofile  655350
*       soft        nproc   655350
*       hard        nproc   655350
root        soft        nofile  655350
root        hard        nofile  655350
root        soft        nproc   655350
root        hard        nproc   655350
EOF
fi

##修复vi命令乱码：
apt remove vim-common -y
apt-get install vim -y
##安装补充包：
apt-get install -y curl gpg net-tools

[ -f /etc/rc.local ] && chmod +x /etc/rc.local

###安装根证书
cat << EOF > /usr/local/share/ca-certificates/ROBOT_Root_CA.crt
xxx
EOF
#import ROBOT domain ac cert
cat << EOF > /usr/local/share/ca-certificates/ROBOT_Domain_CA.crt
xxx
EOF
update-ca-certificates

apt install libnss3-tools -y

##
# Ubuntu无线网络准入
##
echo -e $(date)" computer name is "${hostname} >>/var/log/init_ad.log

# 1. Ubuntu自动加域脚本，上面 JoinAD.sh 已实现
# 2. Ubuntu电脑自动注册证书
[ -d /etc/pki/tls/certs ] || sudo mkdir -p /etc/pki/tls/certs
[ -d /etc/pki/tls/private ] || sudo mkdir -p /etc/pki/tls/private
sudo apt install -y certmonger
getcert add-scep-ca -c ultron_ca -u  http://172.16.110.13/certsrv/mscep/mscep.dll
getcert request -I ultron_cn -c ultron_ca -f /etc/pki/tls/certs/server_cn.crt -k /etc/pki/tls/private/private_cn.key -N cn="${hostname}" -D "${hostname}" -l /tmp/challenge_password

# 检查签发状态
for i in `seq 15`; do
    t=20
    echo -e $(date)" wait ${t}s for cert request ..." | tee -a /var/log/init_ad.log
    sleep ${t}
    getcert list | grep "status: " | tee -a /var/log/init_ad.log
    if [ -f /etc/pki/tls/certs/server_cn.crt ]; then
        echo -e $(date)" found /etc/pki/tls/certs/server_cn.crt." | tee -a /var/log/init_ad.log
        break
    fi
done

if [ ! -f /etc/pki/tls/certs/server_cn.crt ]; then
    echo -e $(date)" not found /etc/pki/tls/certs/server_cn.crt, exit 1" | tee -a /var/log/init_ad.log
    exit 1
fi

# 输出 server_cn.crt
echo -e $(date)" new server_cn.crt text is " >>/var/log/init_ad.log
openssl x509 -noout -text -in /etc/pki/tls/certs/server_cn.crt | tee -a /var/log/init_ad.log

wifi_ifname=$(iwconfig 2>&1 | grep -v "no wireless extensions." | grep "ESSID:" | awk '{print $1}')
if [[ "${wifi_ifname}"x == ""x ]]; then
    wifi_ifname=$(iwconfig 2>&1 | grep "^wl" | awk '{print $1}')
fi
if [[ "${wifi_ifname}"x == ""x ]]; then
    echo "no wifi ifname found, try run iwconfig, exit 1"
    exit 1
fi

# 3. Ubuntu生成WPA配置并自动连接
# 3.1 生成802.1x和EAP的配置文件
# cat << EOF > /etc/wpa_supplicant/wpa_supplicant-${wifi_ifname}.conf
# ctrl_interface=/var/run/wpa_supplicant-${wifi_ifname}
# country=CN
# ap_scan=1
# network={
#      ssid="ultron_WLAN"
#      scan_ssid=1
#      key_mgmt=WPA-EAP
#      eap=TLS
#      identity="host/${hostname}"
#      ca_cert="/usr/local/share/ca-certificates/ROBOT_Domain_CA.crt"
#      client_cert="/etc/pki/tls/certs/server_cn.crt"
#      private_key="/etc/pki/tls/private/private_cn.key"
# }
# EOF

# 3.2 激活对应无线网卡的服务
# sudo systemctl disable "wpa_supplicant@${wifi_ifname}.service"
sudo systemctl enable wpa_supplicant.service
sudo systemctl restart wpa_supplicant.service

# 3.3 配置启动 dhclient 来获得IP（使用NetworkManager自动获取）
# cat << EOF > /etc/systemd/system/dhclient.service
# [Unit]
# Description= DHCP Client
# Before=network.target

# [Service]
# Type=forking
# ExecStart=/sbin/dhclient ${wifi_ifname} -v
# ExecStop=/sbin/dhclient ${wifi_ifname} -r
# Restart=always

# [Install]
# WantedBy=multi-user.target
# EOF

# sudo systemctl enable dhclient.service

# 3.4 添加nmcli 的证书认证配置文件
if [ `nmcli connection show | grep -c "ultron_WLAN"` -eq '0' ];then
    sudo nmcli connection add type wifi ifname ${wifi_ifname} \
        con-name ultron_WLAN \
        802-11-wireless.ssid ultron_WLAN \
        802-11-wireless-security.key-mgmt "wpa-eap" \
        802-1x.eap "tls" \
        802-1x.identity "host/${hostname}" \
        802-1x.ca-cert '/usr/local/share/ca-certificates/ROBOT_Domain_CA.crt' \
        802-1x.client-cert /etc/pki/tls/certs/server_cn.crt \
        802-1x.private-key /etc/pki/tls/private/private_cn.key \
        802-1x.private-key-password "xxxxxxxx"
else
    echo "ultron_WLAN is already created, run:"
    echo "1. 'nmcli connection show ultron_WLAN' for detail"
    echo "2. 'nmcli connection delete ultron_WLAN' to delete it"
fi

# realm list
# systemctl status wpa_supplicant@<name>.service
# wpa_cli show
# nmcli connection show
# nmcli connection delete uuid
```