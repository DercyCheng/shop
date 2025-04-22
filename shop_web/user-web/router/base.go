package router

import (
	"github.com/gin-gonic/gin"
	"web_golang/user-web/api"
)

func InitBaseRouter(Router *gin.Engine) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
		BaseRouter.POST("send_sms", api.SendSms)
	}
}
