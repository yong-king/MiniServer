package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"

	"com.ysh.kit/demo/endpoint"
	"com.ysh.kit/demo/middleware"
	"com.ysh.kit/demo/pb"
	"com.ysh.kit/demo/service"
)

func NewHTTPServer(svc service.AddService, logger log.Logger) http.Handler {
	sum := endpoint.MakeSumEndpoint(svc)
	// 日志记录中间件
	sum = middleware.LoggingMiddleware(log.With(logger, "method", "sum"))(sum)
	// 限流中间件
	sum = middleware.RateMiddleware(rate.NewLimiter(1,1))(sum)
	
	sumHandler := httptransport.NewServer(
		sum,
		decodeSumRequest,
		encodeResponse,
	)

	concat := endpoint.MakeconcatEndpoint(svc)
	concat = middleware.LoggingMiddleware(log.With(logger, "method", "concat"))(concat)
	concatHandler := httptransport.NewServer(
		concat,
		decodeConcatRequest,
		encodeResponse,
	)

	// github.com/gotilla/mux
	// r:=mux.NewRouter()
	// http.Handle("/sum", sumHandler).Methods("POST")
	// http.Handle("/concat", concatHandler).Methods("POST")

	// gin
	r := gin.Default()
	r.POST("/sum", gin.WrapH(sumHandler))
	r.POST("/concat", gin.WrapH(concatHandler))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	return r
}

func NewGRPCServer(svc service.AddService) pb.AddServer {
	return &service.GrpcServer{
		GRPCsum: grpctransport.NewServer(
			endpoint.MakeSumEndpoint(svc),
			decodeGRPCSumRequest,
			encodeGRPCSumResponse,
		),
		GRPCconcat: grpctransport.NewServer(
			endpoint.MakeconcatEndpoint(svc),
			decodeGRPCConcatRequest,
			encodeGRPCConcatResponse,
		),
	}
}
