package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"user_srv/global"
	"user_srv/handler"
	"user_srv/initialize"
	proto "user_srv/proto"
	"user_srv/utils"
	"user_srv/utils/register/consul"
)

func main() {
	//IP := flag.String("ip", "192.168.1.106", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	zap.S().Info(global.ServerConfig)
	initialize.InitDB()

	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
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
		zap.S().Panic("【用户服务-srv】服务注册失败:", err.Error())
	} else {
		zap.S().Info("【用户服务-srv】注册成功")
		zap.S().Info("ip：", global.ServerConfig.Host, ":", *Port)
	}
	//如何启动两个服务
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
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Panic("【用户服务-srv】服务注销失败:", err.Error())
	} else {
		zap.S().Info("【用户服务-srv】注销成功")
	}

}
