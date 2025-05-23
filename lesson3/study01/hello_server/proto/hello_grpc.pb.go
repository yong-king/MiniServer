// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GretterClient is the client API for Gretter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GretterClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
	LotsOfReplines(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (Gretter_LotsOfReplinesClient, error)
	LotsOfGreetings(ctx context.Context, opts ...grpc.CallOption) (Gretter_LotsOfGreetingsClient, error)
	BidiHello(ctx context.Context, opts ...grpc.CallOption) (Gretter_BidiHelloClient, error)
}

type gretterClient struct {
	cc grpc.ClientConnInterface
}

func NewGretterClient(cc grpc.ClientConnInterface) GretterClient {
	return &gretterClient{cc}
}

func (c *gretterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, "/proto.Gretter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gretterClient) LotsOfReplines(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (Gretter_LotsOfReplinesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Gretter_ServiceDesc.Streams[0], "/proto.Gretter/LotsOfReplines", opts...)
	if err != nil {
		return nil, err
	}
	x := &gretterLotsOfReplinesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Gretter_LotsOfReplinesClient interface {
	Recv() (*HelloResponse, error)
	grpc.ClientStream
}

type gretterLotsOfReplinesClient struct {
	grpc.ClientStream
}

func (x *gretterLotsOfReplinesClient) Recv() (*HelloResponse, error) {
	m := new(HelloResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gretterClient) LotsOfGreetings(ctx context.Context, opts ...grpc.CallOption) (Gretter_LotsOfGreetingsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Gretter_ServiceDesc.Streams[1], "/proto.Gretter/LotsOfGreetings", opts...)
	if err != nil {
		return nil, err
	}
	x := &gretterLotsOfGreetingsClient{stream}
	return x, nil
}

type Gretter_LotsOfGreetingsClient interface {
	Send(*HelloRequest) error
	CloseAndRecv() (*HelloResponse, error)
	grpc.ClientStream
}

type gretterLotsOfGreetingsClient struct {
	grpc.ClientStream
}

func (x *gretterLotsOfGreetingsClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gretterLotsOfGreetingsClient) CloseAndRecv() (*HelloResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(HelloResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gretterClient) BidiHello(ctx context.Context, opts ...grpc.CallOption) (Gretter_BidiHelloClient, error) {
	stream, err := c.cc.NewStream(ctx, &Gretter_ServiceDesc.Streams[2], "/proto.Gretter/BidiHello", opts...)
	if err != nil {
		return nil, err
	}
	x := &gretterBidiHelloClient{stream}
	return x, nil
}

type Gretter_BidiHelloClient interface {
	Send(*HelloRequest) error
	Recv() (*HelloResponse, error)
	grpc.ClientStream
}

type gretterBidiHelloClient struct {
	grpc.ClientStream
}

func (x *gretterBidiHelloClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gretterBidiHelloClient) Recv() (*HelloResponse, error) {
	m := new(HelloResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GretterServer is the server API for Gretter service.
// All implementations must embed UnimplementedGretterServer
// for forward compatibility
type GretterServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
	LotsOfReplines(*HelloRequest, Gretter_LotsOfReplinesServer) error
	LotsOfGreetings(Gretter_LotsOfGreetingsServer) error
	BidiHello(Gretter_BidiHelloServer) error
	mustEmbedUnimplementedGretterServer()
}

// UnimplementedGretterServer must be embedded to have forward compatible implementations.
type UnimplementedGretterServer struct {
}

func (UnimplementedGretterServer) SayHello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedGretterServer) LotsOfReplines(*HelloRequest, Gretter_LotsOfReplinesServer) error {
	return status.Errorf(codes.Unimplemented, "method LotsOfReplines not implemented")
}
func (UnimplementedGretterServer) LotsOfGreetings(Gretter_LotsOfGreetingsServer) error {
	return status.Errorf(codes.Unimplemented, "method LotsOfGreetings not implemented")
}
func (UnimplementedGretterServer) BidiHello(Gretter_BidiHelloServer) error {
	return status.Errorf(codes.Unimplemented, "method BidiHello not implemented")
}
func (UnimplementedGretterServer) mustEmbedUnimplementedGretterServer() {}

// UnsafeGretterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GretterServer will
// result in compilation errors.
type UnsafeGretterServer interface {
	mustEmbedUnimplementedGretterServer()
}

func RegisterGretterServer(s grpc.ServiceRegistrar, srv GretterServer) {
	s.RegisterService(&Gretter_ServiceDesc, srv)
}

func _Gretter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GretterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Gretter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GretterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gretter_LotsOfReplines_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(HelloRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GretterServer).LotsOfReplines(m, &gretterLotsOfReplinesServer{stream})
}

type Gretter_LotsOfReplinesServer interface {
	Send(*HelloResponse) error
	grpc.ServerStream
}

type gretterLotsOfReplinesServer struct {
	grpc.ServerStream
}

func (x *gretterLotsOfReplinesServer) Send(m *HelloResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Gretter_LotsOfGreetings_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GretterServer).LotsOfGreetings(&gretterLotsOfGreetingsServer{stream})
}

type Gretter_LotsOfGreetingsServer interface {
	SendAndClose(*HelloResponse) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type gretterLotsOfGreetingsServer struct {
	grpc.ServerStream
}

func (x *gretterLotsOfGreetingsServer) SendAndClose(m *HelloResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gretterLotsOfGreetingsServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Gretter_BidiHello_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GretterServer).BidiHello(&gretterBidiHelloServer{stream})
}

type Gretter_BidiHelloServer interface {
	Send(*HelloResponse) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type gretterBidiHelloServer struct {
	grpc.ServerStream
}

func (x *gretterBidiHelloServer) Send(m *HelloResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gretterBidiHelloServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Gretter_ServiceDesc is the grpc.ServiceDesc for Gretter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gretter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Gretter",
	HandlerType: (*GretterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Gretter_SayHello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "LotsOfReplines",
			Handler:       _Gretter_LotsOfReplines_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "LotsOfGreetings",
			Handler:       _Gretter_LotsOfGreetings_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "BidiHello",
			Handler:       _Gretter_BidiHello_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "hello.proto",
}
