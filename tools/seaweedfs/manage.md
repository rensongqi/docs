- [SeaweedFS管理维护手册](#seaweedfs管理维护手册)
- [1 手动编译](#1-手动编译)
- [2 元数据备份](#2-元数据备份)
- [3 报错处理](#3-报错处理)
  - [3.1 .idx does not exists](#31-idx-does-not-exists)
  - [3.2 df -h hang死](#32-df--h-hang死)
    - [3.3 volume vacuum failed](#33-volume-vacuum-failed)
    - [3.3.1 手动清理](#331-手动清理)
    - [3.3.2 半自动化清理](#332-半自动化清理)
    - [3.3.3 weed mount lost data](#333-weed-mount-lost-data)
- [4 使用elastic7作为Filer store](#4-使用elastic7作为filer-store)
- [5 weed shell](#5-weed-shell)
  - [5.1 标记卷只读/读写](#51-标记卷只读读写)
  - [5.2 手动创建卷](#52-手动创建卷)
  - [5.3 移动卷](#53-移动卷)
- [6 快速定位文件有没有损坏](#6-快速定位文件有没有损坏)

# SeaweedFS管理维护手册

# 1 手动编译

```bash
git clone https://github.com/seaweedfs/seaweedfs.git
cd seaweedfs/docker
make go_build_large_disk
```

# 2 元数据备份

> 备份目录：/disk/upload/devops/seaweedfs_bakcup
> 备份脚本：backup.sh
```bash
#!/bin/bash

LOG_PATH="/var/log/weed"
full_base_time=`date +%Y%m%d`
redis_base_time=`date +%Y%m%d%H%M%S`
day_of_week=`date +%u`
backup_base_path="/disk/upload/devops/seaweedfs_bakcup/<BACKUP_IP>"
full_backup_filepath=${backup_base_path}/full/${full_base_time}
redis_backup_filepath=${backup_base_path}/redis/${redis_base_time}
sunday_backup_filepath=${backup_base_path}/sunday/${full_base_time}

# Create log path
if [[ ! -d "$LOG_PATH" ]]; then
  mkdir -p ${LOG_PATH}
fi

# Backup files daily
if [[ ! -d "$full_backup_filepath" ]]; then
  mkdir -p ${full_backup_filepath}
  rsync -avP /data/* ${full_backup_filepath}/ >> /var/log/weed/full_${full_base_time}.log
fi

# Backup files on sunday
if [ $day_of_week -eq 7 ]; then
  if [[ ! -d "$sunday_backup_filepath" ]]; then
    mkdir -p ${sunday_backup_filepath}
    rsync -avP /data/* ${sunday_backup_filepath}/ >> /var/log/weed/sunday_${full_base_time}.log
  fi
fi

# Backup Redis files every half hour
mkdir -p ${redis_backup_filepath}
rsync -avP /data/redis/data ${redis_backup_filepath}/ >> /var/log/weed/redis_${redis_base_time}.log

# clean backup tar
find ${redis_backup_filepath}/ -type d -ctime +1 -exec rm -f -r {} \;
find ${full_backup_filepath}/ -type d -ctime +7 -exec rm -f -r {} \;
find ${sunday_backup_filepath}/ -type d -ctime +90 -exec rm -f -r {} \;
```

# 3 报错处理

## 3.1 .idx does not exists

报错：
```bash
F1014 17:48:18.794113 volume_loading.go:117 check volume idx file /weed/metadata/477.idx: idx file /weed/metadata/477.idx does not exists
```

修复：

到volume对应的宿主机，找到volume对应的磁盘路径，执行如下命令
```bash
weed fix -volumeId 477 /disk/local_disk1/
weed fix -volumeId 477 /disk/local_disk2/
```

## 3.2 df -h hang死
```bash
mount | grep fuse
umount  -l <mount point>
```

### 3.3 volume vacuum failed

> 卷垃圾自动清理失败，需要手动清理卷垃圾，找到`volume id`对应的Volume节点中具体的存储目录，执行如下命令

### 3.3.1 手动清理
```bash
find /disk/local_disk* -type f -name "collection_1174*"

weed compact -method 1 -dir=/disk/local_disk2 -volumeId=1174 -collection=collection

# 重命名，需要将原有数据跟现有数据改名
mv collection_1174.dat collection_1174.dat.bak
mv collection_1174.idx collection_1174.idx.bak

# 将新生成的文件重命名
mv collection_1174.cpd collection_1174.dat
mv collection_1174.cpx collection_1174.idx

# 重启volume节点
```
### 3.3.2 半自动化清理
> 需要先找出来垃圾比较大或failed compact volume的volume所在服务器的绝对路径

```bash
docker logs seaweedfs-volume --tail=20000 2>&1 | grep 'failed compact volume'

find /disk/local_disk* -type f -name "*_3070.dat"
```

`check1.txt`
```
/disk/local_disk2/collection_3061.dat
/disk/local_disk4/collection_3073.dat
/disk/local_disk5/collection_3094.dat
/disk/local_disk3/collection_3101.dat
/disk/local_disk2/collection_3107.dat
/disk/local_disk4/collection_3121.dat
/disk/local_disk3/collection_3124.dat
/disk/local_disk1/collection_3125.dat
/disk/local_disk3/collection_3130.dat
/disk/local_disk6/collection_3137.dat
```

> 执行完compact.sh脚本之后需要手动重启seaweedfs-volume服务，然后再手动删除.bak备份的文件以释放空间

`compact.sh`

```bash
#!/bin/bash

SOURCE_FILE="check1.txt"

while read -r line; do
    TARGET_DISK=$(echo "$line" | cut -d'/' -f3) # 提取 local_disk 
    FILE_NAME=$(basename "$line" | cut -d'.' -f1) # 提取文件名（不含后缀）
    VOLUME_ID=$(basename "$line" | cut -d'.' -f1 | cut -d'_' -f2) # 获取卷ID
    VOLUME_BUCKET_NAME=$(basename "$line" | cut -d'.' -f1 | cut -d'_' -f1) # 获取bucket_name

    # 拼接源路径和目标路径
    SOURCE_FILE_PATH="/data/seaweedfs/metadata/${FILE_NAME}.idx"
    TARGET_DIR="/disk/${TARGET_DISK}/"

    # 检查源文件是否存在
    if [[ -f "$SOURCE_FILE_PATH" ]]; then
        # 拷贝文件到目标目录
        cp "${SOURCE_FILE_PATH}" "${TARGET_DIR}"
        weed compact -dir=${TARGET_DIR} -volumeId=${VOLUME_ID} -collection=${VOLUME_BUCKET_NAME} >> /var/log/weed_compact1.log 2>&1 
        if [[ $? -eq 0 ]]; then
            mv ${TARGET_DIR}${FILE_NAME}.dat ${TARGET_DIR}${FILE_NAME}.dat.bak
            mv ${TARGET_DIR}${FILE_NAME}.idx ${TARGET_DIR}${FILE_NAME}.idx.bak
            mv ${TARGET_DIR}${FILE_NAME}.cpx ${TARGET_DIR}${FILE_NAME}.idx
            mv ${TARGET_DIR}${FILE_NAME}.cpd ${TARGET_DIR}${FILE_NAME}.dat
            cp ${TARGET_DIR}${FILE_NAME}.idx /data/seaweedfs/metadata/
        fi
    else
        echo "Source file $SOURCE_FILE_PATH does not exist"
    fi
done < "$SOURCE_FILE"
```

### 3.3.3 weed mount lost data

> 报错信息，更新客户端版本至 3.80+ 之后解决问题

```bash
I1212 18:06:12.574545 weedfs_stats.go:38 reading filer stats ttl:"0s": rpc error: code = Unavailable desc = keepalive ping failed to receive ACK within timeout
I1212 18:06:12.575669 wfs_filer_client.go:31 WithFilerClient 0 172.16.90.174:18888: rpc error: code = Unavailable desc = keepalive ping failed to receive ACK within timeout
I1212 18:22:43.750945 weedfs_stats.go:38 reading filer stats ttl:"0s": rpc error: code = Unknown desc = raft.Server: Not current leader
I1212 18:22:43.751036 wfs_filer_client.go:31 WithFilerClient 0 172.16.90.172:18888: rpc error: code = Unknown desc = raft.Server: Not current leader
I1212 18:22:43.780581 weedfs_stats.go:38 reading filer stats ttl:"0s": rpc error: code = Unknown desc = raft.Server: Not current leader
I1212 18:22:43.780624 wfs_filer_client.go:31 WithFilerClient 1 172.16.90.173:18888: rpc error: code = Unknown desc = raft.Server: Not current leader
```

# 4 使用elastic7作为Filer store
>新创建的ES节点需要修改.seaweedfs_bucketsindex Inode type为unsigned_long
```bash
PUT /.seaweedfs_buckets
{
  "settings": {
    "number_of_shards": 5,
    "number_of_replicas": 1
  },
  "mappings": {
      "properties": {
        "Entry": {
          "properties": {
            "Crtime": {
              "type": "date"
            },
            "Extended": {
              "properties": {
                "X-Amz-Meta-S3cmd-Attrs": {
                  "type": "text",
                  "fields": {
                    "keyword": {
                      "type": "keyword",
                      "ignore_above": 256
                    }
                  }
                },
                "x-amz-storage-class": {
                  "type": "text",
                  "fields": {
                    "keyword": {
                      "type": "keyword",
                      "ignore_above": 256
                    }
                  }
                },
                "xattr-system": {
                  "properties": {
                    "posix_acl_access": {
                      "type": "text",
                      "fields": {
                        "keyword": {
                          "type": "keyword",
                          "ignore_above": 256
                        }
                      }
                    }
                  }
                }
              }
            },
            "FileSize": {
              "type": "long"
            },
            "FullPath": {
              "type": "text",
              "fields": {
                "keyword": {
                  "type": "keyword",
                  "ignore_above": 256
                }
              }
            },
            "Gid": {
              "type": "long"
            },
            "HardLinkCounter": {
              "type": "long"
            },
            "Inode": {
              "type": "unsigned_long"
            },
            "Md5": {
              "type": "text",
              "fields": {
                "keyword": {
                  "type": "keyword",
                  "ignore_above": 256
                }
              }
            },
            "Mime": {
              "type": "text",
              "fields": {
                "keyword": {
                  "type": "keyword",
                  "ignore_above": 256
                }
              }
            },
            "Mode": {
              "type": "long"
            },
            "Mtime": {
              "type": "date"
            },
            "Quota": {
              "type": "long"
            },
            "Rdev": {
              "type": "long"
            },
            "SymlinkTarget": {
              "type": "text",
              "fields": {
                "keyword": {
                  "type": "keyword",
                  "ignore_above": 256
                }
              }
            },
            "TtlSec": {
              "type": "long"
            },
            "Uid": {
              "type": "long"
            },
            "UserName": {
              "type": "text",
              "fields": {
                "keyword": {
                  "type": "keyword",
                  "ignore_above": 256
                }
              }
            },
            "chunks": {
              "properties": {
                "e_tag": {
                  "type": "text",
                  "fields": {
                    "keyword": {
                      "type": "keyword",
                      "ignore_above": 256
                    }
                  }
                },
                "fid": {
                  "properties": {
                    "cookie": {
                      "type": "long"
                    },
                    "file_key": {
                      "type": "long"
                    },
                    "volume_id": {
                      "type": "long"
                    }
                  }
                },
                "is_compressed": {
                  "type": "boolean"
                },
                "modified_ts_ns": {
                  "type": "long"
                },
                "offset": {
                  "type": "long"
                },
                "size": {
                  "type": "long"
                }
              }
            }
          }
        },
        "ParentId": {
          "type": "text",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 256
            }
          }
        }
      }
  }
}
```

# 5 weed shell

> 执行 lock 时如果发现被其他节点已经lock了，确认了其它节点并没有手动执行的lock进程，那么可以找到其他节点当前的登录会话，把这些会话tty踢下线即可

```bash
w
pkill -kill -t pts/1
```

## 5.1 标记卷只读/读写
```bash
weed shell
lock
volume.mark -node 172.16.104.124:8080 -readonly -volumeId 42
volume.mark -node 172.16.104.124:8080 -writable -volumeId 42
unlock
```

## 5.2 手动创建卷

```bash
weed shell
volume.grow -collection upload -count 4

# 在指定数据节点创建volume
volume.grow -collection collection -count 6 -dataNode 172.16.90.176:8080
```

## 5.3 移动卷
> 前提条件需要保证卷处于可写状态，只读状态会移动失败

`volume.move`
```bash
lock

volume.move -source 172.16.104.124:8080 -target 172.16.90.177:8080 -volumeId 37

# 多条命令执行加 ; 分号即可
volume.move xxx ; volume.move xxx ; ...

unlock
```

`volumeServer.evacuate`
```bash
lock

# 如下命令会把172.16.90.172节点的所有数据转移到其它节点中
# 该命令会遍历所有卷ID，逐一迁移卷ID的数据，每完成一个卷ID的迁移会删除其对应的数据
volumeServer.evacuate -node 172.16.90.172:8080 -force

# 指定卷迁移的目标机器
volumeServer.evacuate -node 172.16.90.172:8080 -force -target 172.16.90.173:8080

# 标记该Volume离开集群，执行之后在master web界面看不到此volume
volumeServer.leave -node 172.16.90.172:8080

unlock
```

# 6 快速定位文件有没有损坏

```bash
for i in `find /path/to/file -type f`; do timeout 2s head $i > /dev/null 2>&1 ; if [[ $? != 0 ]]; then echo "$i" >> ~/lost.txt ; fi; done &  
```