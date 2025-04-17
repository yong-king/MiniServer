package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"com.ysh/trim/demo/pb"
)


const serviceName = "trim_service"

var (
	port       = flag.Int("port", 8975, "service prot")
	consulAddr = flag.String("consul", "localhost:8500", "consul address")
)

type serve struct {
	pb.UnimplementedTrimServer
}

func (s *serve) TrimSpace(ctx context.Context, req *pb.TrimRequest) (*pb.TrimResponse, error) {
	ov := req.GetS()
	v := strings.ReplaceAll(ov, " ", "")
	fmt.Printf("ov:%#v v:%#v\n", ov, v)
	return &pb.TrimResponse{S: v}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed listen %d, err:%v\n", *port, err)
	}

	s := grpc.NewServer()
	pb.RegisterTrimServer(s, &serve{})

	healthpb.RegisterHealthServer(s, health.NewServer())// consul 发来健康检查的RPC请求，这个负责返回OK

	// 注册服务
	cc, err := NewConsul(*consulAddr)
	if err != nil {
		fmt.Printf("failed to NewConsulClient: %v", err)
		return
	}
	ipInfo, err := GetOutboundIP()
	if err != nil {
		fmt.Printf("getOutboundIP failed, err:%v\n", err)
		return
	}
	if err := cc.RegisterService(serviceName, ipInfo.String(), *port); err != nil {
		fmt.Printf("regToConsul failed, err:%v\n", err)
		return
	}

	go func ()  {
		err = s.Serve(lis)
		if err != nil {
			log.Fatal("server failed, err:%v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	cc.Deregister(fmt.Sprintf("%s-%s-%d", serviceName, ipInfo.String(), *port))

}

type consul struct {
	client *api.Client
}

func NewConsul (addr string) (*consul, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &consul{c}, err
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

func (c *consul) RegisterService(serviceName string, ip string, port int) error {
	// 健康检查
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ip, port), // 这里一定是外部可以访问的地址
		Timeout:                        "10s",  // 超时时间
		Interval:                       "10s",  // 运行检查的频率
		// 指定时间后自动注销不健康的服务节点
		// 最小超时时间为1分钟，收获不健康服务的进程每30秒运行一次，因此触发注销的时间可能略长于配置的超时时间。
		DeregisterCriticalServiceAfter: "1m",
	}
	srv := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ip, port),
		Name:    serviceName,
		Tags:    []string{"ysh", "hello"},
		Address: ip,
		Port:    port,
		Check: check,
	}
	return c.client.Agent().ServiceRegister(srv)
}

// Deregister 注销服务
func (c *consul) Deregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}
