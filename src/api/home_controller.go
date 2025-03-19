package api

import (
	"cz-ip-service/src/service"
	"cz-ip-service/src/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
	SearchService *service.SearchService
}

func InitApiHome(group *gin.RouterGroup, searchService *service.SearchService) {
	controller := &HomeController{
		SearchService: searchService,
	}
	{ // {} 是一个代码块，用于限定变量的作用域
		group.GET("", controller.Home)
	}
}

func (c *HomeController) Home(ctx *gin.Context) {
	ip := ctx.ClientIP()
	ipInfo, err := c.SearchService.Search(ctx, ip)
	if err != nil {
		ipInfo = &vo.IPInfo{IP: ip}
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"IPInfo": ipInfo,
	})
}
