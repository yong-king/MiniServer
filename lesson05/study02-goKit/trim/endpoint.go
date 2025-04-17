package trim

import (
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"com.ysh.kit/demo/pb"
)

type trimRequest struct {
	s string
}

type trimResponse struct {
	s string
}

// 使用grpctransport.NewClient 创建基于gRPC client的endpoint。
func MakeTrimEndpoint(cc *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		cc,
		"pb.Trim",          // 服务名
		"TrimSpace",        // 方法名
		encodeTrimRequest,  // 编码
		decodeTrimResponse, // 解码
		pb.TrimResponse{},	// 接收结果
	).Endpoint()
}


