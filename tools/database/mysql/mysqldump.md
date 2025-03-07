# mysqldump备份及还原

## 备份

锁表导出指定库
```bash
mysqldump -h 172.16.1.104 -u root -p<PASSWORD> --max_allowed_packet=1024M --databases <DATABASE_NAME> > backup.sql 
```

不锁表导出指定库
```bash
mysqldump -h 172.16.1.104 -u root -p<PASSWORD> --max_allowed_packet=1024M --single-transaction --databases <DATABASE_NAME> > backup.sql
```

仅导出表结构
```bash
mysqldump -h 172.16.1.104 -P 13306 -u readonly --single-transaction --no-data -p promise > backup.sql
```

仅导出INSERT插入数据，不需要drop和create table的语句
```bash
mysqldump -h 172.16.1.104 -P 13306 -u readonly --single-transaction --skip-add-drop-table --no-create-info -p promise > backup.sql
```

不锁表导出指定库中指定表
```bash
mysqldump -h 172.16.1.104 -u root -p<PASSWORD> --tables test1 test2 --single-transaction --databases <DATABASE_NAME> > backup.sql
```

## 还原

```bash
# 方法 一
mysql -u root -p <DATABASE_NAME> < backup.sql

# 方法 二
mysql -u root -p
mysql> USE your_database;
mysql> SOURCE /path/to/backup.sql;
```