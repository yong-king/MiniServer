package main

import (
	"context"
	"hello_server/pb"
	"net"
    "fmt"
	"google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedGreeterServer
}

func (s *server) SayHello (ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
    reply := in.GetName()
    return &pb.HelloResponse{Reply: "Hello " + reply}, nil
}

func main() {
    // 监听本地8972端口
    l, e  := net.Listen("tcp", ":8972")
    if e != nil {
        fmt.Printf("failed to listeen %v", e)
        return
    }
    // 创建gprc服务器
    s := grpc.NewServer()
    // 在grpc服务端注册服务
    pb.RegisterGreeterServer(s, &server{})
    // 启动服务
    e = s.Serve(l)
    if e != nil {
        fmt.Printf("failed to server: %v", e)
        return 
    }
}