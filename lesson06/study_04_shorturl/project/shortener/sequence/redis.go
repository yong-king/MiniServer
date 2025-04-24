package sequence

import (
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const sequenceKey = "sequence:id"

type Redis struct{
	rds *redis.Redis
}

func NewRedis(conf redis.RedisConf) *Redis{
	return &Redis{
		rds: redis.MustNewRedis(conf),
	}
}

// Next 获取下一个字增的序列号
func (r *Redis) Next() (seq uint64, err error){
	// 使用 INCR 命令自增序列号
	var val int64
	// val, err = r.conn.Incr(context.Background(), sequenceKey).Result()
	val, err = r.rds.Incr(sequenceKey)
	if err != nil {
		logx.Errorw("r.conn.Incr failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	// 返回自增后的序列号
	return uint64(val), nil
}

// 实现 LinkMapping 接口的 SetShortLink 方法，存储短链接和长链接的映射
func (r *Redis) SetShortLink(shortLink, longLink string) error {
    err := r.rds.Set("shortlink:"+shortLink, longLink)
    if err != nil {
        logx.Errorw("rds.Set failed", logx.LogField{Key: "err", Value: err.Error()})
        return err
    }
    return nil
}

// 实现 LinkMapping 接口的 GetLongLink 方法，查询短链接对应的长链接
func (r *Redis) GetLongLink(shortLink string) (string, error) {
	if shortLink == "" {
		return "", errors.New("need shortlink")
	}
    longLink, err := r.rds.Get("shortlink:" + shortLink)
    if err != nil {
        logx.Errorw("rds.Get failed", logx.LogField{Key: "err", Value: err.Error()})
        return "", err
    }

    if longLink == "" {
        return "", fmt.Errorf("short link not found")
    }

    return longLink, nil
}