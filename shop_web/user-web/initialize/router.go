package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_golang/user-web/middlewares"
	router2 "web_golang/user-web/router"
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
	//ApiGroup := Router.Group("/u/v1")
	router2.InitUserRouter(Router)
	router2.InitBaseRouter(Router)
	return Router
}
