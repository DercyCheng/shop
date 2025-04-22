package router

import (
	"github.com/gin-gonic/gin"
	"web_golang/goods-web/api/goods"
	"web_golang/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.Engine) {
	GoodsRouter := Router.Group("good").Use(middlewares.Trace())
	{
		GoodsRouter.GET("", goods.List) //商品列表
		//一定要注意middlewares路径
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)               //新建商品
		GoodsRouter.GET("/:id", goods.Detail)                                                           //商品详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)      //删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)                                                    //获取库存
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)         //更新库存
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus) //商品状态

	}
}
