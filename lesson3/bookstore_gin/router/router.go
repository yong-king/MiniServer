package router

import (
	"com.bookstore/demo/controller"
	"com.bookstore/demo/logger"
	"com.bookstore/demo/middlewares"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化
	router := gin.New()
	// 插入中间件
	router.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := router.Group("v1/bookstore")
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 获取书架列表
	v1.GET("/shelvs", controller.ShelvesListHandler)
	// 获取指定书架
	v1.GET("/shelf/:id", controller.GetShelfHandler)
	// 获取书架的书籍列表
	v1.GET("/shelf/:sid/book", controller.ListBookHandler)
	// 获取指定书籍
	v1.GET("/shelf/:sid/book/:id", controller.GetBookHandler)

	// 插入中间件
	v1.Use(middlewares.JwtAuthMiddleWare())
	{ // 创建书架
		v1.POST("/shelf", controller.CreateShelfHandler)
		// 删除书架
		v1.DELETE("shelf/:id", controller.DelateShelfHandler)
		// 创建书籍
		v1.POST("/shelf/:sid/book", controller.CreateBookHandler)
		// 删除书籍
		v1.DELETE("/shelf/book/:id", controller.DeleteBookHandler)
	}
	pprof.Register(router)

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	return router
}
