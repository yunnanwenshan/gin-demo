package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexController struct {}

func (self IndexController) RegisterRouter(g *gin.RouterGroup)  {
	g.Any("/", self.Index)
}

func (self IndexController) Index(ctx *gin.Context)  {
	ctx.HTML(http.StatusOK, "index.html", map[string]interface{}{"title": "1234"})
}
