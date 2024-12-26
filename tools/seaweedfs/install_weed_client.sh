#!/bin/bash

if [[ ! -f /usr/bin/weed ]]; then
  wget https://mirrors.rsq.cn/public/pkgs/weed_linux_amd64_full_large_disk.tar.gz -P /tmp/
  tar -xf /tmp/weed_linux_amd64_full_large_disk.tar.gz -C /usr/bin/
fi

mkdir -p /var/log/weed/
mkdir /mount/point

# weed -v 4 (log level)
cat <<EOF > /etc/systemd/system/seaweedfs-mount.service
[Unit]
Description=SeaweedFS FUSE Mount Service
After=network-online.target

[Service]
Type=simple
ExecStart=/bin/sh -c '/usr/bin/weed -v 4 mount -filer=weedfs.rsq.cn:8888 -filer.path=/buckets/test -dir=/mount/point -cacheDir=/mnt -nonempty >> /var/log/weed/client.log 2>&1'
ExecStop=/bin/umount -l /mount/piont
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable seaweedfs-mount.service --now
systemctl start seaweedfs-mount.service

# config log compress
cat <<EOF > /etc/logrotate.d/weed
/var/log/weed/*.log {
    create 0640 root root
    daily
    rotate 5
    dateext
    copytruncate
    compress
}
EOF