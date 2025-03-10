package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "grpc/grpc-client/proto"
)

var conn *grpc.ClientConn

type ClientTokenAuth struct{}

func (c ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appId":  "rensongqi",
		"appKey": "123456",
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires
// transport security.
func (c ClientTokenAuth) RequireTransportSecurity() bool {
	// 不开启安全认证，就返回false
	return false
}

// ConnGrpcServer 连接grpc server
func ConnGrpcServer() (*grpc.ClientConn, error) {
	// 自定义Token认证

	// TLS认证
	//cred, err := credentials.NewClientTLSFromFile("D:\\goprojects\\grpc_service\\key\\server.pem", "*.codewater.com")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	// 无认证
	//conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// 有认证
	//conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(cred))
	// 自定义token认证
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithPerRPCCredentials(new(ClientTokenAuth)))

	conn, err := grpc.Dial("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return conn, nil
}

func SayHello() error {
	// 建立连接
	client := pb.NewSayHelloClient(conn)
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{
		RequestName: "rensongqi",
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp.ResponseMsg)
	return nil
}

func PrintSum(num int64) error {
	// 建立连接
	client := pb.NewSayHelloClient(conn)
	resp, err := client.PrintSum(context.Background(), &pb.PrintSumRequest{
		Num: num,
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp.Nums)
	return nil
}

func SendFSMessage() error {
	userIds := []string{"89fcaab3", "89fcaab3"}
	client := pb.NewSayHelloClient(conn)
	resp, err := client.SendFSMessageToUser(context.Background(), &pb.UserRequest{
		UserIds: userIds,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.Res)
	return nil
}

func init() {
	var err error
	conn, err = ConnGrpcServer()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	_ = SayHello()
	_ = PrintSum(20)
	_ = SendFSMessage()
}
