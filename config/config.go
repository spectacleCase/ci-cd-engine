package config

import (
	"github.com/spf13/viper"
	"os"
)

var Config *Conf

// Conf 全局配置
type Conf struct {
	System *System
	Zap    *Zap
	Mysql  *Mysql
}

// InitConfig 初始化配置
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}
