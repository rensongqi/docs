# check user pass

对Basic auth的用户和密码进行校验
```bash
curl https://<user>:<pass>@test.rsq.com/proxy
```

# 使用方法
```
location /ubuntu {
    access_by_lua_file /usr/local/openresty/site/lualib/check_user_pass/check.lua;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_pass http://172.16.1.108:9999/xxx;
}
```