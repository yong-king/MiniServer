package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"

	"com.ysh.kit/demo/pb"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	// httptransport "github.com/go-kit/kit/transport/http"
)

// 1.业务逻辑层

// 1.1 定义接口
type AddService interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(cte context.Context, a, b string) (string, error)
}

// 1.2实现接口
type addService struct{}

var (
	// ErrTwoEmptyStrings 两个字符串都为空！
	ErrTwoEmptyStrings = errors.New("两个字符串都为空！")
)

// Sum 两数字相加
func (s addService) Sum(_ context.Context, a, b int) (int, error) {
	return a + b, nil
}

// Concat 字符串拼接
func (s addService) Concat(_ context.Context, a, b string) (string, error) {
	if a == "" && b == "" {
		return "", ErrTwoEmptyStrings
	}
	return a + b, nil
}

// 1.3 请求和响应
type SumResquest struct {
	A int `json:"a"`
	B int `json:"b"`
}

type ConcatRequest struct {
	A string `json:"a"`
	B string `jsoin:"b"`
}

type SumResponse struct {
	V   int    `json:"v"`
	Err string `json:"err, omitempty"`
}

type ConcatResponse struct {
	V   string `json:"v"`
	Err string `json:"err, omitempty"`
}


// Endpoints
// type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
func makeSumEndpoint(srv AddService) endpoint.Endpoint {
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

func makeconcatEndpoint(srv AddService) endpoint.Endpoint {
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

// Transports
// 解码
func decodeSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request SumResquest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeConcatRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ConcatRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// 编码
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}


// grpc
// Handler 应该从服务实现的gRPC绑定调用。
// 传入的请求参数和返回的响应参数都是gRPC类型，而不是用户域类型。
// type Handler interface {
// 	ServeGRPC(ctx context.Context, request interface{}) (context.Context, interface{}, error)
// }
type grpcServer struct{
	pb.UnimplementedAddServer

	sum grpctransport.Handler
	concat grpctransport.Handler
}

func (s *grpcServer) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	_, rep, err := s.sum.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SumResponse), nil
}

func  (s *grpcServer) Concat(ctx context.Context, req *pb.ConcatREquest) (*pb.ConcatResponse, error) {
	_, rep, err := s.concat.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ConcatResponse), nil
}


// grpctransport.Handler 接口实现了 --》 ServeGRPC --> 是Server结构体的方法 --> 由NewServer返回
// func (s Server) ServeGRPC(ctx context.Context, req interface{}) (retctx context.Context, resp interface{}, err error)
// func NewServer(...) *Server 
func decodeGRPCSumRequest(_ context.Context, grpcReq interface{}) (interface{}, error){
	req := grpcReq.(*pb.SumRequest)
	return SumResquest{A: int(req.A), B: int(req.B)}, nil
}

func decodeGRPCConcatRequest(_ context.Context, grpcReq interface{}) (interface{}, error){
	req := grpcReq.(*pb.ConcatREquest)
	return ConcatRequest{A: req.A, B: req.B}, nil
}

func encodeGRPCSumResponse(_ context.Context, response interface{}) (interface{}, error){
	resp := response.(SumResponse)
	return &pb.SumResponse{V: int64(resp.V), Err: resp.Err}, nil
}

func encodeGRPCConcatResponse(_ context.Context, response interface{}) (interface{}, error){
	resp := response.(ConcatResponse)
	return &pb.ConcatResponse{V: resp.V, Err: resp.Err}, nil
}

func NewGRPCServer(svc AddService) pb.AddServer {
	return &grpcServer{
		sum: grpctransport.NewServer(
			makeSumEndpoint(svc),
			decodeGRPCSumRequest,
			encodeGRPCSumResponse,
		),
		concat: grpctransport.NewServer(
			makeconcatEndpoint(svc),
			decodeGRPCConcatRequest,
			encodeGRPCConcatResponse,
		),
	}
}


func main() {
	svc := addService{}

	// sumHandler := httptransport.NewServer(
	// 	makeSumEndpoint(svc),
	// 	decodeSumRequest,
	// 	encodeResponse,
	// )

	// concatHandler := httptransport.NewServer(
	// 	makeconcatEndpoint(svc),
	// 	decodeConcatRequest,
	// 	encodeResponse,
	// )

	// http.Handle("/sum", sumHandler)
	// http.Handle("/concat", concatHandler)
	// http.ListenAndServe(":8080", nil)
	gs := NewGRPCServer(svc)
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterAddServer(s, gs)

	fmt.Println(s.Serve(lis))
}