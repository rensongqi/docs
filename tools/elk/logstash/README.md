
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

# 日志解析案例

背景描述

> 会有源源不断的日志文件会生成并需要logstash解析文件到elasticsearch中，文件解析完成后需要删除该文件
> 
> 官方手册: [plugins-inputs-file](https://www.elastic.co/guide/en/logstash/8.17/plugins-inputs-file.html)

- `file_completed_action` 默认值为 `delete` ，需要`mode`的模式为`read`时才生效
- `file_completed_log_path` 记录删除的文件绝对路径，需要对该文件具有写入的权限

```
input {
    file{
        path => "/vehicle-stubs/\$Logger/20*/**/**/**/**/v*/*.log.*"
        exclude => "syslog.log.*"
        type => "stublogger"
        sincedb_path => "/dev/null"
        file_completed_action => "delete"
        file_completed_log_path => "/etc/logstash/delete.log"
        max_open_files => 10000
        mode => "read"
    }
}

input {
    file{
        path => "/vehicle-stubs/\$Logger/20*/**/**/**/**/v*/syslog.log.*"
        type => "syslog-stublogger"
        sincedb_path => "/dev/null"
        file_completed_action => "delete"
        max_open_files => 10000
        mode => "read"
    }
}

filter {
    if [type] == "stublogger" {
        ruby {
            code => "
                path = event.get('log')['file']['path']
                event.set('logfilepath', path.split('/'))
                logfileprefix = path.split('/')[9]
                prefix = logfileprefix.split('.').first(2).join('.')
                event.set('logfile', prefix) 
            "
        }
        multiline {
            pattern => "^(I|E|W|F|D)\d{6}"
            negate => true
            what => "previous"
        }
        grok {
            match => {
                "message" => "^(?<level>[IEWFD]).{0}(?<msg>.{19})"
            }
            tag_on_failure => ["grok_failure"]
        }
        mutate {
            gsub => [ "level", "I", "info" ]
            gsub => [ "level", "E", "error" ]
            gsub => [ "level", "W", "warn" ]
            gsub => [ "level", "F", "fatal" ]
            gsub => [ "level", "D", "debug" ]
            add_field => { 
                "logtime" => "20%{msg}"
                "carname" => "%{[logfilepath][6]}"
                "hostname" => "%{[logfilepath][7]}"
                "version" => "%{[logfilepath][8]}"
                "oss" => "http://ossapi.rsq.cn:9000/vehicle-stubs-bak/$Logger/%{[logfilepath][3]}/%{[logfilepath][4]}/%{[logfilepath][5]}/%{[logfilepath][6]}/%{[logfilepath][7]}/%{[logfilepath][8]}/%{[logfilepath][9]}"
            }
            convert => {
                "carname" => "string"
                "hostname" => "string"
                "logfile" => "string"
                "version" => "string"
            }
        }
        date {
            match => ["logtime","YYYYMMdd HH:mm:ss:SSS"]
            target => "@timestamp"
            remove_field => [ "logfilepath", "msg", "logtime", "event", "log", "host"] 
        }
    }

    if [type] == "syslog-stublogger" {
        ruby {
            code => "
                path = event.get('log')['file']['path']
                event.set('logfilepath', path.split('/'))
                logfileprefix = path.split('/')[9]
                prefix = logfileprefix.split('.').first(2).join('.')
                event.set('logfile', prefix) 
            "
        }
        grok {
            match => { "message" => "%{MONTH:log_month}%{SPACE}%{MONTHDAY:log_day}%{SPACE}%{TIME:log_time}" }
            tag_on_failure => ["grok_failure"]
        }
        mutate {
            add_field => { 
                "carname" => "%{[logfilepath][6]}"
                "hostname" => "%{[logfilepath][7]}"
                "version" => "%{[logfilepath][8]}"
                "oss" => "http://ossapi.rsq.cn:9000/vehicle-stubs-bak/$Logger/%{[logfilepath][3]}/%{[logfilepath][4]}/%{[logfilepath][5]}/%{[logfilepath][6]}/%{[logfilepath][7]}/%{[logfilepath][8]}/%{[logfilepath][9]}"
                "log_timestamp" => "%{log_month} %{log_day} %{log_time} 2024"
                "level" => "info"
            }
            convert => {
                "carname" => "string"
                "hostname" => "string"
                "logfile" => "string"
                "version" => "string"
            }
        }
        date {
            match => ["log_timestamp","MMM dd HH:mm:ss yyyy"]
            target => "@timestamp"
            remove_field => ["logfilepath", "log_month", "log_day", "log_time", "log_timestamp", "log", "event", "host", "msg"] 
            tag_on_failure => ["date_failure"]
        }
    }
    if "grok_failure" in [tags] {
        drop { }
    }
    if "_dateparsefailure" in [tags] {
        drop { }
    }
    if "date_failure" in [tags] {
        drop { }
    }

}

output {
    elasticsearch {
        hosts => ["172.16.100.107:9200", "172.16.100.108:9200", "172.16.100.109:9200"]
        index => "stublogger-%{+YYYY.MM.dd}"
        action => "create"
    }
}
```
