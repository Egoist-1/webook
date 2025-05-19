package web

import "github.com/gin-gonic/gin"

type HomePage struct {
}

func (h *HomePage) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/user")
	g.GET("/click", h.get)
}

func (h *HomePage) get(ctx *gin.Context) {

}
