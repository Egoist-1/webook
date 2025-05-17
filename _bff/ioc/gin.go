package ioc

import (
	"github.com/gin-gonic/gin"
	"webook/_bff/web"
	"webook/pkg/ginx/cors"
	"webook/pkg/ginx/jwtx"
)

func InitWebServer(userhandle *web.UserHandle) *gin.Engine {
	s := gin.Default()
	s.Use(GinMiddlewares()...)
	userhandle.RegisterRouter(s)
	//artHandle.RegisterRouter(s)
	return s
}

func GinMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.CorsHandle(),
		jwtx.NewLoginJwtMiddleware().
			IgnorePath("/user/signup").
			IgnorePath("/user/login").
			IgnorePath("/user/send_sms").
			IgnorePath("/user/login_sms").
			Build(),
	}
}
