package trim

import (
	"context"
	"io"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	sdconsul "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"com.ysh.kit/demo/service"
)

// trim
type withTrimMiddleware struct {
	next service.AddService
	trimService endpoint.Endpoint
}

func NewwithTrimMiddleware(srv service.AddService, trimService endpoint.Endpoint) service.AddService{
	return &withTrimMiddleware{
		next: srv,
		trimService: trimService,
	}
}

func (mw withTrimMiddleware) Sum(ctx context.Context, a, b int) (res int, err error) {
	return mw.next.Sum(ctx, a, b)
}


func (mw withTrimMiddleware) Concat(ctx context.Context, a, b string) (res string, err error) {
	resqA, err := mw.trimService(ctx, trimRequest{s: a})
	if err != nil {
		return "", err
	}
	resqB, err := mw.trimService(ctx, trimRequest{s: b})
	if err != nil {
		return "", err
	}
	trimA := resqA.(trimResponse)
	trimB := resqB.(trimResponse)
	return mw.next.Concat(ctx, trimA.s, trimB.s)

}

func GetTrimServiceFromConsul(consulAddr string, logger log.Logger, srvName string, tags []string) (endpoint.Endpoint, error) {
	// 1. 连接consul
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = consulAddr

	cc, err := consulapi.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	// 2. 使用go kit 提供的适配器
	sdClient := sdconsul.NewClient(cc)
	var passiongOnly = true
	instancer := sdconsul.NewInstancer(sdClient, logger, srvName, tags, passiongOnly)

	// 3. endpointer
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	// 4. balancer
	balancer := lb.NewRoundRobin(endpointer)
	// 5. retry
	retry := lb.Retry(3, time.Second, balancer)
	return retry, nil
}

func factory(instance string) (endpoint.Endpoint, io.Closer, error){
	conn, err := grpc.NewClient(instance, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	e := MakeTrimEndpoint(conn)
	return e, conn, nil
}