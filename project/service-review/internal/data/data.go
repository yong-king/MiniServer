package data

import (
	"errors"
	"service-review/internal/conf"
	"service-review/internal/data/query"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewReviewRepo, NewDB, NewEsclient, NewRedisClient)

// Data .
type Data struct {
	query *query.Query
	log *log.Helper
	es *elasticsearch.TypedClient
	rdb *redis.Client
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger, esClient *elasticsearch.TypedClient, rdb *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 非常重要!为GEN生成的query代码设置数据库连接对象
	query.SetDefault(db)
	return &Data{query: query.Q, log: log.NewHelper(logger), es: esClient, rdb: rdb}, cleanup, nil
}

func NewEsclient(cfg *conf.Elasticsearch) (*elasticsearch.TypedClient, error) {
	// ES 配置
	c := elasticsearch.Config{
		Addresses: cfg.GetAddresses(),
	}
	// 创建客户端
	return elasticsearch.NewTypedClient(c)
}

// NewDB 数据库连接
func NewDB(cfg *conf.Data) (*gorm.DB, error) {
	switch strings.ToLower(cfg.Database.Driver){
	case "mysql":
		return gorm.Open(mysql.Open(cfg.Database.Source))
	case "sqlite":
		return gorm.Open(sqlite.Open(cfg.Database.Source))
	}
	return nil, errors.New("connect db failed unsuppoesd db driver")
}

// NewRedisClient redis连接
func NewRedisClient(cfg *conf.Data) *redis.Client{
	return redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
		WriteTimeout: cfg.Redis.WriteTimeout.AsDuration(),
		ReadTimeout: cfg.Redis.ReadTimeout.AsDuration(),
	})
}