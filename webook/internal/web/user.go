package web

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"net/http"
	"start/webook/internal/domain"
	"start/webook/internal/e"
	"start/webook/internal/service"
	"start/webook/pkg/ginx/jwtx"
)

var _ handle = (*UserHandle)(nil)

type UserHandle struct {
	EmailRegex    *regexp2.Regexp
	PasswordRegex *regexp2.Regexp
	PhoneRegex    *regexp2.Regexp
	svc           service.UserService
	codeSvc       service.CodeService
	biz           string
}

func (h UserHandle) RegisterRouter(server *gin.Engine) {
	g := server.Group("/user")
	g.POST("/signup", h.signup)
	g.POST("/login", h.login)
	g.POST("/edit", h.edit)
	g.POST("/profile", h.profile)
	g.POST("/send_sms", h.sendSms)
	g.POST("/login_sms", h.loginSMS)
}

func NewUserHandle(svc service.UserService, codesvc service.CodeService) *UserHandle {
	var EmailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	var PasswordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	var PhoneRegexPattern = `^1[3-9]\d{9}$`
	return &UserHandle{
		EmailRegex:    regexp2.MustCompile(EmailRegexPattern, 0),
		PasswordRegex: regexp2.MustCompile(PasswordRegexPattern, 0),
		PhoneRegex:    regexp2.MustCompile(PhoneRegexPattern, 0),
		svc:           svc,
		codeSvc:       codesvc,
		biz:           "login",
	}
}
func (h UserHandle) profile(ctx *gin.Context) {
	claims, er := ctx.Get("claims")
	fmt.Println(er)
	uc := claims.(*jwtx.UserClaims)
	u, err := h.svc.Profile(ctx, uc.Id)
	var profile ProfileVO
	copier.Copy(&profile, &u)
	DecideErr(ctx, "简介", profile, err)
}

func (h UserHandle) signup(ctx *gin.Context) {
	type Req struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, ServerErr())
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "密码不一致",
			Code: e.UserInvalidInput.ToInt(),
		})
		return
	}
	ok, err := h.EmailRegex.MatchString(req.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, ServerErr())
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "邮箱有误",
			Code: e.UserInvalidInput.ToInt(),
		})
		return
	}
	//ok, err = h.PasswordRegex.MatchString(req.Email)
	//if err != nil {
	//	ctx.AbortWithStatusJSON(http.StatusOK, ServerErr())
	//}
	//if !ok {
	//	ctx.JSON(http.StatusOK, Result{
	//		Msg:  "密码必须包含数字、特殊字符，并且长度不能小于 8 位",
	//		code: e.UserInvalidInput.code(),
	//	})
	//	return
	//}
	uid, err := h.svc.Signup(ctx, domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	ok = DecideErr(ctx, "注册成功", uid, err)
	if !ok {
		return
	}
	err = jwtx.SetJwt(ctx, uid)
	if err != nil {
		zap.L().Error("jwt设置失败", zap.Error(err))
		ctx.JSON(http.StatusOK, ServerErr())
		return
	}
}

func (h UserHandle) login(ctx *gin.Context) {
	type Req struct {
		Email    string
		Password string
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var uid int
	uid, err = h.svc.LoginEmail(ctx, req.Email, req.Password)
	ok := DecideErr(ctx, "登录成功", nil, err)
	if !ok {
		return
	}
	err = jwtx.SetJwt(ctx, uid)
	if err != nil {
		zap.L().Error("jwt设置失败", zap.Error(err))
		ctx.JSON(http.StatusOK, ServerErr())
		return
	}
}

func (h UserHandle) edit(ctx *gin.Context) {
	type Req struct {
		Name    string
		Phone   string
		AboutMe string
		CTime   int64
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	claims, _ := ctx.Get("claims")
	uc := claims.(*jwtx.UserClaims)
	var domainU domain.User
	copier.Copy(&domainU, &req)
	domainU.Id = uc.Id
	err = h.svc.Edit(ctx, domainU)
	DecideErr(ctx, "编辑成功", nil, err)
}

func (h UserHandle) sendSms(ctx *gin.Context) {
	type Req struct {
		Phone string
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	ok, err := h.PhoneRegex.MatchString(req.Phone)
	if err != nil {
		zap.L().Warn("手机正则校验错误", zap.Error(err))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Msg:  e.UserInvalidInput.String(),
			Code: e.UserInvalidInput.ToInt(),
		})
		return
	}
	err = h.codeSvc.SendSMS(ctx, h.biz, req.Phone)
	DecideErr(ctx, "发送成功", nil, err)
}

func (h UserHandle) loginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string
		Code  string
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	ok, err := h.PhoneRegex.MatchString(req.Phone)
	if err != nil {
		zap.L().Warn("手机正则校验错误", zap.Error(err))
		return
	}
	if !ok {
		zap.L().Warn("手机正则校验错误", zap.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Msg:  e.UserInvalidInput.String(),
			Code: e.UserInvalidInput.ToInt(),
		})
		return
	}
	err = h.codeSvc.Verify(ctx, h.biz, req.Phone, req.Code)
	DecideErr(ctx, "", nil, err)
	if err != nil {
		return
	}
	uid, err := h.svc.LoginByPhone(ctx, req.Phone)
	DecideErr(ctx, "登录成功", nil, err)
	if err != nil {
		return
	}
	err = jwtx.SetJwt(ctx, uid)
	if err != nil {
		zap.L().Error("jwt设置失败", zap.Error(err))
		ctx.JSON(http.StatusOK, ServerErr())
		return
	}
}
