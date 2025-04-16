package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"com.hello.client/pb"
	// "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

var name = flag.String("name", "ysh", "通过-name告诉server你是谁")

func main() {
	flag.Parse()
	/*
		// 1. 连接到consul
		cc, err := api.NewClient(api.DefaultConfig())
		if err != nil{
			fmt.Printf("conn consul failed:%v\n", err)
			return
		}

		// 2. 根据服务名称查询实例
		// cc.Agent().Services()  // 列出所有的
		serviceMap, err := cc.Agent().ServicesWithFilter("Service=`hello`")// 查询服务名称是hello的所有服务节点
		if err != nil {
			fmt.Printf("query `hello` service failed:%v\n",err)
			return
		}
		var addr string
		for k, v := range serviceMap {
			fmt.Printf("%s:%#v\n", k, v)
			addr = fmt.Sprintf("%s:%d", v.Address, v.Port)
			break
		}
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		// conn, err := grpc.NewClient("127.0.0.1:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
	*/
	conn, err := grpc.NewClient("consul://localhost:8500/hello", // grpc中使用consul名称解析器，
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("failed to connect server:%v\n", err)
	}
	defer conn.Close()

	c := pb.NewGretterClient(conn)

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := c.SayHello(ctx, &pb.SayHelloRequerst{Name: *name})

		if err != nil {
			fmt.Printf("c.SayHello failed, err:%v\n", err)
			return
		}
		// 拿到了RPC响应
		fmt.Printf("resp:%v\n", resp.GetReply())
	}

}
