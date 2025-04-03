package jwtx

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type LoginJwtMiddleware struct {
	paths []string
}

func (m *LoginJwtMiddleware) IgnorePath(path string) *LoginJwtMiddleware {
	m.paths = append(m.paths, path)
	return m
}

func NewLoginJwtMiddleware() *LoginJwtMiddleware {
	return &LoginJwtMiddleware{}
}

func (m *LoginJwtMiddleware) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//是否需要 jwt
		for _, path := range m.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		auth := ctx.GetHeader("Authorization")
		if auth == "" {
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}
		split := strings.Split(auth, "")
		j := split[1]
		uc := &UserClaims{}
		token, err := jwt.ParseWithClaims(j, uc, func(token *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}
		if token == nil || uc.Id == 0 || token.Valid == false {
			ctx.AbortWithStatus(http.StatusNonAuthoritativeInfo)
			return
		}
		ctx.Set("claims", uc)
	}
}
