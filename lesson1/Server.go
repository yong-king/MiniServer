package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type ArgsS struct {
	X, Y int
}

// ServiceA 定义一个结构体类型
type ServiceA struct{}

// Add为ServeceA类型新增一个可以导出的Add方法
func (s *ServiceA) Add(args *ArgsS, reply *int) error {
	*reply = args.X + args.Y
	return nil
}

func main() {
	server := new(ServiceA)
	rpc.Register(server) // 注册RPC服务
	//rpc.HandleHTTP()     // 基于HTTP协议
	// 基于tcp协议的RPC
	l, e := net.Listen("tcp", ":9091")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	//http.Serve(l, nil)
	for {
		conn, _ := l.Accept()
		//rpc.ServeConn(conn)
		// json协议
		rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
