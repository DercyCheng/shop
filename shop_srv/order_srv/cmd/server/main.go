package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"order_srv/global"
	"order_srv/handler"
	"order_srv/initialize"
	proto "order_srv/proto/v1order"
	"order_srv/utils"
	"order_srv/utils/otgrpc"
	"order_srv/utils/register/consul"
)

func main() {
	//IP := flag.String("ip", "192.168.1.106", "ip地址")
	Port := flag.Int("port", 50053, "端口号")

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	//initialize.InitRedis()
	initialize.InitSrvConn()
	initialize.InitRocketMQ()
	initialize.Initjaeger()

	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	//zap.S().Info("port：", *Port)

	//初始化jaeger

	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(global.JaegerTracer)))
	proto.RegisterOrderServer(server, &handler.OrderServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("【订单和购物车服务-srv】服务注册失败:", err.Error())
	} else {
		zap.S().Info("ip：", global.ServerConfig.Host, ":", *Port)
		zap.S().Info("【订单和购物车服务-srv】注册成功")
	}

	//启动服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//监听订单超时topic
	//c, _ := rocketmq.NewPushConsumer(
	//	consumer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.ServerConfig.RocketMQConfig.Host, global.ServerConfig.RocketMQConfig.Port)}),
	//	consumer.WithGroupName("mxshop-order"),
	//)

	if err = global.MQPushClient.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeout); err != nil {
		fmt.Println("读取消息失败")
	}
	err = global.MQPushClient.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
	//接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//_ = c.Shutdown()
	_ = global.JaegerCloser.Close()
	initialize.RegisterMQ()
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Panic("【订单和购物车服务-srv】注销失败:", err.Error())
	} else {
		zap.S().Info("【订单和购物车服务-srv】注销成功")
	}
}
