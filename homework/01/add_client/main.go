package main

import (
	"context"
	"flag"
	"log"
	"time"

	"code.ysh.homework.com/miniServer/add_clent/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var x = flag.Int64("x", 1, "输入x")
var y = flag.Int64("y", 1, "输入y")

func main() {
	flag.Parse()
	// 连接服务端
	conn, err := grpc.NewClient("127.0.0.1:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("conn failed: %v", err)
		return
	}
	defer conn.Close()
	c := proto.NewGretterClient(conn)

	// 调用rpc并打印
	ctx ,cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Add(ctx, &proto.AddResquest{X: *x, Y: *y})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	log.Printf("Greeting: %d", r.GetReply())
}