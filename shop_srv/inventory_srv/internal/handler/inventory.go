package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"sync"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"inventory_srv/global"
	"inventory_srv/model"
	proto "inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

var _ proto.InventoryServer = &InventoryServer{}

func (s *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*proto.Empty, error) {
	//设置库存，如果我要更新库存
	var inv model.Inventory
	//只有是主键的情况才能直接用id   goodsid不是主键
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	if inv.Goods == 0 {
		inv.Goods = req.GoodsId
	}
	inv.Stocks = req.Num
	if result := global.DB.Save(&inv); result.Error != nil {
		return nil, result.Error
	}
	return &proto.Empty{}, nil

}

// 获取库存详情
func (s *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "库存信息不存在")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// var m sync.Mutex
var m2 sync.Mutex

func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*proto.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1

	tx := global.DB.Begin()
	//m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？

	//这个时候应该先查询表，然后确定这个订单是否已经扣减过库存了，已经扣减过了就别扣减了
	//并发时候会有漏洞， 同一个时刻发送了重复了多次， 使用锁，分布式锁
	sellDetail := model.StockSellDetail{
		OrderSn: req.OrderSn,
		Status:  1, //已扣减
	}
	var details []model.GoodsDetail

	for _, goodInfo := range req.GoodsInfo {
		details = append(details, model.GoodsDetail{
			Goods: goodInfo.GoodsId,
			Num:   goodInfo.Num,
		})

		var inv model.Inventory
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods:goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback() //回滚之前的操作
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}

		//for {

		mutex := global.RedisRs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
			//for err != nil {
			//	err = mutex.Lock()
			//}
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		//inv.Stocks -= goodInfo.Num
		//tx.Save(&inv)
		if err := tx.Model(&inv).Update("stocks", gorm.Expr("stocks - ?", goodInfo.Num)).Error; err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "更新语句执行失败")
		}
		//tx.Save(&model.Inventory{BaseModel: model.BaseModel{ID: inv.ID}, Stocks: inv.Stocks}) //高并发下更新失败？
		//tx.Exec("update inventory set stocks = ? where goods = ?", inv.Stocks, goodInfo.GoodsId)
		//var inv2 model.Inventory
		//global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv2)
		//for inv2.Stocks == inv.Stocks {
		//	tx.Save(&inv)
		//	global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv2)
		//}
		if ok, err := mutex.Unlock(); !ok || err != nil {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
		//update inventory set stocks = stocks-1, version=version+1 where goods=goods and version=version
		//这种写法有瑕疵，为什么？
		//零值 对于int类型来说 默认值是0 这种会被gorm给忽略掉
		//if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version= ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version+1}); result.RowsAffected == 0 {
		//	zap.S().Info("库存扣减失败")
		//}else{
		//	break
		//}
		//}
		//tx.Save(&inv)
	}
	sellDetail.Detail = details
	//写selldetail表
	if result := tx.Create(&sellDetail); result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "保存库存扣减历史失败")
	}
	tx.Commit() // 需要自己手动提交操作
	//m.Unlock() //释放锁
	return &proto.Empty{}, nil
}

// 归还库存功能由AutoReback重构 这个已废除
func (s *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*proto.Empty, error) {
	//库存归还：1.订单超时，2.订单创建失败 归还之前扣减的库存，3.手动归还
	tx := global.DB.Begin()
	m2.Lock() //获取锁
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := tx.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			//tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "库存信息不存在")
		}
		//扣减，会出现数据不一致的问题 - 分布式锁
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() //需要自己手动提交操作
	m2.Unlock() //释放锁
	return &proto.Empty{}, nil
}
func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range msgs {
		//既然是归还库存，那么我应该具体的知道每件商品应该归还多少,但是有一个问题是什么?重复归还的问题
		//所以说这个接口应该确保幂等性,你不能因为消息的重复发送导致一个订单的库存归还多次,没有扣减的库存你别归还
		//如何确保这些都没有问题，新建一张表，这张表记录了详细的订单扣减细节，以及归还细节
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("解析json失败：%v\n", msgs[i].Body)
			//根据业务来   订单号都解析失败了，感觉是错误的信息
			//ConsumeRetryLater 保证下次还能执行
			//ConsumeSuccess 丢弃
			return consumer.ConsumeSuccess, nil
		}
		//去将inv的库存加回去 将selldetail的status设置为2，要在事务中进行
		tx := global.DB.Begin()
		var sellDetail model.StockSellDetail
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn, Status: 1}).First(&sellDetail); result.Error != nil {
			return consumer.ConsumeSuccess, nil
		}
		//如果查询到 逐个归还库存
		for _, orderGood := range sellDetail.Detail {
			//update怎么用
			//先查询以下inv表在 ， update语句的update xx set stocks=stocks+2
			//mysql会自己锁住这个操作
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: orderGood.Goods}).Update("stocks", gorm.Expr("stocks + ?", orderGood.Num)); result.Error != nil {
				tx.Rollback()
				//ConsumeRetryLater 保证下次还能执行
				return consumer.ConsumeRetryLater, nil
			}
		}
		sellDetail.Status = 2
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn}).Update("status", 2); result.Error != nil {
			tx.Rollback()
			//ConsumeRetryLater 保证下次还能执行
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
