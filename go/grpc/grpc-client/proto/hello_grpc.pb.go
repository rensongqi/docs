// 使用proto3语法

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.0
// source: hello.proto

package service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SayHello_SayHello_FullMethodName            = "/SayHello/SayHello"
	SayHello_PrintSum_FullMethodName            = "/SayHello/PrintSum"
	SayHello_SendFSMessageToUser_FullMethodName = "/SayHello/SendFSMessageToUser"
)

// SayHelloClient is the client API for SayHello service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 定义服务
// 然后我们需要定义一个服务，在这个服务中需要有一个方法，这个方法可以接收客户端的参数，再返回服务端的响应
// 其实很容易可以看出，我们定义了一个service，称为SayHello，这个服务中有一个rpc方法，名为SayHello
// 它将发送一个 HelloRequest 类型的请求消息，并期望服务器返回一个 HelloResponse 类型的响应消息。
type SayHelloClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
	PrintSum(ctx context.Context, in *PrintSumRequest, opts ...grpc.CallOption) (*PrintSumResponse, error)
	SendFSMessageToUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error)
}

type sayHelloClient struct {
	cc grpc.ClientConnInterface
}

func NewSayHelloClient(cc grpc.ClientConnInterface) SayHelloClient {
	return &sayHelloClient{cc}
}

func (c *sayHelloClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, SayHello_SayHello_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sayHelloClient) PrintSum(ctx context.Context, in *PrintSumRequest, opts ...grpc.CallOption) (*PrintSumResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PrintSumResponse)
	err := c.cc.Invoke(ctx, SayHello_PrintSum_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sayHelloClient) SendFSMessageToUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, SayHello_SendFSMessageToUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SayHelloServer is the server API for SayHello service.
// All implementations must embed UnimplementedSayHelloServer
// for forward compatibility.
//
// 定义服务
// 然后我们需要定义一个服务，在这个服务中需要有一个方法，这个方法可以接收客户端的参数，再返回服务端的响应
// 其实很容易可以看出，我们定义了一个service，称为SayHello，这个服务中有一个rpc方法，名为SayHello
// 它将发送一个 HelloRequest 类型的请求消息，并期望服务器返回一个 HelloResponse 类型的响应消息。
type SayHelloServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
	PrintSum(context.Context, *PrintSumRequest) (*PrintSumResponse, error)
	SendFSMessageToUser(context.Context, *UserRequest) (*UserResponse, error)
	mustEmbedUnimplementedSayHelloServer()
}

// UnimplementedSayHelloServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSayHelloServer struct{}

func (UnimplementedSayHelloServer) SayHello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedSayHelloServer) PrintSum(context.Context, *PrintSumRequest) (*PrintSumResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrintSum not implemented")
}
func (UnimplementedSayHelloServer) SendFSMessageToUser(context.Context, *UserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFSMessageToUser not implemented")
}
func (UnimplementedSayHelloServer) mustEmbedUnimplementedSayHelloServer() {}
func (UnimplementedSayHelloServer) testEmbeddedByValue()                  {}

// UnsafeSayHelloServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SayHelloServer will
// result in compilation errors.
type UnsafeSayHelloServer interface {
	mustEmbedUnimplementedSayHelloServer()
}

func RegisterSayHelloServer(s grpc.ServiceRegistrar, srv SayHelloServer) {
	// If the following call pancis, it indicates UnimplementedSayHelloServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SayHello_ServiceDesc, srv)
}

func _SayHello_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SayHelloServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SayHello_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SayHelloServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SayHello_PrintSum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrintSumRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SayHelloServer).PrintSum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SayHello_PrintSum_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SayHelloServer).PrintSum(ctx, req.(*PrintSumRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SayHello_SendFSMessageToUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SayHelloServer).SendFSMessageToUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SayHello_SendFSMessageToUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SayHelloServer).SendFSMessageToUser(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SayHello_ServiceDesc is the grpc.ServiceDesc for SayHello service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SayHello_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SayHello",
	HandlerType: (*SayHelloServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _SayHello_SayHello_Handler,
		},
		{
			MethodName: "PrintSum",
			Handler:    _SayHello_PrintSum_Handler,
		},
		{
			MethodName: "SendFSMessageToUser",
			Handler:    _SayHello_SendFSMessageToUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hello.proto",
}
