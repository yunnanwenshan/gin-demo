package controller

import (
	"github.com/gin-gonic/gin"
	"mysite/logger"
	"github.com/Sirupsen/logrus"
)

type UserController struct {}

func (self UserController) RegisterRouter(g *gin.RouterGroup)  {
	g.GET("/user/:username", self.Home)
	g.GET("/users", self.ReadList)

}

// Home 用户个人首页
func (self UserController) Home(ctx *gin.Context) {
	logger := logger.GetLogger()
	userName := ctx.Param("username")
	logger.WithFields(logrus.Fields{"userName": userName}).Info("test================i12312312==================")
}

// Home 用户个人首页
func (self UserController) ReadList(ctx *gin.Context) {
}