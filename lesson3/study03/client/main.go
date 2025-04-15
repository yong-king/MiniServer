package main

import (
	"context"
	"log"
	"time"

	"client.hello/demo/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// conn, err := grpc.NewClient("dns:///localhost:8972",
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("ysh:///resolver.ysh.com",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		log.Fatal("dial failed: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewGretterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for _ = range 10 {
		r, err := c.SayHello(ctx, &pb.SayHelloRequest{Name: "ysh"})
		if err != nil {
			log.Fatal("could not gretter: %v", err)
		}
		log.Print(r.GetResult())
	}
}
