package logic

import (
	"context"
	"errors"
	"fmt"

	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	// todo: add your logic here and delete this line
	fmt.Printf("user_id:%v\n", l.ctx.Value("userId"))
	fmt.Printf("auth:%v\n", l.ctx.Value("auth"))

	// 1. 获取用户的user_id参数
	userId := req.UserId
	fmt.Printf("userId:%d\n",userId)

	// 2. 通过user_id查询
	user,  err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, userId)
	if err == sqlx.ErrNotFound{
		return &types.UserInfoResponse{Message: "用户不存在"}, nil
	}
	if err != nil {
		logx.Errorf("user_userInfo_UserModel.FindOneByUserId quiry failed: %#v\n", err)
		return nil, errors.New("内部错误!")
	}
	return &types.UserInfoResponse{UserName: user.Username, Gender: int(user.Gender), Message: "success"}, nil
}
