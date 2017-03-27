package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// render html 输出
func render(ctx *gin.Context, contentTpl string, data map[string]interface{}) error {
	return nil
	//return Render(ctx, contentTpl, data)
}

func success(ctx *gin.Context, data interface{}) error {
	result := gin.H{
		"ok":   1,
		"msg":  "操作成功",
		"data": data,
	}

	ctx.JSON(http.StatusOK, result)

	return nil
}

func fail(ctx *gin.Context, code int, msg string) error {
	result := map[string]interface{}{
		"ok":    code,
		"error": msg,
	}

	ctx.JSON(http.StatusOK, result)

	return nil
}