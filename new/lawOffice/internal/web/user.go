package web

import (
	"github.com/gin-gonic/gin"
)

var _ handler = &UserHandler{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

type UserHandler struct {
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/user")
	g.GET("/click", h.hello)
}

func (h *UserHandler) hello(ctx *gin.Context) {

}
