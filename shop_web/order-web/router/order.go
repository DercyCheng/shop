package router

import (
	"github.com/gin-gonic/gin"
	"web_golang/order-web/api/order"
	"web_golang/order-web/api/pay"
	"web_golang/order-web/middlewares"
)

func InitOrderRouter(Router *gin.Engine) {
	//BannerRouter := Router.Group("banners").Use(middlewares.Trace())
	OrderRouter := Router.Group("order").Use(middlewares.JWTAuth()).Use(middlewares.Trace())
	{
		OrderRouter.GET("", order.List)       // 订单列表
		OrderRouter.POST("", order.New)       //新建订单
		OrderRouter.GET("/:id", order.Detail) //订单详情
	}
	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("/alipay/notify", pay.Notify) //支付回调
	}
}
