package main

import (
	"context"
	"log"
	"net"
	"testing"

	"com.ysh.kit/demo/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// 使用bufconn构建测试链接，避免使用实际端口号启动服务

const bufSize = 1024 * 1024

var bufListener *bufconn.Listener

func init() {
	bufListener = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	gs := NewGRPCServer(addService{})
	pb.RegisterAddServer(s, gs)
	go func() {
		if err := s.Serve(bufListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return bufListener.Dial()
}

func TestSum(t *testing.T) {
	// 1. 建立连接
	conn, err := grpc.NewClient(
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 2. 连接客户端
	c := pb.NewAddClient(conn)

	// 3.发起gprc请求
	resq, err := c.Sum(context.Background(), &pb.SumRequest{A: 10, B: 2})
	assert.Nil(t, err)
	assert.NotEmpty(t, resq)
	assert.Equal(t, resq.V, int64(12))
}

func TestConcat(t *testing.T) {
	// 1. 建立连接
	conn, err := grpc.NewClient(
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fail()
	}
	defer conn.Close()

	// 2. 连接客户端
	c := pb.NewAddClient(conn)

	// 3. 发起客户端请求
	resq, err := c.Concat(context.Background(), &pb.ConcatREquest{A: "10", B: "2"})
	assert.Nil(t, err)
	assert.NotEmpty(t, resq)
	assert.Equal(t, resq.V, "102")
}
