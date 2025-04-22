package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"time"
)

type OrderListener struct {
	ID     int32
	Detail string
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	//执行逻辑并返回状态-自己决定
	return primitive.CommitMessageState
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	//执行回查逻辑并返回状态-自己决定
	return primitive.RollbackMessageState
}

func main() {
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{},
		producer.WithNameServer([]string{"192.168.10.130:9876"}),
	)
	if err != nil {
		zap.S().Error("生成producer失败：%s", err.Error())
	}
	//启动
	if err = p.Start(); err != nil {
		zap.S().Error("启动producer失败：%s", err.Error())
	}
	//发送半消息
	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_info", []byte("this is tranJzin")))
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", res.String())
	}
	//if res.State == primitive.CommitMessageState {
	//	fmt.Printf("发送失败: %s\n", err)
	//}
	//阻塞主线程
	time.Sleep(time.Hour)
	if err = p.Shutdown(); err != nil {
		fmt.Println("关闭producer失败")
	}
}
