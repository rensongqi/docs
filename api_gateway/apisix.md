Centos7 安装APISIX 

[各类操作系统依赖安装](https://github.com/apache/apisix/blob/master/docs/en/latest/install-dependencies.md)
[官方安装手册](https://github.com/apache/apisix/blob/master/docs/en/latest/how-to-build.md)

## 1 安装

### 1.1 Install dependencies

```bash
# install etcd
wget https://github.com/etcd-io/etcd/releases/download/v3.4.13/etcd-v3.4.13-linux-amd64.tar.gz
tar -xvf etcd-v3.4.13-linux-amd64.tar.gz && \
    cd etcd-v3.4.13-linux-amd64 && \
    sudo cp -a etcd etcdctl /usr/bin/

# add OpenResty source
sudo yum install yum-utils
sudo yum-config-manager --add-repo https://openresty.org/package/centos/openresty.repo

# install OpenResty and some compilation tools
sudo yum install -y openresty curl git gcc openresty-openssl111-devel unzip

# install LuaRocks
curl https://raw.githubusercontent.com/apache/apisix/master/utils/linux-install-luarocks.sh -sL | bash -

# start etcd server
nohup etcd &
```

### 1.2 Install Apache APISIX

**RPM**

```bash
sudo yum install -y https://github.com/apache/apisix/releases/download/2.7/apisix-2.7-0.x86_64.rpm
```

**Docker**

```
https://hub.docker.com/r/apache/apisix
```

**Helm Chart**

```
https://github.com/apache/apisix-helm-chart
```

### 1.3 Manage Apache APISIX Server

We can initialize dependencies, start service, and stop service with commands in the Apache APISIX directory, we can also view all commands and their corresponding functions with the `make help` command.

### Initializing Dependencies

Run the following command to initialize the NGINX configuration file and etcd.

```shell
# initialize NGINX config file and etcd
make init
```

### Start Apache APISIX

Run the following command to start Apache APISIX.

```shell
# start Apache APISIX server
make run
```

### Stop Apache APISIX

Both `make quit` and `make stop` can stop Apache APISIX. The main difference is that `make quit` stops Apache APISIX gracefully, while `make stop` stops Apache APISIX immediately.

It is recommended to use gracefully stop command `make quit` because it ensures that Apache APISIX will complete all the requests it has received before stopping down. In contrast, `make stop` will trigger a forced shutdown, it stops Apache APISIX immediately, in which case the incoming requests will not be processed before the shutdown.

The command to perform a graceful shutdown is shown below.

```shell
# stop Apache APISIX server gracefully
make quit
```

The command to perform a forced shutdown is shown below.

```shell
# stop Apache APISIX server immediately
make stop
```

### View Other Operations

Run the `make help` command to see the returned results and get commands and descriptions of other operations.

```shell
# more actions find by `help`
make help
```

## Step 4: Run Test Cases

1. Install `cpanminus`, the package manager for `perl`.

2. Then install the test-nginx dependencies via `cpanm`:

  ```shell
  sudo cpanm --notest Test::Nginx IPC::Run > build.log 2>&1 || (cat build.log && exit 1)
  ```

3. Run the `git clone` command to clone the latest source code locally, please use the version we forked out：

  ```shell
  git clone https://github.com/iresty/test-nginx.git
  ```

4. Load the test-nginx library with the `prove` command in `perl` and run the test case set in the `/t` directory.

  - Append the current directory to the perl module directory: `export PERL5LIB=.:$PERL5LIB`, then run `make test` command.

  - Or you can specify the NGINX binary path by running this command: `TEST_NGINX_BINARY=/usr/local/bin/openresty prove -Itest-nginx/lib -r t`.

  :::note Note
  Some of the tests rely on external services and system configuration modification. For a complete test environment build, you can refer to `ci/linux_openresty_common_runner.sh`.
  :::

### Troubleshoot Testing

**Configuring NGINX Path**

The solution to the `Error unknown directive "lua_package_path" in /API_ASPIX/apisix/t/servroot/conf/nginx.conf` error is as shown below.

Ensure that Openresty is set to the default NGINX, and export the path as follows:

* `export PATH=/usr/local/openresty/nginx/sbin:$PATH`
  * Linux default installation path:
    * `export PATH=/usr/local/openresty/nginx/sbin:$PATH`
  * MacOS default installation path via homebrew:
    * `export PATH=/usr/local/opt/openresty/nginx/sbin:$PATH`

**Run a Single Test Case**

Run the specified test case using the following command.

```shell
prove -Itest-nginx/lib -r t/plugin/openid-connect.t
```

## Step 5: Update Admin API token to Protect Apache APISIX

You need to modify the Admin API key to protect Apache APISIX.

Please modify `apisix.admin_key` in `conf/config.yaml` and restart the service as shown below.

```yaml
apisix:
  # ... ...
  admin_key
    -
      name: "admin"
      key: abcdefghabcdefgh # Modify the original key to abcdefghabcdefgh
      role: admin
```

When we need to access the Admin API, we can use the key above, as shown below.

```shell
curl http://127.0.0.1:9080/apisix/admin/routes?api_key=abcdefghabcdefgh -i
```

The status code 200 in the returned result indicates that the access was successful, as shown below.

```shell
HTTP/1.1 200 OK
Date: Fri, 28 Feb 2020 07:48:04 GMT
Content-Type: text/plain
... ...
{"node":{...},"action":"get"}
```

At this point, if the key you enter does not match the value of `apisix.admin_key` in `conf/config.yaml`, for example, we know that the correct key is `abcdefghabcdefgh`, but we enter an incorrect key, such as `wrong- key`, as shown below.

```shell
curl http://127.0.0.1:9080/apisix/admin/routes?api_key=wrong-key -i
```

The status code `401` in the returned result indicates that the access failed because the `key` entered was incorrect and did not pass authentication, triggering an `Unauthorized` error, as shown below.

```shell
HTTP/1.1 401 Unauthorized
Date: Fri, 28 Feb 2020 08:17:58 GMT
Content-Type: text/html
... ...
{"node":{...},"action":"get"}
```

## Step 6: Build OpenResty for Apache APISIX

Some features require additional NGINX modules to be introduced into OpenResty. If you need these features, you can build OpenResty with [this script](https://raw.githubusercontent.com/api7/apisix-build-tools/master/build-apisix-openresty.sh).

## Step 7: Add Systemd Unit File for Apache APISIX

If you are using CentOS 7 and you installed Apache APISIX via the RPM package in step 2, the configuration file is already in place automatically and you can run the following command directly.

```shell
systemctl start apisix
systemctl stop apisix
```

If you installed Apache APISIX by other methods, you can refer to the [configuration file template](https://github.com/api7/apisix-build-tools/blob/master/usr/lib/systemd/system/apisix.service) for modification and put it in the `/usr/lib/systemd/system/apisix.service` path.