# 1 准备工作
## 1.1 需要两个仓库
一个仓库放pipeline，另一个仓库放代码库

| git仓库         | 描述                      |
| --------------- | ------------------------- |
| devops_pipeline | Sharable library pipeline |
| devops          | 代码                      |

**devops_pipeline 目录结构：**
```bash
$ ls -l
total 5
-rw-r--r-- 1 rsq rsq  21 Oct 24 13:48 README.md
-rwxr-xr-x 1 rsq rsq 595 Oct 24 14:20 test.groovy*
drwxr-xr-x 1 rsq rsq   0 Oct 24 13:56 vars/
$ ls -l vars
total 2
-rwxr-xr-x 1 rsq rsq 178 Oct 24 14:17 build.groovy*
-rwxr-xr-x 1 rsq rsq 819 Oct 24 14:07 gitCheckout.groovy*
```
**devops目录结构：**
```bash
$ ls -l
total 2
-rw-r--r-- 1 rsq rsq Oct 24 13:41 README.md
-rw-r--r-- 1 rsq rsq 29 Oct 24 13:52 build.sh
$ cat build.sh
echo "This is a test scripts"
```

## 1.2 Jenkins 配置Share Library
Manager Jnekins --> Configure System --> Global Pipeline Libraries ---> Add Libary(选择devops_pipeline)
![在这里插入图片描述](https://img-blog.csdnimg.cn/cc3baa40e34644238c514793a145fb9b.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBAUlNR5Y2a5a6i,size_20,color_FFFFFF,t_70,g_se,x_16)

## 1.3 相关groovy

**gitCheckout.groovy**
```groovy
#!/usr/bin/env groovy
def call(String gitRepo, String gitBranch, String gitCredential) {
    checkout(
        [   
            $class: 'GitSCM', 
            branches: [[name: gitBranch]], 
            doGenerateSubmoduleConfigurations: false, 
            extensions: [
                [$class: 'CheckoutOption', timeout: 60],
                [$class: 'GitLFSPull'],
                [$class: 'CloneOption', noTags: false, reference: '', shallow: false, timeout: 60],
                [$class: 'SubmoduleOption', timeout: 60, disableSubmodules: false, parentCredentials: true, recursiveSubmodules: true, reference: '', trackingSubmodules: false]
            ], 
            submoduleCfg: [], 
            userRemoteConfigs: [[credentialsId: gitCredential, 
            url: gitRepo]]
        ]
    )
}
```


**build.groovy**
```groovy
#!/usr/bin/env groovy
def call() {
    dir("${WORKSPACE}") {
        sh '''
            echo "This is a test scripts."
            echo "${WROKSPACE}"
        '''
    }
}
```

test.groovy 主要实现跟share library联动
**test.groovy**
```groovy
#!/usr/bin/env groovy
// share library entrypoint.
@Library("devops") _
node("devops_linux") {
    // define const
    def gitRepo = "${env.gitRepo}"
    def gitBranch = "${env.gitBranch}"
    def gitCredential = "${env.gitCredential}"
    // custom workspace
    env.WORKSPACE = "/home/devops/test"

    stage("Clean Env") {
        sh '''
            echo "---------------------"
            echo ${WORKSPACE}
        '''
    }

    stage("Git Clone") {
        gitCheckout(gitRepo, gitBranch, gitCredential)
    }

    stage("Build") {
        build()
    }
}
```
# 2 Jenkins 新建流水线 item 测试 
![在这里插入图片描述](https://img-blog.csdnimg.cn/4a8e9b82cf704850b3317fc41c89d7c0.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBAUlNR5Y2a5a6i,size_20,color_FFFFFF,t_70,g_se,x_16)

Configure item，pipeline SCM选择`devops_pipeline`这个仓库，Script Path选择 `test.groovy`
![在这里插入图片描述](https://img-blog.csdnimg.cn/ebf596db24994a2da17aa578572e6f72.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBAUlNR5Y2a5a6i,size_20,color_FFFFFF,t_70,g_se,x_16)