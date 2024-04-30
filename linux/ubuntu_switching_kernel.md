
```bash
apt install linux-image-5.15.0-25-generic linux-headers-5.15.0-25-generic
rm -rf /etc/default/grub
cat<<E0F >/etc/default/grub
GRUB_DEFAULT="Advanced options for Ubuntu>Ubuntu, with Linux 5.15.0-25-generic"
GRUB_TIMEOUT_STYLE=hidden
GRUB_TIMEOUT=2
GRUB_DISTRIBUTOR=`lsb release -i -s 2>/dev/null || echo Debian
GRUB_CMDLINE_LINUX_DEFAULT=""
GRUB_CMDLINE_LINUX=""
EOF
update-grub
update-grub2

reboot
```