package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"start/webook/internal/domain"
	"start/webook/internal/service"
	"start/webook/pkg/e"
	"start/webook/pkg/ginx/jwtx"
	"strconv"
)

type ArticleHandle struct {
	svc     service.ArticleService
	intrSvc service.InteractiveService
	biz     string
}

func NewArticleHandle(svc service.ArticleService, intrSvc service.InteractiveService) *ArticleHandle {
	return &ArticleHandle{
		svc:     svc,
		intrSvc: intrSvc,
		biz:     "article",
	}
}

func (h ArticleHandle) RegisterRouter(server *gin.Engine) {
	g := server.Group("/article")
	g.POST("/edit", h.edit)
	g.POST("/withdraw")
	g.POST("/publish", h.publish)
	g.POST("list", h.list)
	g.POST("/detail/id", h.detail)
	pub := g.Group("/pub")
	pub.POST("/detail/:id", h.pubDetail)
	pub.POST("/list", h.pubList)
	pub.POST("/like/:id", h.like)
	//点赞接口

}

func (h *ArticleHandle) edit(ctx *gin.Context) {
	type Req struct {
		Id      int
		Title   string
		Content string
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		zap.L().Warn("articleHandle edit 参数绑定失败", zap.Error(err), zap.Any("req", ctx.Request.Body))
	}
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	uc := claims.(*jwtx.UserClaims)
	art := domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Id,
		},
	}

	aid, err := h.svc.Save(ctx, art)
	DecideErr(ctx, "", aid, err)
}

func (h *ArticleHandle) publish(ctx *gin.Context) {
	type Req struct {
		Id      int
		Title   string
		Content string
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		zap.L().Warn("articleHandle edit 参数绑定失败", zap.Error(err), zap.Any("req", ctx.Request.Body))
	}
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	uc := claims.(*jwtx.UserClaims)
	art := domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Id,
		},
	}
	aid, err := h.svc.Publish(ctx, art)
	DecideErr(ctx, "", aid, err)
}

func (h *ArticleHandle) list(ctx *gin.Context) {
	type Req struct {
		Limit  int
		Offset int
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	claims, ok := ctx.Get("claims")
	if !ok {
		return
	}
	uc := claims.(*jwtx.UserClaims)
	ctx.Writer.Size()
	list, err := h.svc.List(ctx, uc.Id, req.Limit, req.Offset)
	switch err.(type) {
	case e.Err:
		er := err.(e.Err)
		ctx.JSON(http.StatusOK, Result{
			Msg:  er.Code().String(),
			Code: er.Code().ToInt(),
		})
	case error:
		ctx.JSON(http.StatusOK, ServerErr())
	}
	ctx.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Data: list,
	})
}

func (h *ArticleHandle) detail(ctx *gin.Context) {
	param := ctx.Param("id")
	aid, err := strconv.Atoi(param)
	if err != nil {
		return
	}
	claims, err := h.claims(ctx)
	if err != nil {
		return
	}
	art, err := h.svc.Detail(ctx, claims.Id, aid)
	switch err.(type) {
	case e.Err:
		e := err.(e.Err)
		ctx.JSON(http.StatusOK, Result{
			Msg:  e.Code().String(),
			Code: e.Code().ToInt(),
		})
	case error:
		ctx.JSON(http.StatusOK, ServerErr())
	}
	ctx.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Data: art,
	})
}
func (h *ArticleHandle) claims(ctx *gin.Context) (*jwtx.UserClaims, error) {
	claims, ok := ctx.Get("claims")
	if !ok {
		return nil, errors.New("failed to get claims")
	}
	uc, ok := claims.(*jwtx.UserClaims)
	return uc, nil
}

func (h *ArticleHandle) pubList(ctx *gin.Context) {
	type Req struct {
		Limit  int
		Offset int
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	claims, ok := ctx.Get("claims")
	if !ok {
		return
	}
	uc := claims.(*jwtx.UserClaims)
	list, err := h.svc.List(ctx, uc.Id, req.Limit, req.Offset)
	switch err.(type) {
	case e.Err:
		e := err.(e.Err)
		ctx.JSON(http.StatusOK, Result{
			Msg:  e.Code().String(),
			Code: e.Code().ToInt(),
		})
	case error:
		ctx.JSON(http.StatusOK, ServerErr())
	}
	ctx.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Data: list,
	})
}

func (h *ArticleHandle) pubDetail(ctx *gin.Context) {
	param := ctx.Param("id")
	aid, err := strconv.Atoi(param)
	if err != nil {
		return
	}
	claims, err := h.claims(ctx)
	if err != nil {
		return
	}
	var (
		eg   errgroup.Group
		art  domain.Article
		intr domain.Interactive
	)
	//文章详情
	eg.Go(func() error {
		art, err = h.svc.PubDetail(ctx, claims.Id, aid)
		return err
	})
	//文章交互
	eg.Go(func() error {
		intr, err = h.intrSvc.GetIntr(ctx, h.biz, int64(aid), claims.Id)
		return err
	})
	err = eg.Wait()
	switch err.(type) {
	case e.Err:
		e := err.(e.Err)
		ctx.JSON(http.StatusOK, Result{
			Msg:  e.Code().String(),
			Code: e.Code().ToInt(),
		})
		return
	case error:
		ctx.JSON(http.StatusOK, ServerErr())
		return
	}

	var vo PubDetailVo = PubDetailVo{
		Id:         aid,
		Title:      art.Title,
		Content:    art.Content,
		ReadCnt:    intr.ReadCnt,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		Liked:      intr.Liked,
		Collected:  intr.Collected,
		Ctime:      art.Ctime,
	}
	ctx.JSON(http.StatusOK, Result{
		Code: http.StatusOK,
		Data: vo,
	})
}

func (h *ArticleHandle) like(ctx *gin.Context) {
	param := ctx.Param("id")
	_, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusBadRequest,
			Msg:  "invalid article id",
		})
		return
	}
	_, err = h.claims(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusUnauthorized,
			Msg:  "unauthorized",
		})
		return
	}
	//err = h.intrSvc.Like(ctx, h.biz, int64(aid), claims.Id)
	switch err.(type) {
	case e.Err:
		e := err.(e.Err)
		ctx.JSON(http.StatusOK, Result{
			Msg:  e.Code().String(),
			Code: e.Code().ToInt(),
		})
	case error:
		ctx.JSON(http.StatusOK, ServerErr())
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusOK,
			Msg:  "like success",
		})
	}
}
