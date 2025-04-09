package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"homework.ysh.com/miniServer/addServer/proto"
)

type server struct {
	proto.UnimplementedGretterServer
}

func (s *server) Add (cnx context.Context, req *proto.AddResquest) (*proto.AddResponse, error) {
	x := req.GetX()
	y := req.GetY()
	return &proto.AddResponse{Reply: x+y}, nil
}

func main() {
	// 监听端口
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("could not listen: %v", err)
		return
	}
	s := grpc.NewServer()
	proto.RegisterGretterServer(s, &server{})
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("server failed: %v", err)
		return
	}
}