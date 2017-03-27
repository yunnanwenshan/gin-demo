package http

import (
	"github.com/boj/redistore"
	"github.com/polaris1119/config"
	"mysite/logger"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gorilla/sessions"
)

type StoreRedis struct {
	store *redistore.RediStore
}

var Store = new(StoreRedis)

func init()  {
	logger := logger.GetLogger()
	st, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte(config.ConfigFile.MustValue("global", "cookie_secret")))
	if err != nil {
		logger.WithFields(logrus.Fields{"err": err}).Info("NewRediStore初始化失败")
		panic("redis初始化失败")
	}
	logger.Info("redis初始化成功")
	Store.store = st
}

func GetStore() *redistore.RediStore {
	return Store.store
}

func Request(ctx *gin.Context) *http.Request {
	return ctx.Request
}

func ResponseWriter(ctx *gin.Context) http.ResponseWriter  {
	return ctx.Writer
}

func GetCookieSession(ctx *gin.Context) *sessions.Session  {
	session, _ := Store.store.Get(Request(ctx), "user")
	return session
}

func SetCookie(ctx *gin.Context, userName string)  {
	Store.store.Options.HttpOnly = true

	session := GetCookieSession(ctx)
	if ctx.PostForm("remember_me") != "1" {
		session.Options = &sessions.Options{
			Path: "/",
			HttpOnly: true,
		}
	}
	session.Values["username"] = userName
	req := Request(ctx)
	res := ctx.Writer
	session.Save(req, res)
}

const (
	LayoutTpl      = "common/layout.html"
	AdminLayoutTpl = "common.html"
)

//func Render(ctx *gin.Context, contentTpl string, data map[string]interface{}) error {
//	if data == nil {
//		data = map[string]interface{}{}
//	}
//
//	logger := logger.GetLogger()
//	contentTpl = LayoutTpl + "," + contentTpl
//	htmlFiles := strings.Split(contentTpl, ",")
//	for i, contentTpl := range htmlFiles {
//		htmlFiles[i] = config.TemplateDir + contentTpl
//	}
//
//	tpl, err := template.New("layout.html").ParseFiles(htmlFiles...)
//	if err != nil {
//		logger.Errorf("解析模板出错（ParseFiles）：[%q] %s\n", Request(ctx).RequestURI, err)
//		return err
//	}
//
//}

