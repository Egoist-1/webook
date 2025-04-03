package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

func CorsHandle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cors.New(cors.Config{
			AllowHeaders:  []string{"Content-Type", "Authorization"},
			ExposeHeaders: []string{"jwt-token"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return true
			},
		})
	}
}
