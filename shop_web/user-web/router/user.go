package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"web_golang/user-web/api"
	"web_golang/user-web/middlewares"
)

func InitUserRouter(Router *gin.Engine) {
	UserRouer := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserRouer.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		//UserRouer.GET("list", api.GetUserList)
		UserRouer.POST("pwd_login", api.PassWorldLogin)
		UserRouer.POST("register", api.Register)

		UserRouer.GET("detail", middlewares.JWTAuth(), api.GetUserDetail)
		UserRouer.PATCH("update", middlewares.JWTAuth(), api.UpdateUser)
	}
	//服务注册和发现
}
