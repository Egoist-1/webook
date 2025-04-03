package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"start/webook/internal/e"
)

type handle interface {
	RegisterRouter(server *gin.Engine)
}

type Result struct {
	Msg  string
	Code int
	Data any
}

func ServerErr() Result {
	return Result{
		Msg:  e.ServerErr.String(),
		Code: e.ServerErr.ToInt(),
	}
}

// DecideErr msg:成功时返回的信息, data:返回前端的数据, err: error
func DecideErr(ctx *gin.Context, msg string, data any, err error) bool {
	ok := false
	switch err.(type) {
	case e.Err:
		ecode := err.(e.Err)
		zap.L().Warn(ecode.Log())
		ctx.JSON(http.StatusOK, Result{
			Msg:  ecode.Code().String(),
			Code: ecode.Code().ToInt(),
		})
	case error:
		zap.L().Error(ctx.Request.URL.Path, zap.Error(err))
		ctx.JSON(http.StatusOK, ServerErr())
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg:  msg,
			Code: 200,
			Data: data,
		})
		ok = true
		return ok
	}
	return ok
}
