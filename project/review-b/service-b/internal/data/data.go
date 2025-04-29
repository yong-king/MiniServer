package data

import (
	"context"
	v1 "service-b/api/review/v1"
	"service-b/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
    "github.com/hashicorp/consul/api"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewBusinessRepo, NewReviewServiceClient, NewDiscovever)

// Data .
type Data struct {
	// TODO wrapped database client
	rc  v1.ReviewClient
	log *log.Helper
}

// NewData .
func NewData(c *conf.Data, rc v1.ReviewClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{rc: rc, log: log.NewHelper(logger)}, cleanup, nil
}

func NewDiscovever(conf *conf.Registry) registry.Discovery{
	// new consul client
	c := api.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	// new dis with consul client
	dis := consul.New(client)
	return dis
}

func NewReviewServiceClient(d registry.Discovery) v1.ReviewClient {
	endpoint := "discovery:///review.service"
	conn, err := grpc.DialInsecure(context.Background(),
		// grpc.WithEndpoint("127.0.0.1:9001"),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.Validator(),
		))
	if err != nil {
		panic(err)
	}
	return v1.NewReviewClient(conn)
}
