package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"com.ysh.kit/demo/pb"
	"com.ysh.kit/demo/service"
	"com.ysh.kit/demo/transport"
	"com.ysh.kit/demo/trim"
)

var (
	httpAddr = flag.Int("http-port", 8090, "HTTP端口")
	grpcAddr = flag.Int("grpc-port", 8972, "GRPC端口")
	trimAddr = flag.String("trim-addr", "127.0.0.1:8975", "trim端口")
	consulAddr = flag.String("consul", "localhost:8500", "consul address")
)

func main() {
	flag.Parse()
	// 前置资源初始化
	logger := log.NewJSONLogger(os.Stderr)

	svc := service.NewService()
	svc = service.NewlogMiddleware(logger, svc)
	svc = service.NewinstrumentingMiddleware(
		requestCount,
		requestLatency,
		countResult,
		svc,
	)

	// conn, err := grpc.NewClient(
	// 	*trimAddr,
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	
	// if err != nil {
	// 	fmt.Printf("connt %s failed, err:%v", *trimAddr, err)
	// 	return
	// }
	// defer conn.Close()

	// trimEndpoint := trim.MakeTrimEndpoint(conn)
	trimEndpoint, err := trim.GetTrimServiceFromConsul(*consulAddr, logger, "trim_service", nil)
	if err != nil {
		fmt.Printf("connect %s failed, err: %v", *trimAddr, err)
		return
	}
	svc = trim.NewwithTrimMiddleware(svc, trimEndpoint)
	var g errgroup.Group

	// HTTP
	g.Go(func() error {
		httpListen, err := net.Listen("tcp", fmt.Sprintf(":%d", *httpAddr))
		if err != nil {
			fmt.Printf("net listen %d failed: %v", *httpAddr, err)
			return err
		}
		defer httpListen.Close()

		logger := log.NewLogfmtLogger(os.Stderr)
		httpHandler := transport.NewHTTPServer(svc, logger)

		return http.Serve(httpListen, httpHandler)
	})

	// GRPC
	g.Go(func() error {
		grpcListen, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcAddr))
		if err != nil {
			fmt.Printf("net listen %d failed: %v", *grpcAddr, err)
			return err
		}
		defer grpcListen.Close()

		s := grpc.NewServer()
		pb.RegisterAddServer(s, transport.NewGRPCServer(svc))

		return s.Serve(grpcListen)
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("server exit with err:%v\n", err)
	}
}
