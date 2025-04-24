package svc

import (
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 空结构体 struct{} 没有字段，不占用任何内存。即使你有上万个 URL 存储在 map 中，使用 struct{} 作为值类型也不会增加内存占用。

type ServiceContext struct {
	Config config.Config
	ShortUrlDB model.ShortUrlMapModel
	SequenceDB sequence.Sequence
	ShortUrlBlackList map[string]struct{}
	ShortDomain string

	LinkMappingDB sequence.LinkMapping

	Filter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlMapDb.DSN)
	// 把配置文件中配置的黑名单加载到map，方便后续判断
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList{
		m[v] = struct{}{}
	}

	// 初始化布隆过滤器
	store := redis.New(c.CacheRedis[0].Host)
	filter := bloom.New(store, "bloom_filter", 20*(1<<20))
	return &ServiceContext{
		Config: c,
		ShortUrlDB: model.NewShortUrlMapModel(conn, c.CacheRedis),
		// SequenceDB: sequence.NewMysql(c.SequenceDB.DNS),
		SequenceDB: sequence.NewRedis(c.SequenceRDB),
		ShortUrlBlackList: m,
		ShortDomain: c.ShortDomain,
		LinkMappingDB: sequence.NewRedis(c.SequenceRDB),
		Filter: filter,
	}
}



// import (
// 	"errors"

// 	"github.com/bits-and-blooms/bloom/v3"
// 	"github.com/zeromicro/go-zero/core/logx"
// 	"github.com/zeromicro/go-zero/core/stores/sqlx"
// )



// var filter = bloom.NewWithEstimates(1<<20, 0.01)

// func loadDataToBloomFilter(conn sqlx.SqlConn, filter *bloom.BloomFilter) error {
// 	if conn == nil || filter == nil{
// 		return errors.New("loadDataToBloomFilter invalid param")
// 	}

// 	// 查总数
// 	total := 0
// 	if err := conn.QueryRow(&total, "select count(*) from short_url_map where is_del=0"); err != nil{
// 		logx.Errorw(" conn.QueryRow failed", logx.LogField{Key: "err", Value: err.Error()})
// 		return err
// 	}
// 	logx.Infow("total data", logx.LogField{Key: "total", Value: total})
// 	if total == 0 {
// 		logx.Info("no data need to load")
// 		return nil
// 	}

// 	pageTotal := 0
// 	pageSize := 20
// 	if total%pageSize == 0{
// 		pageTotal = total/pageSize
// 	}else {
// 		pageTotal = total/pageSize + 1
// 	}
// 	logx.Infow("pageTotal", logx.LogField{Key: "pageTotal", Value: pageTotal})

// 	// 循环查询所有数据
// 	for page := 1; page <= pageTotal; page++ {
// 		offset := pageSize * (pageTotal - 1)
// 		surls := []string{}
// 		if err := conn.QueryRow(&surls, "select surl from short_url_map where is_del=0 limit ?,?", offset, pageSize); err != nil {
// 			return err
// 		}
// 		for _, surl := range surls{
// 			filter.AddString(surl)
// 		}
// 	}
// 	logx.Info("load data to bloom success")
// 	return nil
// }