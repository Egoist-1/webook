//go:build wireinject

package main

import (
	"github.com/google/wire"
	ioc2 "start/webook/bff/ioc"
	"start/webook/bff/web"
	"start/webook/code/_internal/repository"
	"start/webook/code/_internal/repository/cache"
	"start/webook/code/_internal/service"
	"start/webook/sms/_internal/service/sms/memory"
	repository2 "start/webook/user/_internal/repository"
	"start/webook/user/_internal/repository/dao"
	service2 "start/webook/user/_internal/service"
)

var user = wire.NewSet(
	web.NewUserHandle,
	service2.NewUserServiceImpl,
	repository2.NewUserRepo,
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
		ioc2.InitWebServer,
		ioc2.InitGorm,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
