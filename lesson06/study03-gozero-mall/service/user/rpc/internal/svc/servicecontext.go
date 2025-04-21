package svc

import (
	"user/rpc/internal/config"
	"user/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MysqlDb.DbSource)
	return &ServiceContext{
		Config: c,
		UserModel: model.NewUserModel(conn, c.CacheRedis),
	}
}
