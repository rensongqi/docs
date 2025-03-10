package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/chyroc/lark"
	"github.com/chyroc/lark/card"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	pb "grpc/grpc-server/proto"
)

type Server struct {
	pb.UnimplementedSayHelloServer
}

func FsContent() string {
	describe := "【小酷运维】"
	//构造button传递值
	config := lark.MessageContentCardConfig{
		UpdateMulti: true,
	}

	content := card.Card(
		card.Markdown("cdscds"),
	).
		SetHeader(card.Header(describe).SetRed()).SetConfig(&config)
	return content.String()
}

func (s *Server) SendFSMessageToUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	cli := lark.New(
		lark.WithAppCredential("cli_a228448b037ed00d", "xyRevQ1Z069gCnXA3RyBnftyFSVmj8hn"),
		//lark.WithAppCredential("cli_a36dfb87beb8900d", "UysnvQLnaYdHIEYt8wdhCcniQXxK4PZw"),
	)

	for _, userId := range req.UserIds {
		res, _, err := cli.Message.Send().ToUserID(userId).SendCard(ctx, FsContent())
		fmt.Println("res: ", res)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	return &pb.UserResponse{Res: "success"}, nil
}

func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("未传输 token")
	}

	fmt.Println("md: ", md)

	var appId string
	var appKey string
	if v, ok := md["appid"]; ok {
		appId = v[0]
	}
	if v, ok := md["appkey"]; ok {
		appKey = v[0]
	}

	if appId != "rensongqi" || appKey != "123456" {
		return nil, errors.New("token 不正确")
	}
	return &pb.HelloResponse{ResponseMsg: "nihao helloworld " + req.RequestName}, nil
}

func (s *Server) PrintSum(ctx context.Context, req *pb.PrintSumRequest) (*pb.PrintSumResponse, error) {
	sum := make([]int64, 0)
	for i := 0; i < int(req.Num); i++ {
		sum = append(sum, int64(i))
	}
	return &pb.PrintSumResponse{
		Nums: sum,
	}, nil
}

func main() {
	// TLS 认证
	//cred, err := credentials.NewServerTLSFromFile("/home/rensongqi/goprojects/gitlab/testing/grpc_service/key/server.pem", "/home/rensongqi/goprojects/gitlab/testing/grpc_service/key/server.key")
	//if err != nil {
	//	log.Fatal(err)
	//}
	// 创建grpc服务 认证
	//grpcServer := grpc.NewServer(grpc.Creds(cred))

	// 启动服务
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln(err)
	}

	// 创建grpc服务 无认证
	grpcServer := grpc.NewServer()

	// 在grpc中注册服务
	pb.RegisterSayHelloServer(grpcServer, &Server{})

	// 启动服务
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln(err)
	}
}
