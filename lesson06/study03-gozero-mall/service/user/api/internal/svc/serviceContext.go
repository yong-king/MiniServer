package svc

import (
	"user-api/internal/config"
	"user-api/internal/middleware"
	"user-api/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config config.Config
	//UserModel: 类型为 model.UsersModel，表示与用户相关的数据库模型
	//用于处理与用户相关的数据操作（如用户的创建、读取、更新和删除等）
	UserModel model.UserModel

	Cost rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		// UserModel 指针类型 --> userModel 指针类型
		// userModel 包含增删改查
		// defaultUserModel 实现了userModel接口
		// newUserModel 是 defaultUserModel 的构建方法
		//  NewUserModel(conn sqlx.SqlConn) UserModel
		//通过调用 model.NewUsersModel 函数对UserModel 进行初始化
		//sqlx.NewMysql 是数据库连接,链接字符串为config中的MysqlDb.DbSource
		UserModel: model.NewUserModel(sqlx.NewMysql(c.MysqlDb.DbSource), c.CacheRedis),
		Cost:      middleware.NewCostMiddleware().Handle,
	}
}
