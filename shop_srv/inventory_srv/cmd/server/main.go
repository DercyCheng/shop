package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"

	"inventory_srv/global"
	"inventory_srv/handler"
	"inventory_srv/initialize"
	proto "inventory_srv/proto"
	"inventory_srv/utils"
	"inventory_srv/utils/register/consul"
)

func main() {
	//IP := flag.String("ip", "192.168.1.106", "ip地址")
	Port := flag.Int("port", 50052, "端口号")

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	zap.S().Info(global.ServerConfig)
	initialize.InitDB()
	initialize.InitRedis()

	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("ip：", global.ServerConfig.Host, ":", *Port)
	//zap.S().Info("port：", *Port)

	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})
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
		zap.S().Panic("【库存服务-srv】服务注册失败:", err.Error())
		panic(err)
	} else {
		zap.S().Info("【库存服务-srv】注册成功")
	}

	//启动服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//订阅rocketMQ消息队列，监听库存归还topic
	//启动recketmq并设置负载均衡的Group
	c, _ := rocketmq.NewPushConsumer(
		//consumer.WithNameServer([]string{"192.168.10.130:9876"}),
		consumer.WithNameServer([]string{fmt.Sprintf("%s:%d", global.ServerConfig.Rocketmq.Host, global.ServerConfig.Rocketmq.Port)}),
		//consumer.WithGroupName("mxshop-inventory"),
		consumer.WithGroupName(global.ServerConfig.Rocketmq.Group),
	)
	//订阅消息
	if err = c.Subscribe(global.ServerConfig.Rocketmq.Suborder1, consumer.MessageSelector{}, handler.AutoReback); err != nil {
		fmt.Println(err.Error())
	}
	//启动
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
	/*//阻塞主线程
	time.Sleep(time.Hour)
	//关闭连接
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}*/
	//接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Panic("【库存服务-srv】注销失败:", err.Error())
	} else {
		zap.S().Info("【库存服务-srv】注销成功:")
	}
}
