package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"server.hello/demo/pb"
)

var port = flag.Int("port", 8972, "服务端口")

type server struct{
	pb.UnimplementedGretterServer
	Addr string
}

func (s *server) SayHello(ctx context.Context, in *pb.SayHelloRequest) (*pb.SayHelloResponse, error) {
	reply := fmt.Sprintf("hello %s .[form %s]", in.GetName(), s.Addr)
	return &pb.SayHelloResponse{Result: reply}, nil
}

func main() {
	flag.Parse()
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen failed:%v\n", err)
	}

	// 创建grpc服务
	s := grpc.NewServer()
	// 注册服务
	pb.RegisterGretterServer(s, &server{Addr: addr})

	// 启动服务
	fmt.Printf("start to server")
	err = s.Serve(l)
	if err != nil {
		fmt.Printf("failed to serve,err:%v\n", err)
		return
	}
}