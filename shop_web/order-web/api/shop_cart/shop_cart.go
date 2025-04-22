package shop_cart

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"web_golang/order-web/api"
	"web_golang/order-web/forms"
	"web_golang/order-web/global"
	"web_golang/order-web/proto"
)

//获取购物车商品
func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("[shop_cat] 查询【购物车】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ids := make([]int32, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)

	}
	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}
	//请求商品服务获取商品信息
	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("[List] 批量查询【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	reMap := gin.H{
		"total": rsp.Totle,
	}
	/*
		{
			"total":12,
			"data":[
				{
					"id":1,
					"goods_id":421,
					"goods_name":421,
					"goods_price":421,
					"goods_image":421,
					"nums":421,
					"checked":421,
				}
			]
		}
	*/
	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		for _, good := range goodsRsp.Data {
			if good.Id == item.GoodsId {
				tempMap := map[string]interface{}{}
				tempMap["id"] = item.Id
				tempMap["good_id"] = item.GoodsId
				tempMap["good_name"] = good.Name
				tempMap["good_image"] = good.GoodsFrontImage
				tempMap["good_price"] = good.ShopPrice
				tempMap["nums"] = item.Nums
				tempMap["checked"] = item.Checked
				goodsList = append(goodsList, tempMap)
			}
		}
	}
	reMap["data"] = goodsList
	ctx.JSON(http.StatusOK, reMap)

}

//添加商品到购物车
func New(ctx *gin.Context) {
	itemForm := forms.ShopCartItemForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	//为了严谨性，添加商品到购物车之前，记得检查一次商品是否存在
	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[goods] 查询【商品信息】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	invRsp, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[inv] 查询【库存信息】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if invRsp.Num < itemForm.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"nums1": "库存不足",
		})
		return
	}
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  int32(userId.(uint)),
		Nums:    itemForm.Nums,
	})
	if err != nil {
		zap.S().Errorw("[orderCreate] 添加【购物车】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}
	itemForm := forms.ShopCartItemUpdateForm{}
	if err = ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	userId, _ := ctx.Get("userId")

	request := proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
		Nums:    itemForm.Nums,
		//Checked: false,
	}
	if itemForm.Checked != nil {
		request.Checked = *itemForm.Checked
	}
	mutex := global.RedisRs.NewMutex(fmt.Sprintf("shop_%d_%d", int32(userId.(uint)), itemForm.Nums))
	if err = mutex.Lock(); err != nil {
		ctx.Status(http.StatusInternalServerError)
	}

	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), &request)
	if ok, err := mutex.Unlock(); !ok || err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
	if err != nil {
		zap.S().Errorw("[OrderUpdate] 更新【购物车记录】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}
	userId, _ := ctx.Get("userId")
	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("[OrderDelete] 删除【购物车记录】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
