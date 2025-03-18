package api

import (
	"cz-ip-service/src/service"
	"cz-ip-service/src/vo"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func NewMsgPackController(searchService *service.SearchService) *MsgPackController {
	return &MsgPackController{
		SearchService: searchService,
	}
}

type MsgPackController struct {
	SearchService *service.SearchService
}

func InitApiMsgPack(group *gin.RouterGroup, searchService *service.SearchService) {
	controller := NewMsgPackController(searchService)
	{ // {} 是一个代码块，用于限定变量的作用域
		group.GET("", controller.Search)            // 单个IP查询
		group.POST("batch", controller.BatchSearch) // 批量IP查询
	}
}

func (c *MsgPackController) Search(ctx *gin.Context) {
	ctx.Header("Content-Type", binding.MIMEMSGPACK2)
	res := vo.Result[*vo.IPInfo]{}

	ip := ctx.Query("ip")
	ipInfo, err := c.SearchService.Search(ctx, ip)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, binding.MIMEMSGPACK2, res.ErrorMsgpack(err))
		return
	}
	ctx.Data(http.StatusOK, binding.MIMEMSGPACK2, res.SuccessMsgpack(ipInfo))
}

func (c *MsgPackController) BatchSearch(ctx *gin.Context) {
	ctx.Header("Content-Type", binding.MIMEMSGPACK2)
	res := vo.Result[[]*vo.IPInfo]{}
	data := make([]*vo.IPInfo, 0)
	req := vo.Reqs{}
	if err := ctx.Bind(&req); err != nil {
		ctx.Data(http.StatusInternalServerError, binding.MIMEMSGPACK2, res.ErrorMsgpack(err))
		return
	}
	for _, ip := range req.IPs {
		ipInfo, err := c.SearchService.Search(ctx, ip)
		if err != nil {
			ctx.Data(http.StatusInternalServerError, binding.MIMEMSGPACK2, res.ErrorMsgpack(err))
		}
		data = append(data, ipInfo)
	}

	ctx.Data(http.StatusOK, binding.MIMEMSGPACK2, res.SuccessMsgpack(data))
}
