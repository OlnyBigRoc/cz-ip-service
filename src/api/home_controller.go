package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HomeController struct {
}

func InitApiHome(group *gin.RouterGroup) {
	controller := &HomeController{}
	{ // {} 是一个代码块，用于限定变量的作用域
		group.GET("", controller.Home)
	}
}

func (c *HomeController) Home(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}
