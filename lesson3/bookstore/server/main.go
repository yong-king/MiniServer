package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"com.ysh.blog.booksotre/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接数据库
	db, err := NewDB()
	if err != nil {
		log.Fatalf("connect to db failed: %v\n", err)
	}

	// 创建server
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatalf("failed to listen to 8090:%v\n", err)
	}

	srv := server{
		bs: &bookstore{db},
	}
	s := grpc.NewServer()
	// 注册
	pb.RegisterBookstoreServer(s, &srv)
	go func() {
		fmt.Println("server start!")
		err := s.Serve(l)
		if err != nil {
			log.Fatalf("filed to server: %v\n", err)
		}
	}()

	// grpc-ateway
	conn, err :=grpc.NewClient(
		"127.0.0.1:8090",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),)
	if err != nil{
		fmt.Printf("grpc conn failed:%v\n", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterBookstoreHandler(context.Background(), gwmux, conn)
	if err != nil{
		log.Fatalf("failed to register handelr:%v\n", err)
	}

	gwServer := &http.Server{
		Addr: ":8080",
		Handler: gwmux,
	}
	fmt.Println("start to server 8080...")
	err = gwServer.ListenAndServe()
	if err != nil {
		log.Fatal("failed to start http server:%v\n", err)
	}
}
