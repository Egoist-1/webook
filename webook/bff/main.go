package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	initViper()
	initLogger()
	app := InitApp()
	var s = viper.GetString("web.port")
	fmt.Println(s)
	app.web.Run(viper.GetString("web.port"))
}

func initViper() {
	viper.SetConfigName("dev")      // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // path to look for the config file in
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	// 设置了全局的 logger，
	// 你在你的代码里面就可以直接使用 zap.XXX 来记录日志
	zap.ReplaceGlobals(logger)
}
