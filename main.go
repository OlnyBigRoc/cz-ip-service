package main

import (
	"context"
	"cz-ip-service/src/api"
	"cz-ip-service/src/service"
	"cz-ip-service/src/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Router *gin.Engine
var SearchService *service.SearchService

func main() {
	utils.InitLogger()
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	if Router == nil {
		Router = r
	}
	if SearchService == nil {
		SearchService = service.NewSearchService()
	}
	host := "0.0.0.0"
	port := "80"
	if os.Getenv("HOST") != "" {
		host = os.Getenv("HOST")
	}
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
	}
	Router.LoadHTMLGlob("templates/*")
	Router.Static("/static", "./static")
	// 注册api
	api.InitApiHome(Router.Group("/"))
	api.InitApiJson(Router.Group("/json"), SearchService)
	api.InitApiMsgPack(Router.Group("/msgpack"), SearchService)

	// 启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			utils.Log.Infof("listen: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		utils.Log.Errorf("Server forced to shutdown:%v", err)
	}
	utils.Log.Info("Server exiting!")
}
