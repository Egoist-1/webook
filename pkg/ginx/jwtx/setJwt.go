package jwtx

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SetJwt(ctx *gin.Context, id int64) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Id:               id,
		RegisteredClaims: jwt.RegisteredClaims{
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	})
	j, err := token.SignedString([]byte(JWTKey))
	if err == nil {
		ctx.Header("jwt-token", j)
	}
	return err
}
