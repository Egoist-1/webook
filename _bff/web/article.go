package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	domain2 "webook/_internal/article/_internal/domain"
	service2 "webook/_internal/article/_internal/service"
	"webook/_internal/interactive/_internal/domain"
	"webook/_internal/interactive/_internal/service"
	"webook/pkg/er"
	"webook/pkg/ginx/jwtx"
)

type ArticleHandle struct {
	svc     service2.ArticleService
	intrSvc service.InteractiveService
	biz     string
}

func NewArticleHandle(svc service2.ArticleService, intrSvc service.InteractiveService) *ArticleHandle {
	return &ArticleHandle{
		svc:     svc,
		intrSvc: intrSvc,
		biz:     "article",
	}
}

func (h *ArticleHandle) RegisterRouter(server *gin.Engine) {
	g := server.Group("/article")
	g.POST("/edit", h.edit)
	g.POST("/withdraw")
	g.POST("/publish", h.publish)
	g.POST("list", h.list)
	g.POST("/detail/id", h.detail)
	pub := g.Group("/pub")
	pub.POST("/unpublish/:id", h.unpublish)
	pub.POST("/detail/:id", h.pubDetail)
	pub.POST("/list", h.pubList)
	pub.POST("/like", h.like)
	pub.POST("/collected", h.collectd)
	pub.POST("/cancelLike/:id", h.cancelLike)
	pub.POST("/cancelCollection", h.cancelCollection)
}

func (h *ArticleHandle) edit(ctx *gin.Context) {
	type Req struct {
		Id      int64
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
	art := domain2.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain2.Author{
			Id: uc.Id,
		},
	}

	aid, err := h.svc.Save(ctx, art)
	HandleErr(ctx, "", aid, err)
}

func (h *ArticleHandle) publish(ctx *gin.Context) {
	type Req struct {
		Id      int64
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
	art := domain2.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain2.Author{
			Id: uc.Id,
		},
	}
	aid, err := h.svc.Publish(ctx, art)
	HandleErr(ctx, "", aid, err)
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
	case er.Err:
		er := err.(er.Err)
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
	art, err := h.svc.Detail(ctx, claims.Id, int64(aid))
	switch err.(type) {
	case er.Err:
		e := err.(er.Err)
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
	list, err := h.svc.PubList(ctx, uc.Id, req.Limit, req.Offset)
	switch err.(type) {
	case er.Err:
		e := err.(er.Err)
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
		art  domain2.Article
		intr domain.Interactive
	)
	//文章详情
	eg.Go(func() error {
		art, err = h.svc.PubDetail(ctx, claims.Id, int64(aid))
		return err
	})
	//文章交互
	eg.Go(func() error {
		intr, err = h.intrSvc.GetIntr(ctx, h.biz, int64(aid), claims.Id)
		return err
	})
	err = eg.Wait()
	switch err.(type) {
	case er.Err:
		e := err.(er.Err)
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
	type Req struct {
		Aid   int  `json:"aid"`
		Liked bool `json:"liked"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusBadRequest,
			Msg:  "invalid article id",
		})
		return
	}
	claims, err := h.claims(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusUnauthorized,
			Msg:  "unauthorized",
		})
		return
	}
	err = h.intrSvc.Liked(ctx, h.biz, int64(req.Aid), int64(claims.Id))
	HandleErr(ctx, "点赞成功", nil, err)
}

func (h *ArticleHandle) collectd(ctx *gin.Context) {
	type Req struct {
		Aid       int64 `json:"aid"`
		Cid       int64 `json:"cid"`
		Collected bool  `json:"collected"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusBadRequest,
			Msg:  "invalid article id",
		})
		return
	}
	claims, err := h.claims(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: http.StatusUnauthorized,
			Msg:  "unauthorized",
		})
		return
	}
	err = h.intrSvc.Collected(ctx, h.biz, req.Aid, claims.Id, req.Cid)
	HandleErr(ctx, "收藏成功", nil, err)
}

func (h *ArticleHandle) unpublish(ctx *gin.Context) {
	param := ctx.Param("id")
	aid, err := strconv.Atoi(param)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		zap.L().Warn("articleHandle edit 参数绑定失败", zap.Error(err), zap.Any("req", ctx.Request.Body))
	}
	err = h.svc.Unpublish(ctx, int64(aid))
	HandleErr(ctx, "", nil, err)
}

func (h *ArticleHandle) cancelLike(ctx *gin.Context) {
	aid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return
	}
	claims, err := h.claims(ctx)
	if err != nil {
		return
	}
	err = h.intrSvc.CancelLike(ctx, h.biz, claims.Id, int64(aid))
	HandleErr(ctx, "", nil, err)
}

func (h *ArticleHandle) cancelCollection(ctx *gin.Context) {
	type Req struct {
		ArticleId    int64
		CollectionId int64
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	claims, err := h.claims(ctx)
	if err != nil {
		return
	}
	err = h.intrSvc.CancelCollection(ctx, h.biz, claims.Id, req.ArticleId, req.CollectionId)
	HandleErr(ctx, "", nil, err)
}
