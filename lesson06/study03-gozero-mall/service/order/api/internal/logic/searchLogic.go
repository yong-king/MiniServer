package logic

import (
	"context"
	"errors"
	"strconv"
	"user/rpc/pb/user"

	"order/api/internal/interceptor"
	"order/api/internal/svc"
	"order/api/internal/types"
	"order/api/internal/errorx"


	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	// todo: add your logic here and delete this line
	// 1. 根据请求参数中的订单号查询数据库找到订单记录
	orderId, err := strconv.ParseUint(req.OrderID, 10, 64)
	if err != nil {
		logx.Errorw("order.Search.strconv.ParseInt falied", logx.Field("err", err))
		return nil, err
	}
	order, err := l.svcCtx.OrderModel.FindOneByOrderId(l.ctx, orderId)
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, errorx.NewCodeError(errorx.QuerNoFoundErrorCode, "内部错误")
	}
	if err != nil {
		logx.Errorf("order.Search.svcCtx.OrderModel.FindOne query failed: %#v\n", err)
		return nil, errorx.NewCodeError(errorx.SqlErrorCode, "内部错误")
	}
	// 2. 根据订单记录中的 user_id 查询用户数据（通过调用PRC的user服务）
	l.ctx = context.WithValue(l.ctx, interceptor.CtxKeyAdmindID, "666") // 通过上下文传入数据

	userResp, err := l.svcCtx.UserRPC.GetUser(l.ctx, &user.GetUserReq{UserID: int64(order.UserId)})
	if err != nil {
		logx.Errorw("order.Search.svcCtx.UserRPC.GetUser failed", logx.Field("err", err))
		return nil, errorx.NewCodeError(errorx.RpcErroCode, "内部错误")
	}
	// 3. 拼接放回结果（因为我们这个接口的数据不是由我一个服务组成的）
	return &types.SearchResponse{
		OrderID: strconv.FormatUint(orderId, 10),
		Status: int(order.Status),
		Username: userResp.Username,
		TradeID: order.TradeId,
		PayChannel: int(order.PayChannel),
		PayAmount: int(order.PayAmount),
		PayTime: int(order.PayTime.Unix()),
	}, nil
}
