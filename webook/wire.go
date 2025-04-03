//go:build wireinject

package main

import (
	"github.com/google/wire"
	"start/webook/internal/repository"
	"start/webook/internal/repository/cache"
	"start/webook/internal/repository/dao"
	"start/webook/internal/service"
	"start/webook/internal/service/sms/memory"
	"start/webook/internal/web"
	"start/webook/ioc"
)

var user = wire.NewSet(
	web.NewUserHandle,
	service.NewUserServiceImpl,
	repository.NewUserRepo,
	dao.NewUserDao,
)
var code = wire.NewSet(
	service.NewCodeService,
	repository.NewCodeRepo,
	cache.NewCodeCacheRedis,
)
var sms = wire.NewSet(
	memory.NewMemory,
)

func InitApp() *App {
	wire.Build(
		user,
		code,
		sms,
		ioc.InitWebServer,
		ioc.InitGorm,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
