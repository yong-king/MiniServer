package middleware

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"golang.org/x/time/rate"
)

// type Middleware func(Endpoint) Endpoint
// type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error){
			logger.Log("msg", "开始调用")
			defer logger.Log("msg", "调用结束")
			return next(ctx, request)
		}
	}
}


var (
	ErrRateLimit = errors.New("request rate limit")
)
// ratelimit
func RateMiddleware(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint{
		return func (ctx context.Context, request interface{}) (interface{}, error)  {
			if !limit.Allow() {
				return nil, ErrRateLimit
			}
			return next(ctx, request)
		}
	}
}