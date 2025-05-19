package web

import (
	"github.com/gin-gonic/gin"
	"lawOffice/internal/service"
)

type MonitorLink struct {
	svc service.MonitorLink
}

func (h *MonitorLink) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/user")
	g.GET("/click", h.get)
}

func (h *MonitorLink) get(ctx *gin.Context) {
	ctx.Request.URL.Query()
	var str string
	err := h.svc.Analysis(str)
}
