package service

import (
	"context"
	"fmt"
	"time"

	"com.ysh.kit/demo/pb"
)

// Sum 两数字相加
func (s addService) Sum(_ context.Context, a, b int) (int, error) {
	// 使用一个全局的日志记录
	// zap().L() ...
	return a + b, nil
}

// grpc
func (s *GrpcServer) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	_, rep, err := s.GRPCsum.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SumResponse), nil
}

// 日志中间件
func (mw logMiddleware) Sum(ctx context.Context, a, b int) (res int, err error) {
	defer func (start time.Time)  {
		mw.logger.Log(
			"method", "sum",
			"a", a,
			"b", b,
			"output", res,
			"err", err,
			"took", time.Since(start),
		)
	}(time.Now())
	res, err = mw.next.Sum(ctx, a, b)
	return
}

// 指标采集中间件
func (mw instrumentingMiddleware) Sum(ctx context.Context, a, b int) (res int, err error) {
	defer func (start time.Time)  {
		lvs := []string{"method", "sum", "error", fmt.Sprint(err != nil)}
		mw.requrestCount.With(lvs...).Add(1)
		mw.requestLatancy.With(lvs...).Observe(time.Since(start).Seconds())
	}(time.Now())

	res, err = mw.next.Sum(ctx, a, b)
	return
}

