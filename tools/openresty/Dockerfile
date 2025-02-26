FROM centos:7

RUN yum install -y yum-utils \
    && yum groupinstall "Development Tools" \
    && yum install -y pcre-devel openssl-devel gcc gcc-c++ curl zlib-devel autoconf automake \
        gettext \
        gzip \
        make \
        tar \
        unzip \
        perl-IPC-Cmd \
    && cd /usr/local/src \
    && curl -o openresty-1.19.3.1.tar.gz https://openresty.org/download/openresty-1.19.3.1.tar.gz \
    && curl -o ngx_cache_purge-2.3.tar.gz http://labs.frickle.com/files/ngx_cache_purge-2.3.tar.gz \
    && curl -o openssl-3.0.1.tar.gz https://www.openssl.org/source/openssl-3.0.1.tar.gz \
    && tar xzf openresty-1.19.3.1.tar.gz \
    && tar xzf ngx_cache_purge-2.3.tar.gz \
    && tar xzf openssl-3.0.1.tar.gz \
    && cd openssl-3.0.1 \
    && ./config \
        --prefix=/usr/local/openssl \
        --openssldir=/usr/local/openssl shared zlib \
    && make -j$(nproc) 
    && make install
    && cd ../openresty-1.19.3.1 \
    && ./configure \
        --prefix=/usr/local/openresty-1.19.3.1 \
        --with-http_ssl_module \
        --with-http_stub_status_module \
        --with-http_slice_module \
        --with-http_gzip_static_module \
        --with-pcre \
        --with-http_v2_module \
        --with-http_mp4_module \
        --with-http_flv_module \
        --with-http_secure_link_module  \
        --with-http_auth_request_module \
        --with-http_addition_module \
        --add-module=../ngx_cache_purge-2.3 \
        --with-openssl=/usr/local/src/openssl-3.0.1 \
    && gmake -j 4 \
    && gmake install \
    && yum remove -y make \
    && yum clean all \
    && ln -sf /usr/local/openresty-1.19.3.1/nginx /usr/local/nginx \
    && rm -f  /usr/bin/nginx /usr/sbin/nginx \
    && ln -sf /usr/local/openresty-1.19.3.1/nginx/sbin/nginx /usr/bin/nginx \
    && ln -sf /usr/local/openresty-1.19.3.1/nginx/sbin/nginx /usr/sbin/nginx \
    && mkdir -p /var/log/nginx/ /usr/local/nginx/conf/ssl /var/log/openresty/


ENV PATH=$PATH:/usr/local/openresty-1.19.3.1/luajit/bin:/usr/local/openresty-1.19.3.1/nginx/sbin:/usr/local/openresty-1.19.3.1/bin

ENV LUA_PATH="/usr/local/openresty-1.19.3.1/site/lualib/?.ljbc;/usr/local/openresty-1.19.3.1/site/lualib/?/init.ljbc;/usr/local/openresty-1.19.3.1/lualib/?.ljbc;/usr/local/openresty-1.19.3.1/lualib/?/init.ljbc;/usr/local/openresty-1.19.3.1/site/lualib/?.lua;/usr/local/openresty-1.19.3.1/site/lualib/?/init.lua;/usr/local/openresty-1.19.3.1/lualib/?.lua;/usr/local/openresty-1.19.3.1/lualib/?/init.lua;./?.lua;/usr/local/openresty-1.19.3.1/luajit/share/luajit-2.1.0-beta3/?.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua;/usr/local/openresty-1.19.3.1/luajit/share/lua/5.1/?.lua;/usr/local/openresty-1.19.3.1/luajit/share/lua/5.1/?/init.lua"

ENV LUA_CPATH="/usr/local/openresty-1.19.3.1/site/lualib/?.so;/usr/local/openresty-1.19.3.1/lualib/?.so;./?.so;/usr/local/lib/lua/5.1/?.so;/usr/local/openresty-1.19.3.1/luajit/lib/lua/5.1/?.so;/usr/local/lib/lua/5.1/loadall.so;/usr/local/openresty-1.19.3.1/luajit/lib/lua/5.1/?.so"

# Copy nginx configuration files
COPY nginx.conf /usr/local/openresty-1.19.3.1/nginx/conf/nginx.conf
EXPOSE 80

ENTRYPOINT ["/usr/bin/nginx", "-g", "daemon off;"]