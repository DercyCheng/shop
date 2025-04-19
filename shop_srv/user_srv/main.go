package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"nd/user_srv/global"
	"nd/user_srv/initialize"
	"nd/user_srv/proto"
	"nd/user_srv/utils"
	"nd/user_srv/wire"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/middleware"
	"google.golang.org/grpc/middleware/recovery"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号") // 这个修改为0，如果我们从命令行带参数启动的话就不会为0
	MetricsPort := flag.Int("metrics_port", 9100, "Prometheus指标暴露端口")

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.ServerConfig)

	// 初始化Prometheus指标暴露服务器
	initialize.InitPrometheusServer(*MetricsPort)

	flag.Parse()
	zap.S().Info("ip: ", *IP)
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("port: ", *Port)
	zap.S().Info("metrics_port: ", *MetricsPort)

	// 配置gRPC服务器选项
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 15 * time.Minute,  // 如果一个客户端空闲15分钟，发送一个GOAWAY
			Time:              5 * time.Minute,   // 每5分钟检测一次连接是否存活
			Timeout:           1 * time.Second,   // 如果1秒内未响应，则认为连接已断开
			MaxConnectionAge:  30 * time.Minute,  // 如果连接存在超过30分钟，发送一个GOAWAY
		}),
		grpc.UnaryInterceptor(middleware.ChainUnaryServer(
			recovery.UnaryServerInterceptor(),  // 添加panic恢复拦截器
			// 这里可以添加更多的拦截器，如日志、认证等
		)),
	}

	server := grpc.NewServer(opts...)
	
	// 使用Wire进行依赖注入，创建用户处理器
	userHandler := wire.ProvideUserHandler(global.DB)
	proto.RegisterUserServer(server, userHandler)
	
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = []string{"shop", "user", "srv"}
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	// 添加Prometheus指标端点到metadata中，方便服务发现
	registration.Meta = map[string]string{
		"metrics_path": "/metrics",
		"metrics_port": fmt.Sprintf("%d", *MetricsPort),
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	// 启动gRPC服务器
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal, 1) // 修改为buffer为1的channel，防止信号丢失
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	// 优雅关闭gRPC服务器
	server.GracefulStop()
	zap.S().Info("服务器已优雅关闭")
	
	// 注销服务
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Errorf("注销服务失败: %s", err.Error())
	} else {
		zap.S().Info("服务已成功注销")
	}
}
