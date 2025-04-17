package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"com.ysh.kit/demo/service"
)

type ConcatRequest struct {
	A string `json:"a"`
	B string `jsoin:"b"`
}


type ConcatResponse struct {
	V   string `json:"v"`
	Err string `json:"err, omitempty"`
}


func MakeconcatEndpoint(srv service.AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error){
		// 类型断言
		req := request.(ConcatRequest)
		v, err := srv.Concat(ctx, req.A, req.B)
		if err != nil {
			return ConcatResponse{V: v, Err: err.Error()}, nil
		}
		return ConcatResponse{V: v}, nil
	}
}
