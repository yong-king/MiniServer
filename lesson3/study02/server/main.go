package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"server/gateway/demo/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.29.0
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

// protoc \
//   --proto_path=proto \
//   --go_out=proto --go_opt=paths=source_relative \
//   --go-grpc_out=proto --go-grpc_opt=paths=source_relative \
//   hello.proto

type server struct{
	proto.UnimplementedGretterServer
}

func (s *server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{Message: in.Name + " Hello"}, nil
}

func main(){
	// 建立tcp连接
	lis, err := net.Listen("tcp", ":8080")
	if err != nil{
		log.Fatal("main. listen failed: %v", err)
	}

	// 创建grpc对象
	s := grpc.NewServer()
	// 注册
	proto.RegisterGretterServer(s, &server{})

	// 启动grpc Server
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// 创建一个连接到刚刚启动的grpc服务器的客户端连接
	// gRPC-Gateway 就是通过它来代理请求（将HTTP请求转为RPC请求）
	conn, err := grpc.NewClient(
		"127.0.0.1:8080",
		// grpc.WithBlock()：确保连接是 同步 的，直到连接建立成功。
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil{
		log.Fatalln("falied to dial server:", err)
	}

	// 创建 gRPC-Gateway Mux（路由器）：它会将 HTTP 请求映射到 gRPC 服务。
	gwmux := runtime.NewServeMux()

	// 注册 gRPC 服务与 HTTP 路由的映射
	// proto.RegisterGretterHandler 是自动生成的代码，用于将 gRPC 服务注册到 gRPC-Gateway 的HTTP 路由中。
	// 这一行的作用是将 gRPC 服务的接口（conn）与 HTTP 请求的路由 (gwmux) 关联起来。
	err = proto.RegisterGretterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("falied to register gateway:", err)
	}

	// 这一行代码创建了一个 HTTP 服务器，绑定在 8090 端口，
	// 并将之前创建的 gwmux（HTTP 路由器）作为处理请求的 handler。
	gwServer := &http.Server{Addr: ":8090", Handler: gwmux}

	log.Fatalln(gwServer.ListenAndServe())

}