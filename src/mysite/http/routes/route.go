package routes

import (
	"github.com/gin-gonic/gin"
	"mysite/http/controller"
)

func RegisterRouters(g *gin.RouterGroup)  {
	new(controller.UserController).RegisterRouter(g)
	new(controller.AccountController).RegisterRouter(g)
	new(controller.IndexController).RegisterRouter(g)
}