package main

import (
	"context"
	"flag"
	"fmt"

	// "go/build/constraint"

	"user/rpc/internal/config"
	"user/rpc/internal/server"
	"user/rpc/internal/svc"
	"user/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 注册服务端拦截器
	s.AddUnaryInterceptors(yshInterceptor)

	// 注册consul
	_ = consul.RegisterService(c.ListenOn, c.Consul)
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}


func yshInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// 调用前
	fmt.Println("服务端拦截器 in")
	// 拦截器业务逻辑
	// 获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "need metadata")
	}
	fmt.Println("metadata:%#v\n", md)

	// // 根据metadata中的数据进行一些校验处理
	// if md["token"][0] != "ysh&dlrb" {
	// 	return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	// }

	m, err := handler(ctx, req) // 实际RPC方法

	// 调用后
	fmt.Println("服务端拦截器 out")
	return m, err
}