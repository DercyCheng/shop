package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_golang/goods-web/middlewares"
	router2 "web_golang/goods-web/router"
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
	//ApiGroup := Router.Group("/v1")
	router2.InitGoodsRouter(Router)
	router2.InitCategoryRouter(Router)
	router2.InitBannerRouter(Router)
	router2.InitBrandRouter(Router)
	return Router
}
