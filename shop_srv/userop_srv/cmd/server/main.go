package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"userop_srv/handler"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"userop_srv/global"
	"userop_srv/initialize"
	proto "userop_srv/proto"
	"userop_srv/utils"
	"userop_srv/utils/register/consul"
)

func main() {
	//IP := flag.String("ip", "192.168.1.106", "ip地址")
	Port := flag.Int("port", 50054, "端口号")

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	//initialize.InitRedis()

	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	//zap.S().Info("port：", *Port)

	server := grpc.NewServer()
	//proto.RegisterAddressServer(server, &proto.UnimplementedAddressServer{})
	//proto.RegisterMessageServer(server, &proto.UnimplementedMessageServer{})
	//proto.RegisterUserFavServer(server, &proto.UnimplementedUserFavServer{})
	//proto.RegisterAddressServer(server, &handler.UserOpServer{})
	//proto.RegisterMessageServer(server, &handler.UserOpServer{})
	proto.RegisterUserOpServer(server, &handler.UserOpServer{})
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
		zap.S().Panic("【用户操作服务-srv】注册失败:", err.Error())
	} else {
		zap.S().Info("ip：", global.ServerConfig.Host, ":", *Port)
		zap.S().Info("【用户操作服务-srv】注册成功")
	}

	//启动服务
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
		zap.S().Panic("【用户操作服务-srv】注销失败:", err.Error())
	} else {
		zap.S().Info("【用户操作服务-srv】注销成功")
	}
}
