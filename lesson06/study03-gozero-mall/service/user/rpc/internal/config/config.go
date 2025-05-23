package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf

	// mysql
	MysqlDb struct{
		DbSource string `json:"DbSource"`
	}

	// redis
	CacheRedis cache.CacheConf

	Consul consul.Conf
}
