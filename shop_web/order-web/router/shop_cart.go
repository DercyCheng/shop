package router

import (
	"github.com/gin-gonic/gin"
	"web_golang/order-web/api/shop_cart"
	"web_golang/order-web/middlewares"
)

func InitShopCartRouter(Router *gin.Engine) {
	GoodsRouter := Router.Group("shopcart").Use(middlewares.JWTAuth())
	{
		GoodsRouter.GET("", shop_cart.List)          //购物车列表
		GoodsRouter.DELETE("/:id", shop_cart.Delete) //删除条目
		GoodsRouter.POST("/", shop_cart.New)         //添加商品购物车
		GoodsRouter.PATCH("/:id", shop_cart.Update)  //修改商品购物车
	}
}
