package controller

import (
	"github.com/gin-gonic/gin"
	"mysite/http/middleware"
	"mysite/model"
	"net/http"
	"net/url"
	"mysite/service"
	"mysite/logger"
	"github.com/Sirupsen/logrus"
	chttp "mysite/http"
	"github.com/gorilla/sessions"
)

type AccountController struct {}

func (self AccountController) RegisterRouter(g *gin.RouterGroup)  {
	g.Any("/account/register", self.Register)
	g.POST("/account/send_activate_email", self.SendActivateEmail)
	g.GET("/account/activate", self.Activate)
	g.Any("/account/login", self.Login)
	g.Any("/account/edit", self.Edit, middleware.NeedLogin())
	g.POST("/account/change_avatar", self.ChangeAvatar, middleware.NeedLogin())
	g.POST("/account/changepwd", self.ChangePwd, middleware.NeedLogin())
	g.Any("/account/forgetpwd", self.ForgetPasswd)
	g.Any("/account/resetpwd", self.ResetPasswd)
	g.GET("/account/logout", self.Logout, middleware.NeedLogin())
}

func (self AccountController) Register(ctx *gin.Context) {
	logger := logger.GetLogger()
	if user, ok := ctx.Get("user"); ok {
		if _, ok := user.(*model.Me); ok {
			ctx.Redirect(http.StatusSeeOther, "/")
		}
	}

	userName := ctx.PostForm("username")
	if userName == "" || ctx.Request.Method != "POST" {
		logger.WithFields(logrus.Fields{"userName":userName}).Info("process fail")
		fail(ctx, 123, "fail")
		return
	}

	if ctx.PostForm("passwd") != ctx.PostForm("passwd2") {
		fail(ctx, 234, "pass1 not equal to pass2")
		return
	}

	//fields := []string{"username", "email", "passwd"}
	fields := []string{"username", "email"}
	form := url.Values{}
	for _, field := range fields {
		form.Set(field, ctx.PostForm(field))
	}

	//调用服务层的逻辑创建用户信息
	errMsg, err := service.DefaultUserService.RegisterUser(ctx, form)
	if err != nil {
		fail(ctx, 123, errMsg)
		return
	}

	logger.Info("process successful")
	success(ctx, gin.H{"helloworld": "hellowrold"})
}

func (self AccountController) SendActivateEmail(ctx *gin.Context) {
}

func (self AccountController) Activate(ctx *gin.Context) {
}

func (self AccountController) Login(ctx *gin.Context) {
	logger := logger.GetLogger()
	if user, ok := ctx.Get("user"); ok {
		if _, ok := user.(*model.Me); ok {
			ctx.Redirect(http.StatusSeeOther, "/")
		}
	}

	userName := ctx.PostForm("username")
	//userName := ctx.Query("username")
	if userName == "" || ctx.Request.Method != "POST" {
		logger.WithFields(logrus.Fields{"userName":userName}).Info("process fail")
		fail(ctx, 123, "fail")
		return
	}

	userLogin, err := service.DefaultUserService.Login(ctx, userName, ctx.PostForm("passwd"))
	if err != nil {
		logger.WithFields(logrus.Fields{"userLogin": userLogin, "err": err}).Info("登录失败")
		fail(ctx, 123, "登录失败")
		return
	}

	//保存cookie信息
	chttp.SetCookie(ctx, userLogin.Username)

	logger.WithFields(logrus.Fields{"userId": userLogin.Uid}).Info("登录成功")
	success(ctx, userLogin.Uid)
	return
}

func (self AccountController) Edit(ctx *gin.Context) {
}

func (self AccountController) ChangeAvatar(ctx *gin.Context) {
}

func (self AccountController) ChangePwd(ctx *gin.Context) {
}

func (self AccountController) ForgetPasswd(ctx *gin.Context) {
}

func (self AccountController) ResetPasswd(ctx *gin.Context) {
}

func (self AccountController) Logout(ctx *gin.Context) {
	logger := logger.GetLogger()
	session := chttp.GetCookieSession(ctx)
	session.Options = &sessions.Options{
		Path: "/",
		MaxAge: -1,
	}
	session.Save(chttp.Request(ctx), chttp.ResponseWriter(ctx))
	logger.Info("退出登录成功")
	success(ctx, "success")
	return
}
