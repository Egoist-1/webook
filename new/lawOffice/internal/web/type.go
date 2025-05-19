package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lawOffice/internal/e"
	"net/http"
)

type handler interface {
	RegisterRoutes(server *gin.Engine)
}

type Response struct {
	Code int
	Msg  string
	Data any
}

func handleErr(ctx *gin.Context, data any, err error) bool {
	ok := false
	switch err.(type) {
	case e.Err:
		er := err.(e.Err)
		zap.L().Error(er.Error())
		ctx.JSON(http.StatusOK, Response{
			Code: er.Code().ToInt(),
			Msg:  er.Code().String(),
		})
	case error:
		zap.L().Error(ctx.Request.URL.String(), zap.Error(err))
		ctx.JSON(http.StatusOK, Response{
			Code: e.ServerErr.ToInt(),
			Msg:  e.ServerErr.String() + err.Error(),
		})
	case nil:
		ok = true
		ctx.JSON(http.StatusOK, Response{
			Code: e.Success.ToInt(),
			Data: data,
		})
	}

	return ok
}
