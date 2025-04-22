package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"web_golang/goods-web/global"
)

func RemoveTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

//将grpc的code转换成http的状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				//5->404 请求失败
				c.JSON(http.StatusNotFound, gin.H{
					"msg":     e.Message(),
					"message": e.Message(),
				})
			case codes.Internal:
				//13->500 内部错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":     "内部错误",
					"message": e.Message(),
				})
			case codes.InvalidArgument:
				//3->400 参数错误
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":     "参数错误",
					"message": e.Message(),
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":     "商品服务不可用",
					"message": e.Message(),
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":     e.Code(),
					"message": e.Message(),
				})
			}
			return
		}
	}
}

//表单err中文翻译器
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": errs.Translate(global.Trans),
	})
}
