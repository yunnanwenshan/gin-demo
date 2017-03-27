package middleware

import (
	"github.com/gin-gonic/gin"
	"mysite/http"
	"github.com/gorilla/context"
	"mysite/db"
	"mysite/service"
	"mysite/logger"
	"github.com/Sirupsen/logrus"
	"fmt"
)

func NeedLogin() gin.HandlerFunc {
	return needLogin();
}

func needLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer context.Clear(http.Request(ctx))
		logger := logger.GetLogger()
		session := http.GetCookieSession(ctx)
		fmt.Println(session.Values)
		userName, ok := session.Values["username"]
		logger.WithFields(logrus.Fields{"userName": userName}).Info("登录验证")
		if ok {
			if db.MasterDB != nil {
				user := service.DefaultUserService.FindCurrentUser(ctx, userName)
				if user.Uid != 0 {
					ctx.Set("user", user)
				}
			}
		}
		ctx.Next()
	}
}
