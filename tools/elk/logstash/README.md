
- [修改logstash配置](#修改logstash配置)
- [日志解析案例](#日志解析案例)

# 修改logstash配置

1. 修改启动内存参数

    `/etc/logstash/jvm.options` 
    ```bash
    -Xms3g
    -Xmx3g
    ```

2. 修改文件打开数量限制
   
    `/etc/logstash/startup.options`
    ```
    LS_OPEN_FILES=163840
    ```

    `/lib/systemd/system/logstash.service`
    ```
    LimitNOFILE=163840
    ```

    操作系统层面
    ```bash
    cat >>/etc/security/limits.conf<<EOF
    * soft nofile 655360
    * hard nofile 231072
    * soft nproc 655360
    * hard nproc 655360
    * soft memlock unlimited
    * hard memlock unlimited
    EOF

    cat >>/etc/sysctl.conf<<EOF
    fs.file-max = 4194303
    vm.max_map_count=262144
    fs.inotify.max_user_watches=524288
    EOF
    sysctl -p

    # centos
    sed -i 's/4096/unlimited/g' /etc/security/limits.d/20-nproc.conf

    # ubuntu
    echo "* - nofile 102400" >> /etc/security/limits.d/nofile.conf
    echo "root - nofile 102400" >> /etc/security/limits.d/nofile.conf
    ```

3. 安装插件

    ```bash
    /usr/share/logstash/bin/logstash-plugin install logstash-filter-multiline
    ```

4. 升级logstash

    > [官方下载地址](https://www.elastic.co/downloads/past-releases/logstash-8-17-2)

    ```bash
    yum install logstash-8.17.2-x86_64.rpm
    ```

# 日志解析案例

背景描述

> 会有源源不断的日志文件会生成并需要logstash解析文件到elasticsearch中，文件解析完成后需要删除该文件
> 
> 官方手册: [plugins-inputs-file](https://www.elastic.co/guide/en/logstash/8.17/plugins-inputs-file.html)

- `file_completed_action` 默认值为 `delete` ，需要`mode`的模式为`read`时才生效
- `file_completed_log_path` 记录删除的文件绝对路径，需要对该文件具有写入的权限

配置文件参考

初始化
```bash
# 创建日志输出目录
mkdir /var/log/carlog/
chmdo 777 /var/log/carlog/

# 修改新写入日志文件权限，或者将logstash启动用户改为root，否则logstash将无权限读完文件后删除日志文件
vim /lib/systemd/system/logstash.service
[Service]
User=root
Group=root
```

[logstash.yml](./logstash.yml)
