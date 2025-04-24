package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	"shortener/pkg/urltool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 长链接转短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 1. 参数校验
	// 1.1 长链接不能为空
	// validate校验,在路由中处理，如果路由都不能通过，则不会进来

	// 1.2 长链接必须是能连通的
	// http.Get()
	if ok := connect.Get(req.LongUrl); !ok {
		return nil, errors.New("无效链接")
	}

	// 1.3 判断之前是否已经转链过
	// 1.3.1 长链生成md5
	md5Value := md5.Sum([]byte(req.LongUrl))
	// 1.3.2 长链生成的md5是否在数据库中
	u, err := l.svcCtx.ShortUrlDB.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sqlx.ErrNotFound{
		if err == nil {
			return nil, fmt.Errorf("该链接已被转为%s", u.Surl.String)
		}
		logx.Errorw("ShortUrlModel.FindOneByMd5 failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}

	// 1.4 输入的不能是一个短链接（循环转链）
	basePath, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("urltool.GetBasePath failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	_, err = l.svcCtx.ShortUrlDB.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound{
		if err == nil {
			return nil, errors.New("该链接已经是短链接了")
		}
		logx.Errorw("ShortUrlDB.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	// 2. 取号
	// 每来一个转链请求，就使用 REPLACE INTO语句往sequence表中一条数据，并去取出主键id为号码
	// 两种，mysql取号和redis取号
	var short string
	for {
		seq, err := l.svcCtx.SequenceDB.Next()
		if err != nil {
			logx.Errorw("SequenceDB.Next failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, err
		}
		fmt.Println("--->", seq)
		// 3. 号码转链 转为62进制，可以缩短长度
		short = base62.Int2String(seq)
		// 3.1 移除出现的敏感词，如api,convert,fuck,shirt...
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
			break // 生成不在黑名单里的短链接就跳出for循环
		}
	}
	// 4. 存储长短链接映射关系
	l.svcCtx.ShortUrlDB.Insert(l.ctx, &model.ShortUrlMap{
		Lurl: sql.NullString{String: req.LongUrl, Valid: true},
		Md5: sql.NullString{String: md5Value, Valid: true},
		Surl: sql.NullString{String: short, Valid: true},
	})
	// 4.2将生成的短链接加入到布隆过滤器
	err = l.svcCtx.Filter.Add([]byte(short))
	if err != nil {
		logx.Errorw("Filter.Add failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	// err = l.svcCtx.LinkMappingDB.SetShortLink(short, req.LongUrl)
	// if err != nil{
	// 	logx.Errorw("LinkMappingDB.SetShortLink failed", logx.LogField{Key: "err", Value: err.Error()})
	// 		return nil, err
	// }
	// 5. 返回响应
	// 5.1 返回的 域名 + 短链接 
	shortUrl := l.svcCtx.ShortDomain + "/" + short
	return &types.ConvertResponse{ShortUrl: shortUrl}, nil
}
