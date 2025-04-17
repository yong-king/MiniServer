package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"com.ysh.kit/demo/endpoint"
	"com.ysh.kit/demo/pb"
)

func decodeConcatRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.ConcatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// grpc
func decodeGRPCConcatRequest(_ context.Context, grpcReq interface{}) (interface{}, error){
	req := grpcReq.(*pb.ConcatREquest)
	return endpoint.ConcatRequest{A: req.A, B: req.B}, nil
}


func encodeGRPCConcatResponse(_ context.Context, response interface{}) (interface{}, error){
	resp := response.(endpoint.ConcatResponse)
	return &pb.ConcatResponse{V: resp.V, Err: resp.Err}, nil
}