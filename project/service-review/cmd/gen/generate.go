package main

import (
	"errors"
	"flag"
	"fmt"
	"service-review/internal/conf"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gorm.io/gen"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}


func connectDB(cfg *conf.Data_Database) *gorm.DB {
	if cfg == nil{
		panic(errors.New("GEN: connectDB fail, need cfg"))
	}
	switch strings.ToLower(cfg.GetDriver()){
	case "mysql":
		db, err := gorm.Open(mysql.Open(cfg.GetSource()))
		if err != nil{
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db
	case "sqlter":
		db, err := gorm.Open(sqlite.Open(cfg.GetSource()))
		if err != nil{
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db
	}
	panic(errors.New("GEN:connectDB fail unsupported db driver"))

}

func main() {

	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "../../internal/data/query",
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true, // delete_at是可以为空的
	})
	// 通常复用项目中已有的SQL连接配置db(*gorm.DB)
	// 非必需，但如果需要复用连接时的gorm.Config或需要连接数据库同步表信息则必须设置
	g.UseDB(connectDB(bc.Data.Database))

	// 从连接的数据库为所有表生成Model结构体和CRUD代码
	// 也可以手动指定需要生成代码的数据表
	g.ApplyBasic(g.GenerateAllTable()...)

	// 执行
	g.Execute()
}