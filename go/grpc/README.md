# 1 gRPC

- [官方手册](https://grpc.io/docs/languages/go/quickstart/)

## 1.1 安装protoc及protoc-gen-go

```bash
//1、安装转换文件protoc
https://github.com/protocolbuffers/protobuf/releases

//2、安装protoc-gen-go和protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 导入系统环境变量
export PATH="$PATH:$(go env GOPATH)/bin"
```

下边创建一个简单的server端和client端，实现两者之间的交互

## 1.2 grpc server

**目录结构：**

```bash
├── go.mod
├── go.sum
├── pbfiles
│   └── Prod.proto
├── server.go
└── services
    ├── Prod.pb.go
    └── ProdService.go

2 directories, 6 files
```

1. 创建proto配置文件`Prod.proto`

```go
// 指定的当前proto语法的版本，有2和3
syntax = "proto3";
// 指定等会文件生成出来的package
option go_package = ".;../services";
// 定义request model
message ProductRequest{
  int32 prod_id = 1; // 1代表顺序
}
// 定义response model
message ProductResponse{
  int32 prod_stock = 1; // 1代表顺序
}
// 定义服务主体
service ProdService{
  // 定义方法
  rpc GetProductStock(ProductRequest) returns(ProductResponse);
}
```

2. 生成protobuf中间文件（Prod.pb.go）

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helloworld/Prod.proto
```

3. services层编写商品服务（ProdService.go）

```go
package services

import (
	"context"
)

type ProdService struct {

}

func (this *ProdService) GetProductStock(ctx context.Context, request *ProductRequest) (*ProductResponse, error) {
	return &ProductResponse{ProdStock: 20}, nil
}
```

4. 编写server服务（server.go）

```go
package main

import (
	"google.golang.org/grpc"
	"grpc_server/services"
	"net"
)

func main() {
	rpcServer := grpc.NewServer()
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))

	lis, _ := net.Listen("tcp", ":8081")
	rpcServer.Serve(lis)
}
```

## 1.3 grpc_client

目录结构：

```bash
├── go.mod
├── go.sum
├── client.go
└── services
    └── Prod.pb.go

1 directory, 4 files
```

1. 首先先把server端的Prod.pb.go文件拷贝到client端的services目录下
2. client端连接grpc代码

```go
package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc_client/services"
	"log"
)

func main() {
	// 不使用https的认证 grpc.WithInsecure() 连接grpc server端
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// 获取service的client
	prodClient := services.NewProdServiceClient(conn)
	prodResult, err := prodClient.GetProductStock(context.Background(), &services.ProductRequest{ProdId: 15})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("prod result ===> ", prodResult.ProdStock)
}
```

## 1.4 连通性测试

server端先起来起来，client后运行，查看prodResult结果

```bash
# (client) go run main.go
prod result ===>  prod_stock:20
```

# 2 gRPC http证书配置

server端提供http服务：

```go
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		rpcServer.ServeHTTP(writer, request)
	})
	httpServer := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	httpServer.ListenAndServeTLS("cert/server.pem", "cert/server.key")
```

## 2.1 简单加入证书配置

### 2.2.1 Server端

```go
	creds, err := credentials.NewServerTLSFromFile("keys/server.crt", "keys/server_no_passwd.key")
	if err != nil {
		log.Fatal(err)
	}

	rpcServer := grpc.NewServer(grpc.Creds(creds))
```

### 2.2.2 Client端

```go
	creds, err := credentials.NewClientTLSFromFile("keys/server.crt", "localhost")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(creds))
```

## 2.2 自签证书双向认证

### 2.2.1 生成CA证书

```bash
# openssl genrsa -out ca.key 2048
# openssl req -new -x509 -days 3650 -key ca.key -out ca.pem
Country Name (2 letter code) []:cn
State or Province Name (full name) []:shanghai
Locality Name (eg, city) []:shanghai
Organization Name (eg, company) []:rsq  
Organizational Unit Name (eg, section) []:rsq
Common Name (eg, fully qualified host name) []:localhost
Email Address []:
```

### 2.2.2 生成server端证书

```bash
# openssl genrsa -out server.key 2048
# openssl req -new -key server.key -out server.csr
Country Name (2 letter code) []:cn
State or Province Name (full name) []:shanghai
Locality Name (eg, city) []:shanghai
Organization Name (eg, company) []:rsq
Organizational Unit Name (eg, section) []:rsq
Common Name (eg, fully qualified host name) []:localhost 
Email Address []:

Please enter the following 'extra' attributes
to be sent with your certificate request
A challenge password []:
# openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server.pem

```

### 2.2.3 生成client端证书

```bash
# openssl ecparam -genkey -name secp384r1 -out client.key
# openssl req -new -key client.key -out client.csr

# openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem
```

### 2.2.4 server端

拷贝ca.pem、server.key和server.pem至服务端的cert目录下（新建cert目录）

```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc_server/services"
	"io/ioutil"
	"net"
)

func main() {
	cert, _ := tls.LoadX509KeyPair("cert/server.pem", "cert/server.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, //服务端证书
		ClientAuth: tls.RequireAndVerifyClientCert, //双向认证
		ClientCAs: certPool,
	})

	rpcServer := grpc.NewServer(grpc.Creds(creds))
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))

	lis, _ := net.Listen("tcp", ":8081")
	rpcServer.Serve(lis)
}
```

### 2.2.5 client端

拷贝ca.pem、client.key和client.pem至服务端的cert目录下（新建cert目录）
client端运行前需要配置环境变量（GODEBUG=x509ignoreCN=0）运行命令如下：

```bash
# GODEBUG=x509ignoreCN=0 go run main.go 
```

```go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc_client/services"
	"io/ioutil"
	"log"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cert/client.pem", "cert/client.key")
	if err != nil {
		log.Println("11111")
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, //客户端证书
		ServerName: "localhost", //双向认证
		RootCAs: certPool,
	})

	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 获取service的client
	prodClient := services.NewProdServiceClient(conn)
	prodResult, err := prodClient.GetProductStock(context.Background(), &services.ProductRequest{ProdId: 15})
	if err != nil {
		log.Println("22222")
		log.Fatal(err)
	}

	fmt.Println("prod result ===> ", prodResult.ProdStock)
}
```

# 3 gRPC gateway

```bash
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway 
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

```










