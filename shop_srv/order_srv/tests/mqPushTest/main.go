package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"os"
	"time"
)

func main() {
	//启动recketmq并设置负载均衡的Group
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.10.130:9876"}),
		consumer.WithGroupName("jzins"),
	)
	//订阅消息
	if err := c.Subscribe("jzin", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		fmt.Println(err.Error())
	}
	//启动
	err := c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	//阻塞主线程
	time.Sleep(time.Hour)
	//关闭连接
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}

}
