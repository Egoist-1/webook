package jwtx

import "github.com/golang-jwt/jwt/v5"

const JWTKey = "95osj3fUD7fo0mlYdDbncXz4VD2igvf0"

type UserClaims struct {
	Id int64
	jwt.RegisteredClaims
}
