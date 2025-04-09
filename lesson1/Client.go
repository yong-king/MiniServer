package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Args struct {
	X, Y int
}

func main() {
	// 建立HTTP请求
	//client, err := rpc.DialHTTP("tcp", "127.0.0.1:9091")
	// 基于tcp协议的RPC
	//conn, err := rpc.Dial("tcp", "127.0.0.1:9091")
	// 基于json
	conn, err := net.Dial("tcp", "127.0.0.1:9091")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	// 同步调用
	args := &Args{10, 20}
	var reply int
	err = client.Call("ServiceA.Add", args, &reply)
	if err != nil {
		log.Fatal("Add:", err)
	}
	fmt.Println("Add:", reply)

	// 异步调用
	var reply2 int
	divCall := client.Go("ServiceA.Add", args, &reply2, nil)
	replyCall := <-divCall.Done // 接收调用结果
	fmt.Println("Add2:", replyCall)
}

// python 调用
//import socket
//import json
//
//request = {
//"id": 0,
//"params": [{"x":10, "y":20}],  # 参数要对应上Args结构体
//"method": "ServiceA.Add"
//}
//
//client = socket.create_connection(("127.0.0.1", 9091),5)
//client.sendall(json.dumps(request).encode())
//
//rsp = client.recv(1024)
//rsp = json.loads(rsp.decode())
//print(rsp)
