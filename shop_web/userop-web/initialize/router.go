package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_golang/userop-web/middlewares"
	router2 "web_golang/userop-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	//配置跨域
	Router.Use(middlewares.Cors())
	//ApiGroup := Router.Group("/up/v1")
	router2.InitUserFavRouter(Router)
	router2.InitMessageRouter(Router)
	router2.InitAddressRouter(Router)
	return Router
}
