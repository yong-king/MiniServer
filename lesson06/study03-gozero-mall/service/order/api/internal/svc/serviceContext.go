package svc

import (
	"order/api/internal/config"
	"order/api/internal/interceptor"
	"order/model"
	"user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	OrderModel model.OrderModel

	UserRPC userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MysqlDb.DbSource)
	return &ServiceContext{
		Config: c,
		OrderModel: model.NewOrderModel(conn, c.CacheRedis),
		UserRPC: userclient.NewUser(
			zrpc.MustNewClient(
				c.UserRPC, 
				zrpc.WithUnaryClientInterceptor(interceptor.YshInterceptor),
			),
		),
	}
}


