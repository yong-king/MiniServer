package main

import (
	"context"
	"flag"
	"log"
	"time"

	"client/gateway/demo/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var name = flag.String("name", "ysh", "your name")

func main() {
	flag.Parse()
	// 连接server
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("client main. connect failed: %v", err)
	}
	defer conn.Close()
	c := proto.NewGretterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: *name})	
	if err != nil {
		log.Fatal("could not gretter: %v", err)
	}
	log.Print(r.GetMessage())
}