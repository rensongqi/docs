
- [1 部署canal.deployer](#1-部署canaldeployer)
- [2 部署canal.adapter](#2-部署canaladapter)
- [参考文章](#参考文章)

基于Canal实现mysql到mysql之间的数据实时备份

源数据库账号需要如下权限
```bash
SELECT,REPLICATION CLIENT,SUPER, REPLICATION SLAVE,PROCESS
```

# 1 部署canal.deployer

[canal.deployer下载地址](https://github.com/alibaba/canal/releases/download/canal-1.1.8/canal.deployer-1.1.8.tar.gz)

解压后进入`canal.deployer`目录中
```bash
# vim conf/example/instance.properties
canal.instance.master.address=172.16.100.100:13306 # 源数据库地址
canal.instance.dbUsername=canal                  # 源数据库账号
canal.instance.dbPassword=canal                  # 源数据库密码
canal.instance.filter.regex=backup_db\\..*         # 要同步的库，若要同步所有，可以默认*.\\..*
```
启动 canal.deployer
```bash
./bin/startup.sh
```
查看日志
```bash
tail logs/example/example.log
```

# 2 部署canal.adapter

[canal.adapter下载地址](https://github.com/alibaba/canal/releases/download/canal-1.1.8/canal.adapter-1.1.8.tar.gz)

解压后进入`canal.adapter`目录中

注释 `conf/bootstrap.yml` 暂时用不到
```bash
# vim conf/bootstrap.yml
#canal:
#  manager:
#    jdbc:
#      url: jdbc:mysql://127.0.0.1:3306/canal_manager?useUnicode=true&characterEncoding=UTF-8
#      username: canal
#      password: canal
```
编辑 `conf/application.yml` ，修改需要同步到的目的数据库信息
```bash
server:
  port: 8081
spring:
  jackson:
    date-format: yyyy-MM-dd HH:mm:ss
    time-zone: GMT+8
    default-property-inclusion: non_null
canal.conf:
  mode: tcp #tcp kafka rocketMQ rabbitMQ
  flatMessage: true
  zookeeperHosts:
  syncBatchSize: 1000
  retries: -1
  timeout:
  accessKey:
  secretKey:
  consumerProperties:
    # canal tcp consumer 指定deployer地址
    canal.tcp.server.host: 172.16.100.100:11111
    canal.tcp.zookeeper.hosts:
    canal.tcp.batch.size: 500
    canal.tcp.username:
    canal.tcp.password:
  srcDataSources:
    defaultDS:
      # 指定源数据库地址
      url: jdbc:mysql://172.16.100.100:13306/backup_db?useUnicode=true&characterEncoding=utf8&autoReconnect=true&useSSL=false
      username: canal    # 源库用户名
      password: canal    # 源库密码
      maxActive: 100
  canalAdapters:
  - instance: example # canal instance Name or mq topic name
    groups:
    - groupId: g1
      outerAdapters:
      - name: logger
      - name: rdb
        key: mysql1
        properties:
          jdbc.driverClassName: com.mysql.jdbc.Driver
          # 目的库链接信息
          jdbc.url: jdbc:mysql://172.16.100.101:13306/backup_db?useUnicode=true&characterEncoding=utf8&autoReconnect=true&useSSL=false
          jdbc.username: canal  # 目的库账号
          jdbc.password: canal  # 目的库密码
          druid.stat.enable: false
          druid.stat.slowSqlMillis: 1000
```

编辑自定义的rdb配置清单 `conf/rdb/mytest_user.yml`
```bash
# vim conf/rdb/mytest_user.yml
############ 如下配置用于同步单个表和指定的字段 ##########
#dataSourceKey: defaultDS
#destination: example
#groupId: g1
#outerAdapterKey: mysql1
#concurrent: true
#dbMapping:
#  database: mytest
#  table: user
#  targetTable: mytest2.user
#  targetPk:
#    id: id
#  mapAll: true
#  targetColumns:
#    id:
#    name:
#    role_id:
#    c_time:
#    test1:
#  etlCondition: "where c_time>={}"
#  commitBatch: 3000 # 批量提交的大小
############ 如下配置用于同步指定数据库中所有表 ##########
## Mirror schema synchronize config
dataSourceKey: defaultDS
destination: example
groupId: g1
outerAdapterKey: mysql1
concurrent: true
dbMapping:
  mirrorDb: true
  database: backup_db
```
启动 canal.adapter
```bash
./bin/startup.sh
```
查看日志
```bash
tail logs/adapter/adapter.log
```

# 参考文章

- [Canal同步两个mysql库](https://blog.csdn.net/sszdzq/article/details/137463288)
- [Canal同步mysql至mysql](https://www.stonewu.com/archives/canal-synchronization-problem)