package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"com.ysh.blog.booksotre/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
	// go func() {
	// 	fmt.Println("server start!")
	// 	err := s.Serve(l)
	// 	if err != nil {
	// 		log.Fatalf("filed to server: %v\n", err)
	// 	}
	// }()

	// // grpc-ateway
	// conn, err :=grpc.NewClient(
	// 	"127.0.0.1:8090",
	// 	grpc.WithBlock(),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),)
	// if err != nil{
	// 	fmt.Printf("grpc conn failed:%v\n", err)
	// }

	// gwmux := runtime.NewServeMux()
	// err = pb.RegisterBookstoreHandler(context.Background(), gwmux, conn)
	// if err != nil{
	// 	log.Fatalf("failed to register handelr:%v\n", err)
	// }

	// gwServer := &http.Server{
	// 	Addr: ":8080",
	// 	Handler: gwmux,
	// }
	// fmt.Println("start to server 8080...")
	// err = gwServer.ListenAndServe()
	// if err != nil {
	// 	log.Fatal("failed to start http server:%v\n", err)
	// }

	// 基于同一个端口提供HTTP API和gRPC API
	// 创建grpc- Gateway mux
	gwmux := runtime.NewServeMux()
	drops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// 注册
	err = pb.RegisterBookstoreHandlerFromEndpoint(context.Background(), gwmux, "127.0.0.1:8090", drops)
	if err != nil {
		log.Fatalf("falied to register gwmux:%v\n", err)
	}

	// 创建http mux
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	// 定义http server
	gwServer := &http.Server{
		Addr: "127.0.0.1:8090",
		Handler:grpcHandlerFunc(s, mux), // 请求的统一入口
	}

	// 启动服务
	fmt.Println("server on http 127.0.0.1:8090")
	log.Fatalln(gwServer.Serve(l)) // 启动HTTP服务

}


func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc"){
			// 如果是 HTTP/2 且请求的 Content-Type 包含 application/grpc，处理为 gRPC 请求
			grpcServer.ServeHTTP(w, r)
		}else{
			// 否则，交给另一个 HTTP 处理器处理
			otherHandler.ServeHTTP(w,r)
		}
	}), &http2.Server{})
}