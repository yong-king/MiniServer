package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	MysqlDb struct{
		DbSource string `json:"DbSource"`
	}
}
