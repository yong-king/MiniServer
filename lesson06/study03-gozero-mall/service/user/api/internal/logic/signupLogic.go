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

	logx.Debugf("req:%#v\n", req)

	// 2. 注册用户存在
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil && err != sqlx.ErrNotFound{
		fmt.Printf("err:%#v\n",err)
		logx.Errorf("user_signup_UserModel.FindOneByUsername failed: %#v\n", err)
		return nil, errors.New("查询出错!")
	}
	if u != nil {
		return nil, errors.New("注册用户以存在")
	}

	// 3. 密码加密
	passwordS := passwordMd5([]byte(req.Password))
	
	fmt.Printf("resq: %#v\n", req)

	user := &model.User{
		UserId:   time.Now().Unix(),
		Username: req.Username,
		Password: passwordS,
		Gender:   1,
	}

	// 插入数据
	_, err = l.svcCtx.UserModel.Insert(context.Background(), user)
	if err != nil {
		logx.Errorw(
			"user_signup_UserModel.Insert failed",
			logx.Field("err:", err),
		)
		return nil, err
	}

	return &types.SingupResponse{Message: "success"}, nil
}

// 密码加密
func passwordMd5(password []byte) string {
	h := md5.New()
	h.Write([]byte(password))
	h.Write(secret)
	return hex.EncodeToString(h.Sum(nil))
}


