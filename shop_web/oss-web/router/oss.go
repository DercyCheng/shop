package router

import (
	"github.com/gin-gonic/gin"
	"web_golang/oss-web/handler"
)

func InitOssRouter(Router *gin.Engine) {
	OssRouter := Router.Group("oss")
	{
		//OssRouter.GET("token", middlewares.JWTAuth(), middlewares.IsAdminAuth(), handler.Token)
		OssRouter.GET("token", handler.Token)
		OssRouter.POST("/callback", handler.HandlerRequest)
	}
}
