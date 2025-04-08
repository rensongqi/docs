
# 压缩nginx日志
```bash
cat <<EOF > /etc/logrotate.d/openresty
/var/log/openresty/*.log {
    create 0640 root root
    daily
    rotate 5
    dateext
    copytruncate
    compress
}
EOF
```

通配匹配子目录

```bash
cat <<EOF > /etc/logrotate.d/openresty
/var/log/*/*.log {
    create 0640 root root
    daily
    rotate 5
    dateext
    copytruncate
    compress
}
EOF
```

# 参考文章

[Logrotate滚动openresty日志](https://www.cnblogs.com/xiao987334176/p/11190837.html)