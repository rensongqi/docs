# Jenkins配置


## 基础配置

### 安装Jenkins

```bash
docker run -d \
    --name jenkins \
    -p 8080:8080 \
    -p 50000:50000 \
    --restart=on-failure \
    -v jenkins_home:/var/jenkins_home \
    jenkins/jenkins:2.479.3-lts-jdk17
```

修改插件下载为清华源

```bash
docker inspect volume jenkins_home

sed -i 's#updates.jenkins.io/download#mirrors.tuna.tsinghua.edu.cn/jenkins#g' /var/lib/docker/volumes/jenkins_home/_data

docker restart jenkins
```

### 插件安装
- SSH
- pipeline
- LDAP
- Vault

### Config LDAP

```
Server: ldap://RSQ.cn:389
root DN: DC=RSQ,DC=CN
User search base: CN=Users,DC=RSQ,DC=CN
User search filter: cn={0}
Manager DN: CN=devops,CN=Users,DC=RSQ,DC=CN
Manager Password: xxxxxx
Display Name LDAP attribute: displayname 或 sAMAccountName
```

## 从vault中获取密码
```pipeline
def secrets = [
    [
        path: 'devops/jenkins/gitlab', engineVersion: 2, secretValues: [
            [
                envVar: 'USERNAME', vaultKey: 'username'
            ],
            [
                envVar: 'PASSWORD', vaultKey: 'password'
            ]
        ]
    ]
]

def configuration = [
    vaultUrl: 'http://172.16.104.10:8200',
    vaultCredentialId: 'vault-app-role',
    engineVersion: 2
]

pipeline {
    agent any
    stages {
        stage("Test"){
            steps{
                script{
                    withVault([configuration: configuration, vaultSecrets: secrets]) {
                        sh "echo ${USERNAME}"
                        sh "echo ${PASSWORD}"
                    }
                }
            }
        }
    }
}
```