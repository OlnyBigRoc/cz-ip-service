package api

import (
	"cz-ip-service/src/metrics"
	"cz-ip-service/src/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitApiMetrics(group *gin.RouterGroup) {
	group.GET("", GetMetrics)
}

func GetMetrics(ctx *gin.Context) {
	res := vo.Result[map[string]interface{}]{}
	_metrics := metrics.GetMetrics()
	ctx.JSON(http.StatusOK, res.Success(_metrics))
}
