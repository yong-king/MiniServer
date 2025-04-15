package main

import (
	"client_bookstore/pb"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 拨号、连接
	conn , err := grpc.NewClient("127.0.0.1:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("conn failed:%v\n", err)
	}
	defer conn.Close()

	// 创建客户端
	c := pb.NewBookstoreClient(conn)
	res, err := c.ListBooks(context.Background(), &pb.ListBooksRequest{Shelf: 3})
	if err != nil {
		log.Fatal("c.ListBooks failed:%v\n", err)
	}
	fmt.Printf("next page token:%v\n", res.GetNextPageToken())
	for i, book := range(res.GetBooks()){
		fmt.Println("%d: %#v\n", i, book)
	}
}
