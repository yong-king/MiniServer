package main

import (
	"context"
	"flag"
	"log"
	"time"

	"code.ysh.com/miniserver/helo_client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name = flag.String("name", "ysh", "你的名字！")

func main() {
	flag.Parse()
	// 连接服务端，此处禁用安全传输
	conn, err := grpc.NewClient("127.0.0.1:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("did not connect:v", err)
		return
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	log.Printf("greeting: %s", r.GetReply())

}
