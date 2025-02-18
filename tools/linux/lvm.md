- [使用lvm根分区扩容](#使用lvm根分区扩容)
- [lvm合并两块物理磁盘](#lvm合并两块物理磁盘)
- [问题记录](#问题记录)
  - [GPT PMBR size mismatch will be corrected by write错误](#gpt-pmbr-size-mismatch-will-be-corrected-by-write错误)
  - [lvextend之后发现扩容的文件系统不对，需要还原](#lvextend之后发现扩容的文件系统不对需要还原)

# 使用lvm根分区扩容

> 需要把根分区从200G扩容至500G

```bash
root@SHJDSMAPGROUP:~# df -h​
Filesystem                         Size  Used Avail Use% Mounted on​
udev                                32G     0   32G   0% /dev​
tmpfs                              6.3G  1.3M  6.3G   1% /run​
/dev/mapper/ubuntu--vg-ubuntu--lv  196G  6.4G  181G   4% /​
tmpfs                               32G     0   32G   0% /dev/shm​
tmpfs                              5.0M     0  5.0M   0% /run/lock​
tmpfs                               32G     0   32G   0% /sys/fs/cgroup​
/dev/sda2                          976M   80M  829M   9% /boot​
tmpfs                              6.3G     0  6.3G   0% /run/user/1002​
tmpfs                              6.3G     0  6.3G   0% /run/user/1001
```

创建新分区

```bash
# fdisk 创建一个lvm格式分区
root@SHJDSMAPGROUP:~# fdisk /dev/sda
 Command (m for help): n
Partition number (4-128, default 4):  
First sector (419428352-1048575966, default 419428352): 
Last sector, +sectors or +size{K,M,G,T,P} (419428352-1048575966, default 1048575966): 

Created a new partition 4 of type 'Linux filesystem' and of size 300 GiB.

Command (m for help): p
Disk /dev/sda: 500 GiB, 536870912000 bytes, 1048576000 sectors
Units: sectors of 1 * 512 = 512 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 512 bytes
Disklabel type: gpt
Disk identifier: 1F724393-C34E-404D-8560-CB85D9F719E4

Device         Start        End   Sectors  Size Type
/dev/sda1       2048       4095      2048    1M BIOS boot
/dev/sda2       4096    2101247   2097152    1G Linux filesystem
/dev/sda3    2101248  419428351 417327104  199G Linux filesystem
/dev/sda4  419428352 1048575966 629147615  300G Linux filesystem

Command (m for help): t
Partition number (1-4, default 4): 4
Partition type (type L to list all types): 31  # ubuntu选31，centos选8e

Changed type of partition 'Linux filesystem' to 'Linux LVM'.

Command (m for help): w
GPT PMBR size mismatch (419430399 != 1048575999) will be corrected by w(rite).

The partition table has been altered.
Syncing disks.

# 重读分区表
root@SHJDSMAPGROUP:~# partprobe 

# 格式化新分区 Ubuntu18.04
root@SHJDSMAPGROUP:~# mkfs.ext4 /dev/sda4

# 格式化新分区 Centos7(具体还是看文件系统类型，默认centos是xfs)
root@SHJDSMAPGROUP:~# mkfs.xfs /dev/sda4
```

磁盘扩容

```bash
# 1 查看vgdisplay
root@SHJDSMAPGROUP:~# vgdisplay 
  --- Volume group ---
  VG Name               ubuntu-vg
  System ID             
  Format                lvm2
  Metadata Areas        1
  Metadata Sequence No  18
  VG Access             read/write
  VG Status             resizable
  MAX LV                0
  Cur LV                1
  Open LV               1
  Max PV                0
  Cur PV                1
  Act PV                1
  VG Size               <199.00 GiB
  PE Size               4.00 MiB
  Total PE              50943
  Alloc PE / Size       50943 / <199.00 GiB
  Free  PE / Size       0 / 0   
  VG UUID               F0Fpo3-nX8V-UhvF-2B4A-YlLc-w5mo-XsrE1K
  
# 2 创建pv，把刚才创建的新分区创建为pv
root@SHJDSMAPGROUP:~# pvcreate /dev/sda4
  Physical volume "/dev/sda4" successfully created.
  
# 3 将新的pv加入到此vg中
root@SHJDSMAPGROUP:~# vgextend ubuntu-vg /dev/sda4
  Volume group "ubuntu-vg" successfully extended

# 4 扩展逻辑虚拟卷lv的容量
root@SHJDSMAPGROUP:~# lvextend -l +100%FREE /dev/mapper/ubuntu--vg-ubuntu--lv
  Size of logical volume ubuntu-vg/ubuntu-lv changed from <199.00 GiB (50943 extents) to 498.99 GiB (127742 extents).
  Logical volume ubuntu-vg/ubuntu-lv successfully resized.
  
# 5 上述只是对卷扩容，还需要对文件系统实现真正的扩容
# CentOS 7 下面 由于使用的是 XFS，所以要用
xfs_growfs /dev/mapper/centos-root

# Ubuntu18.04 下面 要用
resize2fs /dev/mapper/ubuntu--vg-ubuntu--lv

# 6 查看扩容后的分区
root@SHJDSMAPGROUP:~# df -h
Filesystem                         Size  Used Avail Use% Mounted on
udev                                32G     0   32G   0% /dev
tmpfs                              6.3G  1.3M  6.3G   1% /run
/dev/mapper/ubuntu--vg-ubuntu--lv  491G  6.4G  464G   2% /
tmpfs                               32G     0   32G   0% /dev/shm
tmpfs                              5.0M     0  5.0M   0% /run/lock
tmpfs                               32G     0   32G   0% /sys/fs/cgroup
/dev/sda2                          976M   80M  829M   9% /boot
tmpfs                              6.3G     0  6.3G   0% /run/user/1002
tmpfs                              6.3G     0  6.3G   0% /run/user/1001
```

# lvm合并两块物理磁盘

> 合并`/dev/sda`和`/dev/sdb`

安装lvm2工具包(如果尚未安装)
```bash
sudo apt install lvm2   # Debian/Ubuntu系统
# 或
sudo yum install lvm2   # RHEL/CentOS系统
```

将硬盘格式化为LVM格式：

```bash
# 使用fdisk创建LVM分区
sudo fdisk /dev/sda
# 输入n新建分区
# 输入p创建主分区
# 分区号使用默认的1
# 其他选项使用默认值
# 输入t更改分区类型
# 输入8e选择Linux LVM类型
# 输入w保存更改

# 对sdb重复相同操作
sudo fdisk /dev/sdb

```

创建物理卷(PV):
```bash
sudo pvcreate /dev/sda1
sudo pvcreate /dev/sdb1

# 检查物理卷
sudo pvs
```

创建卷组(VG):
```bash
# 创建名为vg_data的卷组，并将第一个物理卷添加进去
sudo vgcreate vg_data /dev/sda1

# 将第二个物理卷扩展到卷组中
sudo vgextend vg_data /dev/sdb1

# 检查卷组
sudo vgs
```
创建逻辑卷(LV):
```bash
# 创建使用全部空间的逻辑卷
sudo lvcreate -l 100%VG -n lv_data vg_data

# 检查逻辑卷
sudo lvs
```
格式化逻辑卷：
```bash
# 使用ext4文件系统格式化
sudo mkfs.ext4 /dev/vg_data/lv_data
```

创建挂载点并挂载：
```bash
# 创建挂载点
sudo mkdir /data

# 挂载逻辑卷
sudo mount /dev/vg_data/lv_data /data

# 检查挂载情况
df -h
```
设置开机自动挂载：
```bash
# 编辑/etc/fstab文件
sudo nano /etc/fstab

# 添加以下行
/dev/vg_data/lv_data /data ext4 defaults 0 0

```

# 问题记录
## GPT PMBR size mismatch will be corrected by write错误
> 原因：在对虚拟机扩容时候，由于linux系统没有对其磁盘信息进行更新，导致了磁盘实际容量和linux系统容量不一致
修复
```bash
# 执行命令：
sudo parted -l

# 然后输入：
Fix
```

## lvextend之后发现扩容的文件系统不对，需要还原
> 根分区有1.75TB实际可用空间，需要把多余的两块8T的盘加入根分区，在lvextend之后发现原来8T的盘被格式化为ext4，然后根分区是xfs文件系统，这会导致xfs_growfs resizefs的时候报错，所以需要把加入的两个pv，分别是/dev/sda1和/dev/sdb1给删掉。

修复

> 需要按照以下步骤严格执行，否则会造成文件系统错误
1. 查看现有vg状态，使用命令：vgdisplay
```bash
[root@hadoop111 ~]# vgdisplay                                             
  --- Volume group ---
  VG Name               centos
  System ID             
  Format                lvm2
  Metadata Areas        3
  Metadata Sequence No  6
  VG Access             read/write
  VG Status             resizable
  MAX LV                0
  Cur LV                2
  Open LV               2
  Max PV                0
  Cur PV                3
  Act PV                3
  VG Size               <15.72 TiB
  PE Size               4.00 MiB
  Total PE              4120269
  Alloc PE / Size       4120269 / <15.72 TiB
  Free  PE / Size       0 / 0   
  VG UUID               3Z4pzJ-iXuk-RuhH-UU2w-JCZR-1Khy-ir35yb
```

2. 需要找出之前的状态，查看之前的数据盘大小，翻历史记录发现之前的大小是1.75TB，这时候需要缩小lv的大小为1.7TB，保证有充足的扇区可以还原给/dev/sda1和/dev/sdb1

```bash
  [root@hadoop112 ~]# vgdisplay 
  --- Volume group ---
  VG Name               centos
  System ID             
  Format                lvm2
  Metadata Areas        1
  Metadata Sequence No  3
  VG Access             read/write
  VG Status             resizable
  MAX LV                0
  Cur LV                2
  Open LV               2
  Max PV                0
  Cur PV                1
  Act PV                1
  VG Size               <1.75 TiB
  PE Size               4.00 MiB
  Total PE              457445
  Alloc PE / Size       457445 / <1.75 TiB
  Free  PE / Size       0 / 0   
  VG UUID               xDwaya-VH9C-wqLh-m2I3-j0SB-7Obe-cH53fK
```

3. 缩小根分区lv大小，使用命令：lvreduce -L 1.7T /dev/mapper/centos-root
执行完命令之后再查看下vg，需要保证Alloc PE / Size的值小于等于原来的457445

```bash
[root@hadoop111 ~]# lvreduce -L 1.7T /dev/mapper/centos-root
  Rounding size to boundary between physical extents: 1.70 TiB.
  WARNING: Reducing active and open logical volume to 1.70 TiB.
  THIS MAY DESTROY YOUR DATA (filesystem etc.)
Do you really want to reduce centos/root? [y/n]: y
  Size of logical volume centos/root changed from <1.75 TiB (457445 extents) to 1.70 TiB (445645 extents).
  Logical volume centos/root successfully resized.
  

[root@hadoop111 ~]# vgdisplay 
  --- Volume group ---
  VG Name               centos
  System ID             
  Format                lvm2
  Metadata Areas        2
  Metadata Sequence No  16
  VG Access             read/write
  VG Status             resizable
  MAX LV                0
  Cur LV                2
  Open LV               2
  Max PV                0
  Cur PV                2
  Act PV                2
  VG Size               8.73 TiB
  PE Size               4.00 MiB
  Total PE              2288857
  Alloc PE / Size       446669 / 1.70 TiB
  Free  PE / Size       1842188 / <7.03 TiB
  VG UUID               3Z4pzJ-iXuk-RuhH-UU2w-JCZR-1Khy-ir35yb
```

4. 还原/dev/sda1和/dev/sdb1的pv大小为默认的
```bash
[root@hadoop111 ~]# pvresize /dev/sda1
  Physical volume "/dev/sda1" changed
  1 physical volume(s) resized or updated / 0 physical volume(s) not resized
[root@hadoop111 ~]# pvresize /dev/sdb1
  Physical volume "/dev/sdb1" changed
  1 physical volume(s) resized or updated / 0 physical volume(s) not resized
```

5. 把原来vgextend的两块盘去掉

```bash
[root@hadoop111 ~]# vgreduce centos /dev/sda1
  Removed "/dev/sda1" from volume group "centos"
[root@hadoop111 ~]# vgreduce centos /dev/sdb1
  Removed "/dev/sdb1" from volume group "centos"
```

6. 删除pv

```bash
[root@hadoop111 ~]# pvremove /dev/sda1
  Labels on physical volume "/dev/sda1" successfully wiped.
[root@hadoop111 ~]# pvremove /dev/sdb1
  Labels on physical volume "/dev/sdb1" successfully wiped.
```

7. 扩容根分区为它原本的最大空间

```bash
[root@hadoop111 ~]# lvextend -l +100%FREE /dev/mapper/centos-root
  Size of logical volume centos/root changed from 1.70 TiB (445645 extents) to 1.74 TiB (456421 extents).
  Logical volume centos/root successfully resized.
[root@hadoop111 ~]# vgdisplay 
  --- Volume group ---
  VG Name               centos
  System ID             
  Format                lvm2
  Metadata Areas        1
  Metadata Sequence No  18
  VG Access             read/write
  VG Status             resizable
  MAX LV                0
  Cur LV                2
  Open LV               2
  Max PV                0
  Cur PV                1
  Act PV                1
  VG Size               <1.75 TiB
  PE Size               4.00 MiB
  Total PE              457445
  Alloc PE / Size       457445 / <1.75 TiB
  Free  PE / Size       0 / 0   
  VG UUID               3Z4pzJ-iXuk-RuhH-UU2w-JCZR-1Khy-ir35yb
# 对文件系统实现真正的扩容
[root@hadoop112 ~]# xfs_growfs /dev/mapper/centos-root
meta-data=/dev/mapper/centos-root isize=512    agcount=32, agsize=14605504 blks
         =                       sectsz=512   attr=2, projid32bit=1
         =                       crc=1        finobt=0 spinodes=0
data     =                       bsize=4096   blocks=467375104, imaxpct=5
         =                       sunit=64     swidth=64 blks
naming   =version 2              bsize=4096   ascii-ci=0 ftype=1
log      =internal               bsize=4096   blocks=228224, version=2
         =                       sectsz=512   sunit=64 blks, lazy-count=1
realtime =none                   extsz=4096   blocks=0, rtextents=0
```