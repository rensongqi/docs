local http = require("resty.http")
local json = require("cjson")

-- 定义共享内存字典的名称和大小
local b_cache = "b_cache"
local c_cache = "c_cache"
local token_size = 1024 
local auth_header = ngx.var.http_authorization

local function http_post(url, body, headers)
    local request = http.new()
    request:set_timeout(10000)

    return request:request_uri(url, {
        method = "POST",
        body = body,
        headers = headers,
        ssl_verify = false
    })
end

local function check_pass(token)
    -- 如果 Token 是以 C 开头，则进行验证
    if token and string.sub(token, 1, 1) == "C" then
        -- 从共享内存字典中获取 Token
        ngx.log(ngx.INFO, "Token: ", token)
        local cache = ngx.shared[c_cache]
        local cached_value, err = cache:get(token)

        if cached_value then
            -- Token 在缓存中存在，直接放行
            ngx.log(ngx.INFO, "Token found in cache. Allowing access.")
        else
            -- Token 不在缓存中，进行验证
            local body = { token = token }
            local json_data = json.encode(body)
            local url = "https://check.rsq.com/api/v1/access/auth"
            local headers = {
                ["Content-Type"] = "application/json; charset=utf-8",
            }
            local res, err = http_post(url, json_data, headers)
            ngx.log(ngx.ERR, err)
            ngx.log(ngx.ERR, res.status)

            if res and res.status == 200 then
                -- 认证通过，将 Token 存入缓存
                cache:set(token, true)
                ngx.log(ngx.INFO, "Token verification successful. Caching token.")
            else
                -- 认证失败，返回 403 错误
                ngx.log(ngx.ERR, "Token verification failed")
                ngx.exit(ngx.HTTP_FORBIDDEN)
            end
        end
    else
        ngx.exit(ngx.HTTP_FORBIDDEN)
    end
end

if auth_header then
    local _, _, encoded_creds = string.find(auth_header, "Basic%s+(.+)")
    if encoded_creds then
        local decoded_creds = ngx.decode_base64(encoded_creds)
        local username, password = decoded_creds:match("(.*):(.*)")
        if username == "rsq" then
            check_pass(password)
            return
        else
            ngx.say("Authentication failed, user name is invalid.")
            ngx.exit(ngx.HTTP_UNAUTHORIZED)
        end
    end
end

-- default deny all requests
ngx.exit(ngx.HTTP_UNAUTHORIZED)