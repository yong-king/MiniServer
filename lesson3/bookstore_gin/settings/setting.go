package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Version      string `mapstructure:"version"`
	Mode         string `mapstructure:"mode"`
	Port         string `mapstructure:"port"`
	*MySQLConfig `mapstructure:"mysql"`
	*LogConfig   `mapstructure:"log"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type LogConfig struct {
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Level      string `mapstructure:"level"`
}

var Conf = new(AppConfig)

func Init(filename string) (err error) {
	viper.SetConfigName(filename)

	err = viper.ReadInConfig()
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}

	// 反序列化到结构体
	err = viper.Unmarshal(Conf)
	if err != nil {
		fmt.Printf("viper.Unmarshal() failed, err:%v\n", err)
		return
	}
	// 热监督
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件被修改!", e.Name)
		// 发现修改后重新反序列化
		err = viper.Unmarshal(Conf)
		if err != nil {
			fmt.Printf("viper.Unmarshal() failed, err:%v\n", err)
			return
		}
	})
	return
}
