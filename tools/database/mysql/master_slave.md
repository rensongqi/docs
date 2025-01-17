
# 1 master

```bash
# mkdir /data/mysql-master/volumes/{log,conf,data} -p

# vim /data/mysql-master/docker-compose.yml
version: "2"
services:
  mysql-master:
    container_name: mysql-master
    image: harbor.rsq.cn/library/mysql/mysql:5.7
    restart: always
    ports:
      - "3306:3306"
    privileged: true
    volumes:
      - ./volumes/log:/var/log/mysql
      - ./volumes/conf/my.cnf:/etc/mysql/my.cnf
      - ./volumes/data:/var/lib/mysql
    environment:
      MYSQL_ROOT_PASSWORD: "xxxxxx"
      TZ: "Asia/Shanghai"
    command: [
      "--character-set-server=utf8mb4",
      "--collation-server=utf8mb4_general_ci",
      "--max_connections=30000"
    ]

# 编辑配置文件
# vim /data/mysql-master/volumes/conf/my.cnf
[mysqld]
server_id=1000
binlog-ignore-db=mysql
log-bin=replicas-mysql-bin
binlog_cache_size=1M
binlog_format=mixed
expire_logs_days=7
slave_skip_errors=1062
innodb_buffer_pool_size = 8G
join_buffer_size = 8M
sort_buffer_size = 8M
read_buffer_size = 4M
read_rnd_buffer_size = 8M

# 启动容器
# cd /data/mysql-master
# docker-compose up -d
```

# 2 slave

```bash
mkdir /data/mysql-slave/volumes/{log,conf,data} -p

vim /data/mysql-slave/docker-compose.yml
version: "2"
services:
  mysql-slave:
    container_name: mysql-slave
    image: harbor.rsq.cn/library/mysql/mysql:5.7
    restart: always
    ports:
      - "3306:3306"
    privileged: true
    volumes:
      - ./volumes/log:/var/log/mysql
      - ./volumes/conf/my.cnf:/etc/mysql/my.cnf
      - ./volumes/data:/var/lib/mysql
    environment:
      MYSQL_ROOT_PASSWORD: "xxxxxx"
      TZ: "Asia/Shanghai"
    command: [
      "--character-set-server=utf8mb4",
      "--collation-server=utf8mb4_general_ci",
      "--max_connections=30000"
    ]

# 编辑配置文件
vim /data/mysql-slave/volumes/conf/my.cnf
[mysqld]
server_id=2000
binlog-ignore-db=mysql
log-bin=replicas-mysql-slave1-bin
binlog_cache_size=1M
binlog_format=mixed
expire_logs_days=7
slave_skip_errors=1062
relay_log=replicas-mysql-relay-bin
log_slave_updates=1
read_only=1
innodb_buffer_pool_size = 8G
join_buffer_size = 8M
sort_buffer_size = 8M
read_buffer_size = 4M
read_rnd_buffer_size = 8M

# 启动容器
docker-compose up -d
```

# 3 配置主从
进入master容器，查看 file 和 pos 数据
```bash
mysql -uroot -p

mysql> show master status;
+---------------------------+----------+--------------+------------------+-------------------+
| File                      | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+---------------------------+----------+--------------+------------------+-------------------+
| replicas-mysql-bin.000003 |      154 |              | mysql            |                   |
+---------------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)
```

进入slave容器

```bash
mysql -uroot -p

mysql> change master to master_host='172.16.104.104',master_user='root',master_password='xxxxxx',master_port=3306,master_log_file='replicas-mysql-bin.000003', master_log_pos=154,master_connect_retry=30;
mysql> start slave;
mysql> show slave status \G；
```

测试主从同步情况，master创建数据库，查看slave有没有