package initialize

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"order_srv/global"
	"order_srv/handler"
)

func InitRocketMQ() {

	//发送事务消息

	var err error
	//fmt.Println(global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)

	global.MQSendTranClient, err = rocketmq.NewTransactionProducer(
		&handler.OrderListener{},
		producer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)}),
		producer.WithGroupName("shiwu"),
	)
	if err != nil {
		fmt.Println("【事务消息】生成producer失败：%s", err.Error())
		panic(err.Error())
		//return nil, err
		//fmt.Println("连接错误：", err)
	}
	//启动
	if err = global.MQSendTranClient.Start(); err != nil {
		fmt.Println("【事务消息】启动producer失败：%s", err.Error())
		panic(err.Error())
		//return nil, err
	}

	//发送延时消息
	global.MQSendClient, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{fmt.Sprintf("%s:%d", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)})),
		producer.WithGroupName("yanshi"),
	)
	if err != nil {
		fmt.Println("【同步消息】生成producer失败：", err)
		panic(err.Error())
	}
	//启动
	err = global.MQSendClient.Start()
	if err != nil {
		fmt.Println("【同步消息】启动producer错误：", err)
	}

	//订阅消息 - 订单超时
	global.MQPushClient, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)}),
		consumer.WithGroupName(global.ServerConfig.RocketMQConfig.Name),
	)
	if err != nil {
		fmt.Printf("【订阅消息】生成producer失败：%s\n", err.Error())
		panic(err.Error())
	}
	if err = global.MQPushClient.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeout); err != nil {
		fmt.Printf("【订阅消息】失败：%s\n", err.Error())
		panic(err.Error())
	}
	//启动
	if err = global.MQPushClient.Start(); err != nil {
		fmt.Printf("【订阅消息】启动producer失败：%s\n", err.Error())
	}
}
func RegisterMQ() {
	err := global.MQSendClient.Shutdown()
	if err != nil {
		zap.S().Error("【同步消息】注销失败")
	}
	err = global.MQSendTranClient.Shutdown()
	if err != nil {
		zap.S().Error("【事务消息】注销失败")
	}
	err = global.MQPushClient.Shutdown()
	if err != nil {
		zap.S().Error("【订阅消息】注销失败")
	}
}
