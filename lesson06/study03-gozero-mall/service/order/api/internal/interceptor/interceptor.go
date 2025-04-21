package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type CtxKey string
const(
	CtxKeyAdmindID CtxKey = "adminID"
)

func YshInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("客户端拦截器 in")
	// RPC调用前
	// 编写客户端拦截器的逻辑
	adminID := ctx.Value(CtxKeyAdmindID).(string)
	md := metadata.Pairs(
		"token", "ysh&dlrb",
		"adminID", adminID, 
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	err := invoker(ctx, method, req, reply, cc, opts...) // 实际的RPC调用

	// RPC调用后
	fmt.Println("客户端拦截器 out")
	return err
}