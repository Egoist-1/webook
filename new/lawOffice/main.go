package main

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

func main() {
	initViper()
	initZapLog()
	app := InitWebServer()
	app.Server.Run(viper.GetString("server.host"))
}

func initViper() {
	// 配置文件的名字，但是不包含文件扩展名
	// 不包含 .go, .yaml 之类的后缀
	viper.SetConfigName("dev")
	// 告诉 viper 我的配置用的是 yaml 格式
	// 现实中，有很多格式，JSON，XML，YAML，TOML，ini
	viper.SetConfigType("yaml")
	// 当前工作目录下的 config 子目录
	viper.AddConfigPath("./config")
	// 读取配置到 viper 里面，或者你可以理解为加载到内存里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func initZapLog() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
