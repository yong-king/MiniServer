package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"com.ysh.kit/demo/service"
)

// 1.3 请求和响应
type SumResquest struct {
	A int `json:"a"`
	B int `json:"b"`
}


type SumResponse struct {
	V   int    `json:"v"`
	Err string `json:"err, omitempty"`
}


// Endpoints
// type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
func MakeSumEndpoint(srv service.AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error){
		// 类型断言
		req := request.(SumResquest)
		v, err := srv.Sum(ctx, req.A, req.B)
		if err != nil {
			return SumResponse{V: v, Err: err.Error()}, nil
		}
		return SumResponse{V: v}, nil
	}
}