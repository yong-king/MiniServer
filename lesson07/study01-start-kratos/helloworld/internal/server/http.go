package server

import (
	"context"
	"fmt"
	v1 "helloworld/api/bubble/v1"
	"helloworld/internal/conf"
	"helloworld/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func Middleware1() middleware.Middleware{
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			// 执行前
			fmt.Println("middle start")
			// 做token校验
			if tr, ok := transport.FromServerContext(ctx); ok{
				token := tr.RequestHeader().Get("token")
				fmt.Printf("toerkn:%v", token)
			}
			defer func() {
				fmt.Println("middle end")
			}()
			return handler(ctx, req)
		}
	}
}

// Middleware1 自定义中间件
func Middleware2(opts ...string) middleware.Middleware {
	return func(middleware.Handler) middleware.Handler {
		// opts
		return nil
	}
}

// Middleware2 自定义中间件3，相比Middleware1 失去了一些灵活性
func Middleware3(middleware.Handler) middleware.Handler {
	return nil
}


// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, todo *service.TodoService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(), // 全局中间件
			selector.Server(Middleware1(),
		).
			Path("/api.bubble.v1.Todo/CreateTodo").Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	// 替换默认到http响应编码器
	opts = append(opts, http.ResponseEncoder(responseEncoder))

	// 替换默认到错误编码响应
	opts = append(opts, http.ErrorEncoder(errorEncoder))

	srv := http.NewServer(opts...)
	v1.RegisterTodoHTTPServer(srv, todo)
	return srv
}
