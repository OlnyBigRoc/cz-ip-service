package api

import (
	"cz-ip-service/src/constant"
	"cz-ip-service/src/metrics"
	"cz-ip-service/src/service"
	"cz-ip-service/src/utils"
	"cz-ip-service/src/vo"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	startTime := time.Now()
	ctx.Header("Content-Type", binding.MIMEMSGPACK2)
	res := vo.Result[*vo.IPInfo]{}
	success := false
	defer func() {
		metrics.RecordRequest(success)
		metrics.RecordResponseTime(time.Since(startTime))
	}()

	ip := ctx.Query("ip")
	if ip == "" {
		ctx.Data(http.StatusBadRequest, binding.MIMEMSGPACK2, res.ErrorMsgpack(fmt.Errorf("ip parameter is required")))
		return
	}
	if !utils.IsValidIP(ip) {
		ctx.Data(http.StatusBadRequest, binding.MIMEMSGPACK2, res.ErrorMsgpack(fmt.Errorf("invalid ip format: %s", ip)))
		return
	}
	ipInfo, err := c.SearchService.Search(ctx, ip)
	if err != nil {
		ctx.Data(http.StatusInternalServerError, binding.MIMEMSGPACK2, res.ErrorMsgpack(err))
		return
	}
	success = true
	metrics.RecordIPQuery(ip)
	ctx.Data(http.StatusOK, binding.MIMEMSGPACK2, res.SuccessMsgpack(ipInfo))
}

func (c *MsgPackController) BatchSearch(ctx *gin.Context) {
	startTime := time.Now()
	ctx.Header("Content-Type", binding.MIMEMSGPACK2)
	res := vo.Result[[]*vo.IPInfo]{}
	success := false
	defer func() {
		metrics.RecordRequest(success)
		metrics.RecordResponseTime(time.Since(startTime))
	}()

	req := vo.Reqs{}
	if err := ctx.Bind(&req); err != nil {
		ctx.Data(http.StatusBadRequest, binding.MIMEMSGPACK2, res.ErrorMsgpack(err))
		return
	}

	valid, errMsg := utils.ValidateIPList(req.IPs, 100)
	if !valid {
		ctx.Data(http.StatusBadRequest, binding.MIMEMSGPACK2, res.ErrorMsgpack(fmt.Errorf(errMsg)))
		return
	}

	// 使用channel控制并发
	semaphore := make(chan struct{}, constant.MaxConcurrent)
	var wg sync.WaitGroup
	data := make([]*vo.IPInfo, len(req.IPs))
	errChan := make(chan error, len(req.IPs))

	for i, ip := range req.IPs {
		wg.Add(1)
		go func(index int, ipAddr string) {
			defer wg.Done()
			semaphore <- struct{}{}        // 获取信号量
			defer func() { <-semaphore }() // 释放信号量

			ipInfo, err := c.SearchService.Search(ctx, ipAddr)
			if err != nil {
				errChan <- err
				return
			}
			data[index] = ipInfo
			metrics.RecordIPQuery(ipAddr)
		}(i, ip)
	}

	// 等待所有goroutine完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误发生
	for err := range errChan {
		ctx.Data(http.StatusInternalServerError, binding.MIMEMSGPACK2, res.ErrorMsgpack(err))
		return
	}

	success = true
	ctx.Data(http.StatusOK, binding.MIMEMSGPACK2, res.SuccessMsgpack(data))
}
