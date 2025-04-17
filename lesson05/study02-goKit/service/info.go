package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/metrics"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"

	"com.ysh.kit/demo/pb"
)

// 1.业务逻辑层

// 1.1 定义接口
type AddService interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(cte context.Context, a, b string) (string, error)
}

// 1.2实现接口
type addService struct {
	// db db.Conn
	// log log.Logger // 嵌入一个来用
}

var (
	// ErrTwoEmptyStrings 两个字符串都为空！
	ErrTwoEmptyStrings = errors.New("两个字符串都为空！")
)

// grpc
// Handler 应该从服务实现的gRPC绑定调用。
// 传入的请求参数和返回的响应参数都是gRPC类型，而不是用户域类型。
//
//	type Handler interface {
//		ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error)
//	}
type GrpcServer struct {
	pb.UnimplementedAddServer

	GRPCsum    grpctransport.Handler
	GRPCconcat grpctransport.Handler
}

// NewService addService的构造函数
func NewService() AddService {
	return &addService{
		// db:db
	}
}

// 日志中间件，方案三
type logMiddleware struct {
	logger log.Logger
	next   AddService
}

func NewlogMiddleware(logger log.Logger, svc AddService) AddService {
	return &logMiddleware{
		logger: logger,
		next:   svc,
	}
}

// metrics 指标采集
type instrumentingMiddleware struct {
	requrestCount  metrics.Counter
	requestLatancy metrics.Histogram
	countResult    metrics.Histogram
	next           AddService
}

func NewinstrumentingMiddleware(count metrics.Counter, latancy, result metrics.Histogram, srv AddService) AddService{
	return &instrumentingMiddleware{
		requrestCount: count,
		requestLatancy: latancy,
		countResult: result,
		next: srv,
	}
}

