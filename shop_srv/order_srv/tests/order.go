package main

import (
	"context"
	"fmt"
	proto "order_srv/proto/v1order"
)

func TestCreateCartItem(userId, nums, goodsId int32) {
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		Nums:    nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}
func TestCartItemList(userId int32) {
	rsp, err := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: userId,
	})
	if err != nil {
		panic(err)
	}
	for _, item := range rsp.Data {
		fmt.Println(item.Id, item.GoodsId, item.Nums)
	}
}
func TestUpdateCartItem(id int32) {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      id,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("购物车更新成功")
}
func TestCreateOrder() {
	_, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  12,
		Address: "寨里河镇",
		Name:    "jzin",
		Mobile:  "15263331478",
		Post:    "发货啊啊啊啊啊啊啊",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("新建成功")
}
func TestGetOrderDetail(orderId int32) {
	rsp, err := orderClient.OrderDetail(context.Background(), &proto.OrderRequest{
		Id: orderId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.OrderInfo.OrderSn)
	for _, good := range rsp.Goods {
		fmt.Println(good.GoodsName)
	}
}
func TestOrderList() {
	rsp, err := orderClient.OrderList(context.Background(), &proto.OrderFilterRequest{
		UserId: 12,
	})
	if err != nil {
		panic(err)
	}
	for _, order := range rsp.Data {
		fmt.Println(order.OrderSn)
	}
}
