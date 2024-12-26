# Mongo Cluster部署

# K8s
mongodb-rs.yaml
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mongodb-rs-cm
  namespace: mongo
data:
  keyfile: |
        dGhpcyBpcyBycyBzdXBlciBzZWNyZXQga2V5Cg==
  mongod_rs.conf: |+
    systemLog:
      destination: file
      logAppend: true
      path: /data/mongod.log
    storage:
      dbPath: /data
      journal:
        enabled: true
      directoryPerDB: true
      wiredTiger:
        engineConfig:
          cacheSizeGB: 2
          directoryForIndexes: true
    processManagement:
      fork: true
      pidFilePath: /data/mongod.pid
    net:
      port: 27017
      bindIp: 0.0.0.0
      maxIncomingConnections: 5000
    security:
      keyFile: /data/configdb/keyfile
      authorization: enabled
    replication:
      oplogSizeMB: 1024
      replSetName: rs0    
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mongodb-rs
  namespace: mongo
spec:
  serviceName: mongodb-rs
  replicas: 3
  selector:
    matchLabels:
      app: mongodb-rs
  template:
    metadata:
      labels:
        app: mongodb-rs
    spec:
      containers:
      - name: mongo
        image: harbor.rsq.cn/mongo/mongo:4.4.1
        ports:
        - containerPort: 27017
          name: client
        command: ["sh"]
        args:
        - "-c"
        - |
          set -ex
          mongod --config /data/configdb/mongod_rs.conf
          sleep infinity              
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        volumeMounts:
        - name: conf
          mountPath: /data/configdb
          readOnly: false
        - name: data
          mountPath: /data
          readOnly: false
      volumes:
      - name: conf
        configMap:
          name: mongodb-rs-cm
          defaultMode: 0600
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 100Gi
      storageClassName: local-storage
---
apiVersion: v1
kind: Service
metadata:
  name: mongodb-rs
  labels:
    app: mongodb-rs
  namespace: mongo
spec:
  ports:
    - port: 27017
      targetPort: 27017
      nodePort: 30717
  selector:
    app: mongodb-rs
  type: NodePort
```
mongodb-pv.yaml
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv.mongodb.cluster01
spec:
  capacity:
    storage: 200Gi
  volumeMode: Filesystem
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: local-storage
  local:
    path: /data/mongodb/data
  nodeAffinity:
     required:
       nodeSelectorTerms:
       - matchExpressions:
         - key: kubernetes.io/hostname
           operator: In
           values:
           - master1
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv.mongodb.cluster02
spec:
  capacity:
    storage: 200Gi
  volumeMode: Filesystem
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: local-storage
  local:
    path: /data/mongodb/data
  nodeAffinity:
     required:
       nodeSelectorTerms:
       - matchExpressions:
         - key: kubernetes.io/hostname
           operator: In
           values:
           - master2
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv.mongodb.cluster03
spec:
  capacity:
    storage: 200Gi
  volumeMode: Filesystem
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: local-storage
  local:
    path: /data/mongodb/data
  nodeAffinity:
     required:
       nodeSelectorTerms:
       - matchExpressions:
         - key: kubernetes.io/hostname
           operator: In
           values:
           - master3
```

# docker-compose
创建mongo-keyfile

```bash
mkdir /data/mongodb/mongo/{data,configdb} -p
cd /data/mongodb/mongo

openssl rand -base64 745 > mongo-keyfile
echo "rsq" | sudo -S chmod 600 mongo-keyfile
echo "rsq" | sudo -S chown 999 mongo-keyfile
```
docker-compose.yml
```bash
version: "3"

services:
  #主节点
  mongodb1:
    image: harbor.rsq.cn/library/mongo/mongo:5.0.6
    container_name: mongo1
    restart: always
    ports:
      - 27017:27017
    environment:
      - MONGO_INITDB_ROOT_USERNAME=root
      - MONGO_INITDB_ROOT_PASSWORD=123456
    command: mongod --replSet rs0 --keyFile /mongo-keyfile
    volumes:
      - /etc/localtime:/etc/localtime
      - /data/mongodb/mongo/data:/data/db
      - /data/mongodb/mongo/configdb:/data/configdb
      - /data/mongodb/mongo/mongo-keyfile:/mongo-keyfile
    entrypoint:
      - bash
      - -c
      - |
        chmod 400 /mongo-keyfile
        chown 999:999 /mongo-keyfile
        exec docker-entrypoint.sh $$@
  # 副节点
  mongodb2:
    image: harbor.rsq.cn/library/mongo/mongo:5.0.6
    container_name: mongo2
    restart: always
    ports:
      - 27017:27017
    environment:
      - MONGO_INITDB_ROOT_USERNAME=root
      - MONGO_INITDB_ROOT_PASSWORD=123456
    command: mongod --replSet rs0 --keyFile /mongo-keyfile
    volumes:
      - /etc/localtime:/etc/localtime
      - /data/mongodb/mongo/data:/data/db
      - /data/mongodb/mongo/configdb:/data/configdb
      - /data/mongodb/mongo/mongo-keyfile:/mongo-keyfile
    entrypoint:
      - bash
      - -c
      - |
        chmod 400 /mongo-keyfile
        chown 999:999 /mongo-keyfile
        exec docker-entrypoint.sh $$@
  # 副节点
  mongodb3:
    image: harbor.rsq.cn/library/mongo/mongo:5.0.6
    container_name: mongo3
    restart: always
    ports:
      - 27017:27017
    environment:
      - MONGO_INITDB_ROOT_USERNAME=root
      - MONGO_INITDB_ROOT_PASSWORD=123456
    command: mongod --replSet rs0 --keyFile /mongo-keyfile
    volumes:
      - /etc/localtime:/etc/localtime
      - /data/mongodb/mongo/data:/data/db
      - /data/mongodb/mongo/configdb:/data/configdb
      - /data/mongodb/mongo/mongo-keyfile:/mongo-keyfile
    entrypoint:
      - bash
      - -c
      - |
        chmod 400 /mongo-keyfile
        chown 999:999 /mongo-keyfile
        exec docker-entrypoint.sh $$@
```

在不同节点执行如下命令启动容器
```bash
# node1
docker-compose up -d mongodb1
# node2
docker-compose up -d mongodb2
# node3
docker-compose up -d mongodb3
```

初始化mongodb cluster
```bash
# 1.进入容器
docker exec -it mongo1 bash

# 2.进入数据库
mongo -u root -p 123456

# 3.集群init
config={_id:"rs0",members:[ 
    {_id:0,host:"172.16.100.107:27017"}, 
    {_id:1,host:"172.16.100.108:27017"}, 
    {_id:2,host:"172.16.100.109:27017"}] 
}

rs.initiate(config);

# 4.查看集群状态
rs.status();

# 5.增加节点的权重
cfg = rs.conf()
# 修改权重
cfg.members[0].priority=5
cfg.members[1].priority=3
# 从新配置
rs.reconfig(cfg)
```