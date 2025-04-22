package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	//连接recketmq
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.10.130:9876"}))
	if err != nil {
		fmt.Println("生成producer失败：", err)
	}
	//启动
	err = p.Start()
	if err != nil {
		fmt.Println("启动producer错误：", err)
	}
	//实例化消息
	msg := &primitive.Message{
		Topic: "jzin",
		Body:  []byte("this is jzin"),
	}
	//同步发送
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Printf("send message error: %s\n", err)
	} else {
		fmt.Printf("send message success: result=%s\n", res.String())
	}
	//关闭连接
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}
