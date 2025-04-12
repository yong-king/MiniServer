package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"hello_client/demo/proto"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	// "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var name = flag.String("name", "yk", "please input name")

// 一元拦截器
func unaryInterceptor(ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption) error {
	var credsConfigured bool
	// 遍历 opts，检查是否已经设置了认证凭据（PerRPCCreds）
	for _, o := range opts {
		_, ok := o.(*grpc.PerRPCCredsCallOption) // 类型断言
		if ok {
			credsConfigured = true
			break
		}
	}
	// 如果调用时没有显式设置 token，则自动附加一个默认的 AccessToken
	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(
			oauth.TokenSource{
				TokenSource: oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: "app-test-ysh"},
				)}))
	}
	// 记录请求开始时间
	statr := time.Now()
	// 执行真正的 RPC 请求
	err := invoker(ctx, method, req, reply, cc, opts...)
	// 记录请求结束时间
	end := time.Now()
	// 打印请求日志（方法名、时间、错误）
	fmt.Printf("RPC: %s, satrt time: %s, end time: %s, err: %v\n",
		method, statr.Format("Basic"), end.Format(time.RFC3339), err)
	// 返回 RPC 请求的结果（错误或 nil）
	return err
}

// 流式拦截器
type wrappedStream struct {
	grpc.ClientStream // 嵌入原始的 grpc.ClientStream，实现包装
}

// 重写 RecvMsg 方法，当接收到消息时打印日志，然后调用原始方法
func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

// 重写 SendMsg 方法，当发送消息时打印日志，然后调用原始方法
func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("Send a message(Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

// newWrappedStream 是一个辅助函数，用于创建一个 wrappedStream 实例
func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

// streamInterceptor 是一个 gRPC 客户端流式拦截器（类似中间件），用于在发送/接收消息前后添加自定义逻辑。
func streamInterceptor(ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (
	grpc.ClientStream, error) {
	var credsConfigured bool
	for _, o := range opts {
		_, ok := o.(*grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}

	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(oauth.TokenSource{
			TokenSource: oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: "app-test-ysh"},
			),
		}))
	}
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return newWrappedStream(s), nil
}

func main() {
	flag.Parse()
	// 连接服务器
	creds, err := credentials.NewClientTLSFromFile("crets/server.crt", "yuanshuhao.com")
	if err != nil {
		log.Fatalf("credentials NewClientTLSFromFile failed: %v\n", err)
		return
	}
	conn, err := grpc.NewClient("127.0.0.1:8972",
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(unaryInterceptor),
		grpc.WithStreamInterceptor(streamInterceptor))
	// conn, err := grpc.NewClient("127.0.0.1:8972",
	// 	grpc.WithTransportCredentials(insecure.NewCredentials())) // 不安全连接
	if err != nil {
		log.Fatalf("grpc.NewClient falied: %v\n", err)
		return
	}
	defer conn.Close()
	c := proto.NewGretterClient(conn)

	cnx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 创建metadata
	// 基于metadata创建cnx
	cnx = metadata.AppendToOutgoingContext(cnx, "token", "app-test-ysh")

	var header, trailer metadata.MD
	r, err := c.SayHello(cnx,
		&proto.HelloRequest{Name: *name},
		grpc.Header(&header),
		grpc.Trailer(&trailer))
	if err != nil {
		s := status.Convert(err)
		for _, d := range s.Details() {
			switch info := d.(type) {
			case *errdetails.QuotaFailure:
				fmt.Printf("Quota failuer: %s\n", info.String())
			default:
				fmt.Printf("Unexpected type: %v\n", info)
			}
		}
		log.Fatalf("c.SayHello falied: %v", err)
		return
	}
	fmt.Printf("header:%v\n", header)
	log.Printf("gretiing: %s\n", r.GetReply())
	fmt.Printf("trailer:%v\n", trailer)
	// runLotsOfReplies(c)
	// runLotsOfGrettings(c)
	// runBihiHello(c)

}

func runLotsOfReplies(c proto.GretterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	streamm, err := c.LotsOfReplines(ctx, &proto.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("c.LotsOfReplines falied: %v\n", err)
		return
	}
	for {
		res, err := streamm.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("streamm.Recv falied: %v\n", err)
		}
		log.Printf("got reply: %q\n", res.GetReply())
	}

}

func runLotsOfGrettings(c proto.GretterClient) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.LotsOfGreetings(ctx)
	if err != nil {
		log.Fatalf("c.LotsOfGreetings falied: %v\n", err)
		return
	}
	names := []string{"ysh", "zrn", "dlrb", "sy"}
	for _, name := range names {
		err := stream.Send(&proto.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("stream.Send falied: %v\n", err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("stream.CloseAndRecv faile: %v", err)
		return
	}
	log.Printf("got reply: %v", res)
}

func runBihiHello(c proto.GretterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "token", "app-test-ysh")
	// 使用 GreeterClient 的 BidiHello 方法启动一个双向流
	stream, err := c.BidiHello(ctx)
	if err != nil {
		log.Fatalf("c.BidiHello failed: %v\n", err)
		return
	}

	// 接受
	// 创建一个等待通道，用于等待接收消息的 goroutine 完成
	waitc := make(chan struct{})
	go func() {
		header, err := stream.Header()
		if err != nil {
			log.Fatalf("falied to get header form stream: %v\n", err)
		}
		fmt.Printf("header:%v\n", header)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				// 读完了
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("stream.Recv failed %v\n", err)
			}
			log.Printf("AI: %s\n", res.GetReply())
		}
	}()

	// 发送
	// 创建一个 reader 用于从标准输入中读取用户输入
	reader := bufio.NewReader(os.Stdin) // 从标准输入创建一个读取对象
	for {
		// 读取用户输入直到遇到换行符
		cmd, _ := reader.ReadString('\n')
		// 去掉输入字符串两端的空格和换行符
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		if strings.ToUpper(cmd) == "QUIT" {
			break
		}
		if err := stream.Send(&proto.HelloRequest{Name: cmd}); err != nil {
			log.Fatalf("client bihihello stream.Send failed: %v\n", err)
		}
	}
	// 关闭流的发送端，表示没有更多的数据要发送
	stream.CloseSend()
	<-waitc
	trailer := stream.Trailer()
	fmt.Printf("tralier: %v", trailer)
}
