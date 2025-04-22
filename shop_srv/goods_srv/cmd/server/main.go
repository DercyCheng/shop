package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"goods_srv/global"
	"goods_srv/handler"
	"goods_srv/initialize"
	proto "goods_srv/proto"
	"goods_srv/utils"
)

func main() {
	//IP := flag.String("ip", "192.168.1.106", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	zap.S().Info(global.ServerConfig)
	initialize.InitDB()
	//initialize.InitEs()

	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	//zap.S().Info("port：", *Port)

	server := grpc.NewServer()
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Panic("【商品服务-srv】注册失败")
		panic(err)
	} else {
		zap.S().Info("ip：", global.ServerConfig.Host, ":", *Port)
		zap.S().Info("【商品服务-srv】注册成功")
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check
	err = client.Agent().ServiceRegister(registration)
	//如何启动两个服务
	//即使我能够通过终端启动两个服务，但是注册到consul中的时候也会被覆盖
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()
	//接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Panic("【商品服务-srv】注销失败")
	} else {
		zap.S().Info("【商品服务-srv】注销成功")
	}
}
