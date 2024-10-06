
# 1 部署
500人以内 4c8g
```bash
mkdir /data/gitlab/{data,config,logs} -p
cd /data/gitlab/

vim docker-compose.yml
version: '3'
services:
  web:
    image: 'harbor.rsq.cn/library/gitlab-ce:15.1.1-ce.0'
    restart: always
    environment:
      TZ: 'Asia/Shanghai'
      GITLAB_OMNIBUS_CONFIG: |
              external_url 'http://172.16.0.32'
              gitlab_rails['gitlab_shell_ssh_port'] = '10080'
              gitlab_rails['time_zone'] = 'Asia/Shanghai'
    ports:
      - '80:80'
      - '443:443'
      - '10080:22'
    volumes:
      - './config:/etc/gitlab'
      - './logs:/var/log/gitlab'
      - './data:/var/opt/gitlab'
    shm_size: '256m'

# 启动
docker-compose up -d
```

# 2 配置LDAP

编辑 `/data/gitlab/config/gitlab.rb`
```bash
external_url 'https://gitlab.RSQ.cn'
nginx['redirect_http_to_https'] = true
nginx['ssl_certificate'] = "/etc/gitlab/ssl/8444378__RSQ.cn.pem"
nginx['ssl_certificate_key'] = "/etc/gitlab/ssl/8444378__RSQ.cn.key"
gitlab_rails['gitlab_shell_ssh_port'] = '10080'
gitlab_rails['time_zone'] = 'Asia/Shanghai'
# Disable the bundled Omnibus provided PostgreSQL
postgresql['enable'] = false

# PostgreSQL connection details
gitlab_rails['db_adapter'] = 'postgresql'
gitlab_rails['db_encoding'] = 'unicode'
gitlab_rails['db_database'] = "postgres"
gitlab_rails['db_host'] = '172.16.100.21'
gitlab_rails['db_port'] = '5432'
gitlab_rails['db_username'] = 'postgres'
gitlab_rails['db_password'] = '123456'

# Gitlab配置AD，使用sAMAccountName登录
gitlab_rails['ldap_enabled'] = true
gitlab_rails['ldap_servers'] = YAML.load <<-'EOS'
   main: # 'main' is the GitLab 'provider ID' of this LDAP server
     label: 'LDAP'
     host: 'RSQ.cn'
     port: 389
     uid: 'sAMAccountName'
     bind_dn: 'CN=gitlab,CN=Users,DC=RSQ,DC=CN'
     password: 'xxxxxx'
     encryption: 'plain'
     active_directory: true
     allow_username_or_email_login: false
     block_auto_created_users: false
     base: 'DC=RSQ,DC=CN'
     user_filter: ''
     attributes:
       username: ['uid', 'sAMAccountName']
       email:    ['mail']
       name:     'sAMAccountName'
       first_name: 'givenName'
       last_name:  'sn'
EOS

# 进入容器中执行如下命令reload配置
gitlab-ctl reconfigure
```

# 3 修改root密码

```bash
# 进入到容器中
gitlab-rails console
u=User.find(1)

# root用户密码设置为root123456
u.password='root123456'

# 确认密码
u.password_confirmation = 'root123456'

# 保存配置并退出控制台
u.save!
exit
```