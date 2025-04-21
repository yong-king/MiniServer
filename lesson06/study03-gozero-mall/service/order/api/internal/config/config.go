package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	MysqlDb struct{
		DbSource string `json:"DbSource"`
	}

	CacheRedis cache.CacheConf

	UserRPC zrpc.RpcClientConf	// 连接其他微服务的RPC客户端
}
