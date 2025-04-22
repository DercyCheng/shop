package router

import (
	"github.com/gin-gonic/gin"
	"web_golang/userop-web/api/message"
	"web_golang/userop-web/middlewares"
)

func InitMessageRouter(Router *gin.Engine) {
	MessageRouter := Router.Group("message").Use(middlewares.JWTAuth())
	{
		MessageRouter.GET("", message.List) // 轮播图列表页
		MessageRouter.POST("", message.New) //新建轮播图
	}
}
