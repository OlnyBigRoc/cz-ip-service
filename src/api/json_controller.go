package api

import (
	"cz-ip-service/src/service"
	"cz-ip-service/src/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewSearchController(searchService *service.SearchService) *SearchController {
	return &SearchController{
		SearchService: searchService,
	}
}

type SearchController struct {
	SearchService *service.SearchService
}

func InitApiJson(group *gin.RouterGroup, searchService *service.SearchService) {
	controller := NewSearchController(searchService)
	{ // {} 是一个代码块，用于限定变量的作用域
		group.GET("", controller.Search)            // 单个IP查询
		group.POST("batch", controller.BatchSearch) // 批量IP查询
	}
}

func (c *SearchController) Search(ctx *gin.Context) {
	res := vo.Result[*vo.IPInfo]{}

	ip := ctx.Query("ip")
	ipInfo, err := c.SearchService.Search(ctx, ip)
	if err != nil {
		ctx.JSON(500, res.Error(err))
		return
	}
	ctx.JSON(200, res.Success(ipInfo))
}

func (c *SearchController) BatchSearch(ctx *gin.Context) {
	res := vo.Result[[]*vo.IPInfo]{}
	data := make([]*vo.IPInfo, 0)
	req := vo.Reqs{}
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, res.Error(err))
		return
	}
	for _, ip := range req.IPs {
		ipInfo, err := c.SearchService.Search(ctx, ip)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, res.Error(err))
		}
		data = append(data, ipInfo)
	}
	res.Success(data)
	ctx.JSON(http.StatusOK, res)
}
