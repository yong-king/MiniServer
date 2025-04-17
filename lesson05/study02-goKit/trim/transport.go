package trim

import (
	"context"

	"com.ysh.kit/demo/pb"
)

func encodeTrimRequest(ctx context.Context, request interface{}) (interface{}, error){
	req := request.(trimRequest)
	return &pb.TrimRequest{S: req.s}, nil
}

func decodeTrimResponse(ctx context.Context, response interface{}) (interface{}, error){
	resq := response.(*pb.TrimResponse)
	return trimResponse{s: resq.S}, nil
}