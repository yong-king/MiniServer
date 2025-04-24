package logic

import (
	"context"
	"database/sql"
	"errors"

	"shortener/internal/svc"
	"shortener/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowResquest) (resp *types.ShowResponse, err error) {
	// todo: add your logic here and delete this line
	// 1. 根据短链到数据库查询原始长链
	// ul, err := l.svcCtx.LinkMappingDB.GetLongLink(req.ShortUrl)
	// if err != nil {
	// 	if err == sqlx.ErrNotFound{
	// 		return nil, errors.New("404")
	// 	}
	// 	logx.Errorw("ShortUrlDB.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
	// 	return nil, err
	// }
	// fmt.Println("----->ul", ul)
	// 1.0 布隆过滤器
	// 不存在的短链接直接返回404,不需要后续处理
	// a. 基于内存版本
	// b. 基于redis版本
	exist, err := l.svcCtx.Filter.Exists([]byte(req.ShortUrl))
	if err != nil {
		logx.Errorw("Filter.Exists failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	if !exist{
		return nil, errors.New("404")
	}
	u, err := l.svcCtx.ShortUrlDB.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
	// fmt.Println("----->u", u.Lurl.String)
	if err != nil {
		if err == sqlx.ErrNotFound{
			return nil, errors.New("404")
		}
		logx.Errorw("ShortUrlDB.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	// 2. 返回查询到长链，在调用handler层返回重定位响应
	return &types.ShowResponse{LongUrl: u.Lurl.String}, nil
}
