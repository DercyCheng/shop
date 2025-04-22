package main

import (
	"context"
	"fmt"
	"sync"
)

// 获取商品详情页
func TestSetInv(goodsId, num int32) {
	_, err := brandClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
		Num:     num,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")
}
func TestInvDetail(goodsId int32) {
	rsp, err := brandClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Num)

}
func TestSell(wg *sync.WaitGroup) {
	/*
		1.第一件扣减成功 第二件：1.没有库存信息 2.库存不足
		2.两件都扣减成功
	*/
	defer wg.Done()
	_, err := brandClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 2},
			{GoodsId: 422, Num: 2},
			{GoodsId: 423, Num: 2},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存扣减成功")
}

func TestReback() {
	_, err := brandClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 10},
			{GoodsId: 422, Num: 30},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存归还成功")
}
