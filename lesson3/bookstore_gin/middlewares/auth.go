package middlewares

import (
	"com.bookstore/demo/controller"
	"com.bookstore/demo/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

func JwtAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization
		// Authorization: Bearer xxx.xxx.xxx
		authHead := c.Request.Header.Get("Authorization")
		if authHead == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			return
		}
		// 按空格划分
		parse := strings.SplitN(authHead, " ", 2)
		if len(parse) != 2 && parse[0] != "Bearer" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// parse[1]为token,解析验证token
		mc, err := jwt.ParseToken(parse[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}

		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next()
	}
}
