package logic

import (
	"context"

	"shorurl/internal/svc"
	"shorurl/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShorurlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShorurlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShorurlLogic {
	return &ShorurlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShorurlLogic) Shorurl(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	// 将请求的端连接转换为长连接
	// ysh&dlrb --> https://baidu.com
	if req.ShortURL == "ysh&dlrb" {
		return &types.Response{LongURL: "https://baidu.com"}, nil
	}
	// 如果查询不到，就跳转为https://google.com
	return &types.Response{LongURL: "https://google.com"}, nil
}
