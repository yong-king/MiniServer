package main

import (
	"context"
	"fmt"
	"hello_server/demo/proto"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedGretterServer
	mu    sync.Mutex     // count 并发锁
	count map[string]int // map并发不安全
}

func (s *server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloResponse, error) {
	// 从上下问中读取
	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		grpc.SetTrailer(ctx, trailer)
	}()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "UnarySayHello: failed to get metadata")
	}
	token := md.Get("token")
	if len(token) < 1 || token[0] != "app-test-ysh" {
		return nil, status.Error(codes.Unauthenticated, "错误token")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 每个name只能调用一次
	s.count[in.Name]++        // 记录被访问的次数
	if s.count[in.Name] > 1 { // 访问次数超过一次
		st := status.New(codes.ResourceExhausted, "Request limit exceede.")
		dt, err := st.WithDetails(
			&errdetails.QuotaFailure{
				Violations: []*errdetails.QuotaFailure_Violation{{
					Subject:     fmt.Sprintf("names:%s", in.Name),
					Description: "限制每个name调用一次",
				}},
			},
		)
		if err != nil {
			return nil, st.Err()
		}
		return nil, dt.Err()
	}

	reply := in.GetName()
	header := metadata.New(map[string]string{"location": "Guangzhou"})
	grpc.SendHeader(ctx, header)
	return &proto.HelloResponse{Reply: "hello " + reply}, nil
}

// 验证函数
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// 执行token认证逻辑
	return token == "app-test-ysh"
}

// 一元拦截器
// 服务端定义一个一元拦截器，对从请求元数据中获取的authorization进行校验。
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	m, err := handler(ctx, req)
	if err != nil {
		fmt.Printf("RPC failed with error %v\n", err)
	}
	return m, err
}

// 流式拦截器
type wrappedStream struct {
	grpc.ServerStream // 嵌入原始的 grpc.ClientStream，实现包装
}

// 重写 RecvMsg 方法，当接收到消息时打印日志，然后调用原始方法
func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

// 重写 SendMsg 方法，当发送消息时打印日志，然后调用原始方法
func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("Send a message(Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

// newWrappedStream 是一个辅助函数，用于创建一个 wrappedStream 实例
func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// streamInterceptor 服务端流拦截器
func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// authentication (token verification)
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	if !valid(md["authorization"]) {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		fmt.Printf("RPC failed with error %v\n", err)
	}
	return err
}

func main() {
	// 监听
	l, err := net.Listen("tcp", ":8972")
	if err != nil {
		log.Fatalf("net.Listen failed: %v\n", err)
		return
	}
	creds, err := credentials.NewServerTLSFromFile("crets/server.crt", "crets/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile failed: %v\n", err)
		return
	}
	// 创建grpc服务器
	s := grpc.NewServer(grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor))
	// 注册服务
	proto.RegisterGretterServer(s, &server{count: map[string]int{}}) //  assignment to entry in nil map
	// 启动服务
	err = s.Serve(l)
	if err != nil {
		log.Fatalf("s.Serve falied: %\n", err)
		return
	}
}

func (s *server) LotsOfReplines(in *proto.HelloRequest, stream proto.Gretter_LotsOfReplinesServer) error {
	words := []string{
		"你好",
		"hello",
		"こんにちは",
		"안녕하세요",
	}

	for _, word := range words {
		data := &proto.HelloResponse{
			Reply: word + in.GetName(),
		}
		// 使用Send方法返回多个数据
		if err := stream.Send(data); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) LotsOfGreetings(stream proto.Gretter_LotsOfGreetingsServer) error {
	reply := "您好，"
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.HelloResponse{
				Reply: reply,
			})
		}
		if err != nil {
			return err
		}
		reply += res.GetName()
	}
}

func (s *server) BidiHello(stream proto.Gretter_BidiHelloServer) error {

	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		stream.SetTrailer(trailer)
	}()

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "BidirectionalStreamingSayHello: failed to get metadata")
	}
	token := md.Get("token")
	if len(token) < 1 || token[0] != "app-test-ysh" {
		return status.Error(codes.Unauthenticated, "valid token!")
	}

	header := metadata.New(map[string]string{"location": "js"})
	stream.SendHeader(header)
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		reply := magic(res.GetName()) // 对收到数数据进行处理

		// 返回流式响应
		if err := stream.Send(&proto.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
}

// magic 一段价值连城的“人工智能”代码
func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}
