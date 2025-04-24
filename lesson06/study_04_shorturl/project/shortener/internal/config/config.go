package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// sql
	ShortUrlMapDb ShortUrlMapDb
	SequenceDB struct{
		DSN string
	}

	// redis
	CacheRedis cache.CacheConf

	SequenceRDB redis.RedisConf

	Base62String string
	ShortUrlBlackList []string
	ShortDomain string
}

type ShortUrlMapDb struct{
	DSN string
}
