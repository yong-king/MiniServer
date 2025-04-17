package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"com.ysh.kit/demo/endpoint"
	"com.ysh.kit/demo/pb"
)

// Transports
// 解码
func decodeSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.SumResquest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}


// grpc 
// grpctransport.Handler 接口实现了 --》 ServeGRPC --> 是Server结构体的方法 --> 由NewServer返回
// func (s Server) ServeGRPC(ctx context.Context, req interface{}) (retctx context.Context, resp interface{}, err error)
// func NewServer(...) *Server 
func decodeGRPCSumRequest(_ context.Context, grpcReq interface{}) (interface{}, error){
	req := grpcReq.(*pb.SumRequest)
	return endpoint.SumResquest{A: int(req.A), B: int(req.B)}, nil
}

func encodeGRPCSumResponse(_ context.Context, response interface{}) (interface{}, error){
	resp := response.(endpoint.SumResponse)
	return &pb.SumResponse{V: int64(resp.V), Err: resp.Err}, nil
}