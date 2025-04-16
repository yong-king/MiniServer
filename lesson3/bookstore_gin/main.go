package main

import (
	"com.bookstore/demo/dao/mysql"
	"com.bookstore/demo/logger"
	"com.bookstore/demo/router"
	"com.bookstore/demo/settings"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

func main() {
	// 配置文件路径
	// 配置文件路径
	var fileName string
	flag.StringVar(&fileName, "configPath", "./conf/config.yaml", "配置文件路径！")
	flag.Parse()

	// 1.获取配置信息 vipper
	err := settings.Init(fileName)
	if err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	// 2.设置日志文件
	err = logger.Init(settings.Conf.LogConfig)
	if err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	// 3.mysql初始化
	err = mysql.Init(settings.Conf.MySQLConfig)
	if err != nil {
		zap.L().Error("init mysql failed", zap.Error(err))
		return
	}
	defer mysql.Close()

	r := router.Setup(settings.Conf.Mode)
	srv := &http.Server{
		Addr:    settings.Conf.Port,
		Handler: r,
	}
	// 监听
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("listen failed", zap.Error(err))
			return
		}
	}()
	//// 7.优雅关机
	//quit := make(chan os.Signal, 1)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
	//zap.L().Info("Shutdown Server ...")
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//err = srv.Shutdown(ctx)
	//if err != nil {
	//	zap.L().Error("Server Shutdown Failed", zap.Error(err))
	//	return
	//}
}
