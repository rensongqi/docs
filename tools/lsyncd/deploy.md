
# 配置

```bash
yum install lsyncd

# 配置本地目录之间实时同步
settings {
    logfile = "/var/log/lsyncd/lsyncd.log",
    statusFile = "/var/log/lsyncd/lsyncd.status",
    inotifyMode = "CloseWrite",
    maxProcesses = 7,
    maxDelays = 200,
    insist = true
}

sync {
    default.rsync,
    source    = "/test1",
    target    = "/test2",
    rsync     = {
        binary    = "/bin/rsync",
        archive   = true,
        compress  = true,
        verbose   = true
    }
}
```