//go:build wireinject

package main

import (
	"github.com/google/wire"
	"lawOffice/internal/web"
	"lawOffice/ioc"
)

var userServer = wire.NewSet(
	web.NewUserHandler,
)

func InitWebServer() *App {
	wire.Build(
		userServer,
		ioc.InitWebServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
