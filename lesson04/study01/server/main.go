package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	// "net/http"
	// "strings"

	"com.hello.server/pb"
	// "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	// "golang.org/x/net/http2"
	// "golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"

	"github.com/hashicorp/consul/api"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const(
	serviceName = "hello"
)

type service struct{
	pb.UnimplementedGretterServer
	addr string
}


func (s *service)  SayHello(ctx context.Context, in *pb.SayHelloRequerst) (*pb.SayHelloResponse, error) {
	reply := s.addr + " Hello " + in.GetName()
	return &pb.SayHelloResponse{Reply: reply}, nil
}

var port int = 8972

func main(){
	flag.IntVar(&port, "port", port, "input port")
	flag.Parse()
	fmt.Printf("-->port: %d\n", port)
	// 启动服务
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil{
		log.Fatalf("failed to listen %d:%v\n", port, err)
	}

	// 创建grpc服务
	s := grpc.NewServer()
	// 注册服务
	ipinfo, err := GetOutboundIP()
	if err != nil {
		log.Fatalf("failed to get local ip:%v\n",err)
	}
	fmt.Println(ipinfo.String())

	addr := fmt.Sprintf("%s:%d", ipinfo.String(), port)
	pb.RegisterGretterServer(s, &service{addr: addr})
	// 给我们的gprc的服务增减增加注册健康检查
	healthpb.RegisterHealthServer(s, health.NewServer())// consul 发来健康检查的RPC请求，这个负责返回OK

	// 连接到consul
	cc, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("failed to conn consul")
	}
	// 将我们的gprc服务注册到consul
	// 1.定义服务
	// 配置健康检查
	// 获取本机的出口ip

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ipinfo.String(), port), // 外网地址
		Timeout:                        "5s",
		Interval:                       "5s",  // 间隔
		DeregisterCriticalServiceAfter: "10m", // 10分钟后注销掉不健康的服务节点
	}

	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, ipinfo.String(), port)
	srv := &api.AgentServiceRegistration{
		ID:      serviceID, // 服务唯一ID
		Name:    serviceName,
		Tags:    []string{"ysh"},
		Address: ipinfo.String(),
		Port:    port,
		Check:   check,
	}

	// 2. 注册到consun
	err = cc.Agent().ServiceRegister(srv)
	if err != nil {
		fmt.Printf("ServiceRegister failed, err:%v\n", err)
	}

	go func ()  {
		fmt.Printf("server on http 127.0.0.1:%d", port)
		err = s.Serve(lis)
		if err != nil{
			log.Fatalf("failed to server:%v", err)
		}
	}()

	// ctrl + c 退出程序
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)
	fmt.Println("wait quit signal...")
	<-quitCh // 没收到信号就阻塞

	// 程序退出，注销服务
	fmt.Println("service quit...")
	err = cc.Agent().ServiceDeregister(serviceID) // 注销服务
	if err != nil{
		fmt.Printf("service deregister failed: %v\n", err)
	}

	

	// gwmux := runtime.NewServeMux()
	// drops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// err = pb.RegisterGretterHandlerFromEndpoint(context.Background(), gwmux, "127.0.0.1:8972", drops)
	// if err != nil {
	// 	log.Fatalf("falied to register gwmux:%v\n", err)
	// }

	// mux := http.NewServeMux()
	// mux.Handle("/", gwmux)

	// gwServer := &http.Server{
	// 	Addr: "127.0.0.1:8972",
	// 	Handler: grpcHandlerFunc(s, mux),
	// }

	// // 启动服务
	// fmt.Println("server on http 127.0.0.1:8972")
	// log.Fatalln(gwServer.Serve(lis)) // 启动HTTP服务

}

// GetOutboundIP 获取本机的出口IP
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
// 	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc"){
// 			// 如果是 HTTP/2 且请求的 Content-Type 包含 application/grpc，处理为 gRPC 请求
// 			grpcServer.ServeHTTP(w, r)
// 		}else{
// 			// 否则，交给另一个 HTTP 处理器处理
// 			otherHandler.ServeHTTP(w,r)
// 		}
// 	}), &http2.Server{})
// }


