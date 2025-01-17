# mysqldump备份及还原

## 备份

锁表
```bash
mysqldump -h 172.16.1.104 -u root -p<PASSWORD> --max_allowed_packet=1024M --databases <DATABASE_NAME> > backup.sql 
```

不锁表
```bash
mysqldump -h 172.16.1.104 -u root -p<PASSWORD> --max_allowed_packet=1024M --single-transaction --databases <DATABASE_NAME> > backup.sql
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