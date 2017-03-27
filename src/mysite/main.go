package main

import (
	"github.com/gin-gonic/gin"
	"github.com/polaris1119/config"
	_ "mysite/http"
	"mysite/http/routes"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//静态文件
	serveStatic(r)

	//加载模板文件
	r.LoadHTMLGlob("template/*")

	//注册路由
	group := r.Group("")
	routes.RegisterRouters(group)

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

type staticRootConf struct {
	root   string
	isFile bool
}

var staticFileMap = map[string]staticRootConf{
	"/static/":     {"/static", false},
	"/favicon.ico": {"/static/img/go.ico", true},
}

var filterPrefixs = make([]string, 0, 3)

func serveStatic(e *gin.Engine) {
	for prefix, rootConf := range staticFileMap {
		filterPrefixs = append(filterPrefixs, prefix)

		if rootConf.isFile {
			e.StaticFile(prefix, config.ROOT+rootConf.root)
		} else {
			e.Static(prefix, config.ROOT+rootConf.root)
		}
	}
}
