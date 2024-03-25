
# 背景

需要解析生成的日志文件到elasticsearch中

```
input {
    file{
        path => "/vehicle-stubs/\$Logger/20*/**/**/**/**/v*/*.log.*"
        exclude => "syslog.log.*"
        type => "stublogger"
        start_position => "beginning"
        sincedb_path => "/dev/null"
        file_completed_action => "delete"
        max_open_files => 120000
        close_older => 600
    }
}

input {
    file{
        path => "/vehicle-stubs/\$Logger/20*/**/**/**/**/v*/syslog.log.*"
        type => "syslog-stublogger"
        start_position => "beginning"
        sincedb_path => "/dev/null"
        file_completed_action => "delete"
        max_open_files => 120000
        close_older => 600
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

定时获取日志，定时清理日志

```bash
58 * * * * curl -X POST -d '{"bucket_name": "vehicle-stubs","filter": "$Logger/20"}' http://127.0.0.1:8893/api/oss/get

*/10 * * * * find /vehicle-stubs/\$Logger/ -type f -cmin +60 -exec rm -f -r {} \; >>/dev/null 2>&1
```