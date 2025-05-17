//go:build wireinject

package main

import (
	"github.com/google/wire"
	"webook/_basicServer/sms/_internal/service/sms/memory"
	ioc2 "webook/_bff/ioc"
	"webook/_bff/web"
	repository2 "webook/_internal/user/_internal/repository"
	"webook/_internal/user/_internal/repository/cache"
	"webook/_internal/user/_internal/repository/dao"
	service2 "webook/_internal/user/_internal/service"
)

var email = wire.NewSet()
var user = wire.NewSet(
	web.NewUserHandle,
	service2.NewUserServiceImpl,
	repository2.NewUserRepo,
	dao.NewUserDao,
)
var code = wire.NewSet(
	service2.NewCodeService,
	repository2.NewCodeRepo,
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
		email,
		ioc2.InitWebServer,
		ioc2.InitGorm,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
