package initialize

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	
	"web_api/user_web/middlewares"
	"web_api/user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	
	// 添加全局中间件
	Router.Use(middlewares.Cors())                     // 跨域中间件
	Router.Use(middlewares.RequestMetricsMiddleware()) // 请求监控
	Router.Use(middlewares.SecurityHeadersMiddleware()) // 安全头
	Router.Use(middlewares.ErrorHandlerMiddleware())   // 统一错误处理
	
	// 健康检查路由
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	
	// 404处理
	Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "API endpoint not found",
		})
	})
	
	// API路由组
	ApiGroup := Router.Group("/u/v1")
	
	// 为敏感API添加限流中间件
	sensitiveGroup := ApiGroup.Group("/auth")
	sensitiveGroup.Use(middlewares.RateLimitMiddleware())
	
	// 注册各种路由
	router.InitUserRouter(ApiGroup)
	//router.InitBaseRouter(ApiGroup) // 测试的时候先过滤掉图片验证码和ali短信验证
	
	zap.S().Info("路由初始化完成")
	return Router
}
