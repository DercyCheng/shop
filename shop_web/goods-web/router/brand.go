package router

import (
	"github.com/gin-gonic/gin"
	"web_golang/goods-web/api/brands"
	"web_golang/goods-web/middlewares"
)

//1. 商品的api接口开发完成
//2. 图片的坑
func InitBrandRouter(Router *gin.Engine) {
	BrandRouter := Router.Group("brand").Use(middlewares.Trace())
	{
		BrandRouter.GET("", brands.BrandList)          // 品牌列表页
		BrandRouter.DELETE("/:id", brands.DeleteBrand) // 删除品牌
		BrandRouter.POST("", brands.NewBrand)          //新建品牌
		BrandRouter.PUT("/:id", brands.UpdateBrand)    //修改品牌信息
	}

	CategoryBrandRouter := Router.Group("categorybrand").Use(middlewares.Trace())
	{
		CategoryBrandRouter.GET("", brands.CategoryBrandList)          // 类别品牌列表页
		CategoryBrandRouter.DELETE("/:id", brands.DeleteCategoryBrand) // 删除类别品牌
		CategoryBrandRouter.POST("", brands.NewCategoryBrand)          //新建类别品牌
		CategoryBrandRouter.PUT("/:id", brands.UpdateCategoryBrand)    //修改类别品牌
		CategoryBrandRouter.GET("/:id", brands.GetCategoryBrandList)   //获取分类的品牌
	}
}
