- [1 用户管理](#1-用户管理)
  - [1.1 角色](#11-角色)
  - [1.2 注意事项](#12-注意事项)
  - [1.3 给单个数据库创建用户](#13-给单个数据库创建用户)
  - [1.4 给一个用户授权多个数据库](#14-给一个用户授权多个数据库)
  - [1.5 其它命令](#15-其它命令)
  - [1.6 对比两个mongo重复的collections](#16-对比两个mongo重复的collections)
  - [1.7 mongo间数据迁移](#17-mongo间数据迁移)
- [参考文章](#参考文章)


# 1 用户管理
## 1.1 角色
1. 数据库用户角色：read、readWrite
2. 数据库管理角色：dbAdmin、dbOwner、userAdmin
3. 集群管理角色：clusterAdmin、clusterManager、clusterMonitor、hostManager
4. 备份恢复角色：backup、restore
5. 所有数据库角色： readAnyDatabase、readWriteAnyDatabase、userAdminAnyDatabase、
dbAdminAnyDatabase
6. 超级用户角色：root

## 1.2 注意事项
授权常见的模式分两种
1. 指定单个数据库给用户授权，此时需要先切换到指定数据库名称空间下，以test为例，use test，然后进行相关的授权grantRolesToUser，以rsq用户为例，第一次创建，则需要在此数据库中进行用户的创建，注意此用户跟其它数据库的用户不通用，如果需要在test2数据库中给rsq用户授权，则需要在test2中也创建rsq用户。此时mongo连接串需要指定相对应的数据库才可以登录，否则会报认证失败。
2. 指定多个数据库给一个用户授权，这种情况需要切换到admin数据库中，先对用户进行admin数据库的权限授权，如dbAdmin，授权完再给普通数据库授权，这种情况不需要再进入到普通数据库创建用户，用户是全局生效的。

## 1.3 给单个数据库创建用户
连接串：mongodb://<user>:<pass>@mongo.rsq.cn:27017/?tls=false&authSource=ai-rsq3d-xsl-v56

```bash
# 新建数据库
use ai-rsq3d-xsl-v56

# 新建一个测试collection
db.site.insert({"name":"ai-rsq3d-xsl-v56"})

# 查看数据库
show dbs

# 给新数据库创建用户，密码跟之前保持一致
db.createUser({user: "user",pwd: "xxxxx",roles: [ { role: "dbOwner", db: "ai-rsq3d-xsl-v56" } ]})

# 如果需要角色授权，则有如下命令
db.runCommand(
    {
        grantRolesToUser: "user",
        roles:
            [
                { role: "dbOwner", db: "ai-rsq3d-corolla-v56" }
            ]
    }
)
```

## 1.4 给一个用户授权多个数据库
连接串：mongodb://user:xxxxx@mongo.rsqrobot.cn:27017/?tls=false
```bash
# 需要进入到admin数据库授权多个数据库给同一个用户
use admin

# 创建用户
db.createUser({
    user: "user",
    pwd: "xxxxx",
    roles: [
        { 
            role: "dbAdmin", db: "admin" 
        },
        { 
            role: "dbOwner", db: "user" 
        },
        { 
            role: "dbOwner", db: "ai-data-cleaning-test" 
        },
        { 
            role: "dbOwner", db: "ai-rsq3d-x3s-v56" 
        },
        { 
            role: "dbOwner", db: "ai-rsq3d-corolla-v56" 
        },
        { 
            role: "dbOwner", db: "ai-rsq3d-xsl-v56" 
        } 
    ]
})
```

## 1.5 其它命令
```bash
# 查看db相关权限
db.runCommand(
    {
        rolesInfo: { role: "readWrite", db: "zgf-lidar" },
        showPrivileges: true
    }
)

# 查询用户
db.getUsers();
show users

# 删除用户
db.dropUser('admin')

# 修改用户密码
db.updateUser('admin', {pwd: '111111'})

# 创建管理员用户，需要使用admin数据库
use admin  # 
db.createUser({user: "rensongqi",pwd: "xxxxx",roles: [ { role: "dbAdmin", db: "perception" } ]})

# 回收用户角色
db.revokeRolesFromUser(
    "perception",
    [
        { role: "root", db: "admin" }
    ]
)

# 删除数据库
use percepai-rsq3d-corolla-v56tion
db.dropDatabase()

# 添加新角色
db.createRole({
  role: "aiPolicyRole",
  privileges: [
    {
        resource: { db: "", collection: "" },
        actions: [ "find", "insert", "update", "remove", "createCollection", "dropCollection" ],
        resource: { db: "ai-", collection: "" }
    }
  ],
  roles: []
})

# 删除角色
db.dropRole( "aiPolicyRole", { w: "majority" } ) 
```

## 1.6 对比两个mongo重复的collections
使用mongosh，[下载地址](https://www.mongodb.com/try/download/shell)

compare.sc
```
// 配置两个 MongoDB 集群的连接信息
const cluster1 = new Mongo("mongodb://user1:password1@host1:27017");
const cluster2 = new Mongo("mongodb://user2:password2@host2:27017");

// 获取所有集合的全路径列表
function getAllCollections(cluster) {
    const collections = [];
    const databases = cluster.getDB("admin").adminCommand("listDatabases").databases;

    databases.forEach((dbInfo) => {
        const dbName = dbInfo.name;
        const db = cluster.getDB(dbName);
        const collNames = db.getCollectionNames();

        collNames.forEach((collName) => {
            collections.push(`${dbName}.${collName}`);
        });
    });

    return collections;
}

// 获取两个集群的集合列表
const collectionsCluster1 = getAllCollections(cluster1);
const collectionsCluster2 = getAllCollections(cluster2);

// 比对两个集合列表，找出重复项
const duplicateCollections = collectionsCluster1.filter((coll) =>
    collectionsCluster2.includes(coll)
);

// 打印结果
if (duplicateCollections.length > 0) {
    print("重复的集合列表：");
    duplicateCollections.forEach((coll) => print(coll));
} else {
    print("没有发现重复的集合。");
}
```

执行脚本

```bash
mongosh -f compare.sc
```

## 1.7 mongo间数据迁移

需要使用`mongodump`和`mongorestore`，[下载地址](https://www.mongodb.com/try/download/database-tools)

collection.list
```
ai-3d-test
ai-2d-test
...
```

dump.sh
```bash
mkdir -p /data/mongo_dump
while read line; do
    mongodump -u root -p <pass> --authenticationDatabase admin -d $line -o ./dump/
done < collection.list
```

restore.sh
```bash
mongorestore  -u root -p <pass> --host mongo.rsq.cn:27017 --authenticationDatabase admin --dir=/data/mongo_dump --drop
```

# 参考文章
- [mongodb操作指令](https://www.cnblogs.com/dbabd/p/10811523.html)
- [manage-users-and-roles](https://www.yiibai.com/mongodb/manage-users-and-roles.html)
- [MongoDB 4.X 用户和角色权限管理总结](https://www.cnblogs.com/dbabd/p/10811523.html)