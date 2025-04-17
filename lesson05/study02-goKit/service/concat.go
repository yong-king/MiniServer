package service

import (
	"context"
	"fmt"
	"time"

	"com.ysh.kit/demo/pb"
)

// Concat 字符串拼接
func (s addService) Concat(_ context.Context, a, b string) (string, error) {
	if a == "" && b == "" {
		return "", ErrTwoEmptyStrings
	}
	return a + b, nil
}



// grpc
func  (s *GrpcServer) Concat(ctx context.Context, req *pb.ConcatREquest) (*pb.ConcatResponse, error) {
	_, rep, err := s.GRPCconcat.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ConcatResponse), nil
}


func (mw logMiddleware) Concat(ctx context.Context, a, b string) (res string, err error) {
	defer func (start time.Time)  {
		mw.logger.Log(
			"method", "concat",
			"a", a,
			"b", b,
			"output", res,
			"err", err,
			"took", time.Since(start),
		)
	}(time.Now())
	res, err = mw.next.Concat(ctx, a, b)
	return
}

// 指标采集中间件
func (mw instrumentingMiddleware) Concat(ctx context.Context, a, b string) (res string, err error) {
	defer func (start time.Time)  {
		lvs := []string{"method", "sum", "error", fmt.Sprint(err != nil)}
		mw.requrestCount.With(lvs...).Add(1)
		mw.requestLatancy.With(lvs...).Observe(time.Since(start).Seconds())
	}(time.Now())

	res, err = mw.next.Concat(ctx, a, b)
	return
}
