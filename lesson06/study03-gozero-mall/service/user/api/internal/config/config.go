package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	MysqlDb struct{
		DbSource string `json:"DbSource"`
	}

	CacheRedis cache.CacheConf

	Auth struct {// JWT 认证需要的密钥和过期时间配置
        AccessSecret string
        AccessExpire int64
    }
}
