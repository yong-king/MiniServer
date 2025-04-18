package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"user-api/internal/svc"
	"user-api/internal/types"
	"user-api/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

var secret  = []byte("ysh&dlrb forever")

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignupLogic) Signup(req *types.SignupRequest) (resp *types.SingupResponse, err error) {
	// todo: add your logic here and delete this line
	// 添加业务逻辑
	// 1.参数校验
	if req.Password != req.RePassword {
		return nil, errors.New("两次密码不正确!")
	}

	// 2. 注册用户存在
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil && err != sqlx.ErrNotFound{
		return nil, errors.New("查询出错!")
	}
	if u != nil {
		return nil, errors.New("注册用户以存在")
	}

	// 3. 密码加密
	h := md5.New()
	h.Write([]byte(req.Password))
	h.Write(secret)
	passwordS := hex.EncodeToString(h.Sum(nil))

	fmt.Printf("resq: %#v\n", req)
	resp = &types.SingupResponse{
		Message: "success",
	}
	user := &model.User{
		UserId:   time.Now().Unix(),
		Username: req.Username,
		Password: passwordS,
		Gender:   1,
	}
	_, err = l.svcCtx.UserModel.Insert(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
