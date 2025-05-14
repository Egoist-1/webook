package ioc

import (
	"go.uber.org/zap"
	"start/webook/pkg/logger"
)

func InitLogger() logger.Logger {
	cfg := zap.NewDevelopmentConfig()
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	zapLogger := logger.NewZapLogger(l)
	return zapLogger
}
