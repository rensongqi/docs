# 1 安装postgresql
```bash
sudo yum install -y https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm

sudo yum install -y postgresql14-server
sudo /usr/pgsql-14/bin/postgresql-14-setup initdb
sudo systemctl enable postgresql-14
sudo systemctl start postgresql-14
```

# 2 修改配置文件
```bash
mkdir /var/lib/pgsql/14/data/pg_archive
vim /var/lib/pgsql/14/data/postgresql.conf
#修改监听所有ip
listen_addresses = '*'
archive_mode = on
archive_command = 'test ! -f /var/lib/pgsql/14/data/pg_archive/%f && cp %p /var/lib/pgsql/14/data/pg_archive/%f'
wal_level = replica
wal_sender_timeout = 60s

vim /var/lib/pgsql/14/data/pg_hba.conf
#增加一行
host    all             all             0.0.0.0/0               password

sudo systemctl restart postgresql-14
```

# 3 修改postgres用户密码，并配置复制流用户
```bash
su - postgres
psql
alter user postgres with password 'cowa2022';

# 创建复制
CREATE ROLE replica login replication encrypted password 'replica';

# 修改配置文件
vim /var/lib/pgsql/14/data/pg_hba.conf
# 添加一行
host    replication     replica         172.16.104.0/24         trust

```

# 4 配置从数据库
```bash
su - postgres
rm -rf /var/lib/pgsql/14/data/*

pg_basebackup -h 172.16.104.109 -p 5432 -U replica -Fp -Xs -Pv -R -D /var/lib/pgsql/14/data/
```

修改standby.signal
```bash
vim /var/lib/pgsql/14/data/standby.signal
# 添加
standby_mode = 'on'
```

修改主配置文件
```bash
vim /var/lib/pgsql/14/data/postgresql.conf
primary_conninfo = 'host=172.16.104.109 port=5432 user=replica password=replica'
recovery_target_timeline = latest
max_connections = 1200
hot_standby = on
max_standby_streaming_delay = 30s
wal_receiver_status_interval = 10s
hot_standby_feedback = on

# 如果是root用户进行操作，那么需要把data目录下的文件属主和属组全部改为postgres
chown -R postgres. /var/lib/pgsql/14/data/*

# 重启数据库
systemctl restart postgresql-14
```


# 5 安装postgis
```bash
# 安装postgis插件
sudo yum install epel-release
sudo yum install postgis33_14.x86_64

# 开启postgis扩展
sudo su - postgres
psql
postgres=# create extension postgis;
postgres=# create extension postgis_raster;
postgres=# create extension postgis_topology;
postgres=# create extension postgis_sfcgal;
postgres=# create extension fuzzystrmatch;
postgres=# create extension address_standardizer;
postgres=# create extension address_standardizer_data_us;
postgres=# create extension postgis_tiger_geocoder;

# 查看扩展
postgres=# \dx
```
