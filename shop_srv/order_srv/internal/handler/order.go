package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"math/rand"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"order_srv/global"
	"order_srv/model"
	v1goods "order_srv/proto/goods"
	v1inventory "order_srv/proto/inventory"
	proto "order_srv/proto/v1order"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func GenerateOrderSn(userId int32) string {
	//订单号的生成规则
	/*
		年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	//rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}
func (s *OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	//获取用户的购物车列表
	var shopCarts []model.ShoppingCart
	var rsp proto.CartItemListResponse
	if result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	} else {
		rsp.Totle = int32(result.RowsAffected)
	}
	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return &rsp, nil
}
func (s *OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//将商品添加到购物车 1.购物车中原本没有这件商品 - 新建一个记录 2.这个商品之前添加到了购物车 - 合并
	var shopCart model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); result.RowsAffected == 1 {
		//如果记录已经存在，则合并购物车记录 - 更新操作
		shopCart.Nums += req.Nums
	} else {
		//插入操作
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}
	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}
func (s *OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.Empty, error) {
	//更新购物车记录，更新数量和选中状态
	var shopCart model.ShoppingCart
	if result := global.DB.First(&shopCart, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	global.DB.Save(&shopCart)
	return &proto.Empty{}, nil
}
func (s *OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.Empty, error) {
	if result := global.DB.Delete(&model.ShoppingCart{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &proto.Empty{}, nil
}
func (s *OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp proto.OrderListResponse
	//你是后台管理系统查询  还是电商系统中心查询   srv层一般不关心这些   微服务要做到通用 底层尽量不要和业务挂钩  关心的查询条件
	//gorm框架中 如果传递过来是0值  则不会构建sql语句  正好符合查询条件
	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)
	//分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
		})
	}
	return &rsp, nil
}
func (s *OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse
	//这个订单的id是否是当前用户的订单  如果在web层用户传递过来一个id的订单，web层应该先查询一下订单id是否是当前用户的
	//在个人中心可以这样做，如果是后台管理系统，web层如果是后代管理系统 那么值传递order的id，如果是电商系统还需要一个用户的id
	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	orderInfo.Name = order.SignerName
	orderInfo.Mobile = order.SingerMobile

	rsp.OrderInfo = &orderInfo
	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}
	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			Nums:       orderGood.Nums,
		})
	}
	return &rsp, nil
}

type OrderListener struct {
	//Code        codes.Code
	//Detail      string
	//ID          int32
	//OrderAmount float32
	//Ctx         context.Context
}
type OrderListenerInfo struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	Ctx         context.Context
}

var myVar = OrderListenerInfo{}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	//var orderlistenerinfo OrderListenerInfo
	//orderlistenerinfoPrt:=&orderlistenerinfo
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	//var listener OrderListener
	//fmt.Println(listener.Ctx)
	//fmt.Println(o.Ctx)
	parentSpan := opentracing.SpanFromContext(myVar.Ctx)
	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	shopCartSpan := opentracing.GlobalTracer().StartSpan("select_shopcart", opentracing.ChildOf(parentSpan.Context()))
	if result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		myVar.Code = codes.InvalidArgument
		myVar.Detail = "没有选中结算的商品"
		return primitive.RollbackMessageState
		//return nil, status.Errorf(codes.InvalidArgument, "没有选中结算的商品")
	}
	shopCartSpan.Finish()
	for _, shopCarat := range shopCarts {
		goodsIds = append(goodsIds, shopCarat.Goods)
		goodsNumsMap[shopCarat.Goods] = shopCarat.Nums
	}
	//跨服务调用 - 商品微服务
	queryGoodsSpan := opentracing.GlobalTracer().StartSpan("query_goods", opentracing.ChildOf(parentSpan.Context()))
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &v1goods.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		myVar.Code = codes.Internal
		myVar.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
		//return nil, status.Errorf(codes.Internal, "批量查询商品信息失败")
	}
	queryGoodsSpan.Finish()
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*v1inventory.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &v1inventory.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}
	//跨服务调用 - 库存微服务
	queryinvSpan := opentracing.GlobalTracer().StartSpan("query_inv", opentracing.ChildOf(parentSpan.Context()))
	if _, err = global.InventorySrvClient.Sell(context.Background(), &v1inventory.SellInfo{GoodsInfo: goodsInvInfo, OrderSn: orderInfo.OrderSn}); err != nil {
		//如果是因为网络问题状态码变更，这种如何避免误判：改写sell逻辑   判断返回的状态码是否是在sell里的  如果不是就判断   如果是 也要判断
		myVar.Code = codes.ResourceExhausted
		myVar.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
		//return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
	}
	queryinvSpan.Finish()
	/*
		o.Code = codes.Internal
		o.Detail = "【测试】失败-1"
		return primitive.UnknowState
	*/
	//生成订单表
	//20210308xxxx
	tx := global.DB.Begin()
	orderInfo.OrderMount = orderAmount
	saveOrderSpan := opentracing.GlobalTracer().StartSpan("save_order", opentracing.ChildOf(parentSpan.Context()))
	//order := model.OrderInfo{
	//	OrderSn:      GenerateOrderSn(orderInfo.User),//不能重新生成
	//	OrderMount:   orderAmount,
	//	Address:      req.Address,
	//	SignerName:   req.Name,
	//	SingerMobile: req.Mobile,
	//	Post:         req.Post,
	//	User:         req.UserId,
	//}
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()
		myVar.Code = codes.Internal
		myVar.Detail = "创建订单失败"
		return primitive.RollbackMessageState
		//return nil, status.Errorf(codes.Internal, "创建订单失败")
	}
	saveOrderSpan.Finish()

	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}
	//批量插入orderGoods
	saveOrderGoodsSpan := opentracing.GlobalTracer().StartSpan("save_order_goods", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		myVar.Code = codes.Internal
		myVar.Detail = "创建订单失败"
		return primitive.RollbackMessageState
		//return nil, status.Errorf(codes.Internal, "创建订单失败")
	}
	saveOrderGoodsSpan.Finish()

	deleteShopCartSpan := opentracing.GlobalTracer().StartSpan("delete_shopcart", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		myVar.Code = codes.Internal
		myVar.Detail = "删除购物车记录失败"
		return primitive.RollbackMessageState
		//return nil, status.Errorf(codes.Internal, "创建订单失败")
	}
	deleteShopCartSpan.Finish()

	//p, err := rocketmq.NewProducer(
	//	producer.WithNsResolver(primitive.NewPassthroughResolver([]string{fmt.Sprintf("%s:%d", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)})),
	//	producer.WithGroupName("yanshi"),
	//)
	//if err != nil {
	//	panic("生成producer失败")
	//}
	//
	////不要在一个进程中使用多个producer， 但是不要随便调用shutdown因为会影响其他的producer
	//if err = p.Start(); err != nil {
	//	fmt.Println(err.Error())
	//	panic("启动producer失败")
	//}

	msg = primitive.NewMessage("order_timeout", msg.Body)
	msg.WithDelayTimeLevel(3)
	_, err = global.MQSendClient.SendSync(context.Background(), msg)
	if err != nil {
		zap.S().Errorf("发送延时消息失败: %v\n", err)
		tx.Rollback()
		myVar.Code = codes.Internal
		myVar.Detail = "发送延时消息失败"
		return primitive.RollbackMessageState
	}
	//err = p.Shutdown()
	//if err != nil {
	//	panic("Shutdown:")
	//}

	//提交事务
	tx.Commit()
	myVar.OrderAmount = orderAmount
	myVar.ID = orderInfo.ID
	myVar.Code = codes.OK
	return primitive.UnknowState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	if result := global.DB.Where(&model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState
}
func (s *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		新建订单：
			1.从购物车中获取到选中的商品
			2.商品的价格自己查询 - 访问商品服务（跨微服务）
			3.库存的扣减 - 访问库存服务（跨微服务）
			4.订单的基本信息表 - 订单的商品信息表
			5.从购物车中删除已购买的记录
	*/
	//orderListener := OrderListener{Ctx: ctx}
	myVar.Ctx = ctx
	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}
	//应该在消息中具体指明一个订单的具体的商品的扣减情况
	jsonString, _ := json.Marshal(&order)

	_, err := global.MQSendTranClient.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", jsonString))
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
		return nil, status.Error(codes.Internal, "发送消息队列失败")
	}
	if myVar.Code != codes.OK {
		return nil, status.Error(myVar.Code, myVar.Detail)
	}
	//var lin *OrderListener
	//fmt.Println(lin.ID)
	//newOrderListener := &OrderListener{}
	//var orderLInfo model.OrderInfo
	//if result := global.DB.Where(&model.OrderInfo{OrderSn: order.OrderSn}).First(&orderLInfo); result.RowsAffected == 0 {
	//	return nil, status.Error(codes.Internal, "查无此单号")
	//}

	return &proto.OrderInfoResponse{Id: myVar.ID, OrderSn: order.OrderSn, Total: myVar.OrderAmount}, nil
}
func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*proto.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &proto.Empty{}, nil
}
func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var orderInfo model.OrderInfo
		_ = json.Unmarshal(msgs[i].Body, &orderInfo)
		fmt.Printf("获取到订单超时消息：%v\n", time.Now())
		//查询订单的支付状态，如果已支付什么都不做，如果未支付，归还库存
		var order model.OrderInfo
		if result := global.DB.Model(&model.OrderInfo{}).Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&order); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		if order.Status != "TRADE_SUCCESS" {
			tx := global.DB.Begin()

			//p, err := rocketmq.NewProducer(producer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)}))
			//if err != nil {
			//	panic("生成producer失败")
			//}
			//
			//if err = p.Start(); err != nil {
			//	panic("启动producer失败")
			//}
			//归还库存 我们可以模仿order中发送一个消息到order_reback中去
			_, err := global.MQSendClient.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))
			if err != nil {
				tx.Rollback()
				zap.S().Errorf("【超时归还】发送失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}
			//修改订单的状态为已支付
			order.Status = "TRADE_CLOSED"
			if result := tx.Save(&order); result.Error != nil {
				tx.Rollback()
				zap.S().Errorf("【超时归还】修改支付失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}
		}
	}
	return consumer.ConsumeSuccess, nil
}
