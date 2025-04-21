package logic

import (
	"context"
	"errors"
	"time"

	"user-api/internal/svc"
	"user-api/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line
	// 1. 获取用户名和密码
	userName, password := req.Username, req.Password

	// 2. 到数据库查询
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, userName)
	if err == sqlx.ErrNotFound{
		return &types.LoginResponse{Message: "用户不存在!"}, nil
	}
	if err != nil{
		logx.Errorf("user_login_UserModel.FindOneByUsername query failed: %#v\n", err)
		return &types.LoginResponse{Message: "用户不存在!"}, errors.New("数据库内部错误!")
	}

	// jwt
	now := time.Now().Unix()
	expire := l.svcCtx.Config.Auth.AccessExpire
	token, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, expire, user.UserId)
	if err != nil {
		logx.Errorf("user_login_getJwtToken failed :%#v\n", err)
		return nil, err
	}
	
	// 3. 对获取到密码进行比较
	// 3.1 对密码进行加密
	// 3.2 比较
	password = passwordMd5([]byte(password))
	if password == user.Password {
		return &types.LoginResponse{
			Message: "登录成功!",
			 AccessToken: token, 
			 AccseeExpire: int(expire), 
			 RefreshAfter: int(now+expire/2),
			}, nil
	}
	return &types.LoginResponse{Message: "密码错误!"}, nil
}

// 生成JWT方法
// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func (l *LoginLogic)getJwtToken(secretKey string, iat, seconds int64, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["auth"] = "ysh"
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
  }