package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"shop/backend/pkg/consul"
	"shop/backend/pkg/logger"
	"shop/backend/user/configs"
	"shop/backend/user/internal/service"
	"shop/backend/user/internal/web/grpc"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// 1. 加载配置
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化日志
	zapLogger := logger.InitLogger(cfg.LogLevel, cfg.LogFile)
	defer zapLogger.Sync()

	zap.ReplaceGlobals(zapLogger)
	zap.S().Info("User service starting...")

	// 3. 初始化依赖 (使用构造函数链而不是wire，简化实现)
	repo := initRepository(cfg)
	userService := service.NewUserService(repo)
	authService := service.NewAuthService(repo)

	// 4. 初始化服务器
	server := initGRPCServer(zapLogger, userService, authService)

	// 5. 服务注册
	serviceID, err := consul.RegisterService(cfg.Consul.Address, &consul.ServiceDefinition{
		ID:      fmt.Sprintf("user-service-%s-%d", cfg.Server.Host, cfg.Server.Port),
		Name:    "user-service",
		Address: cfg.Server.Host,
		Port:    cfg.Server.Port,
		Check: consul.ServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
		Tags: []string{"user", "auth", "api"},
	})
	if err != nil {
		zapLogger.Fatal("Failed to register service", zap.Error(err))
	}

	// 6. 启动gRPC服务器
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		zapLogger.Fatal("Failed to listen", zap.Error(err))
	}

	zap.S().Infof("User service is running on %s:%d", cfg.Server.Host, cfg.Server.Port)

	// 优雅退出处理
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := server.Serve(listen); err != nil {
			zapLogger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	// 7. 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.S().Info("Shutting down user service...")
	server.GracefulStop()

	// 8. 取消服务注册
	if err := consul.DeregisterService(cfg.Consul.Address, serviceID); err != nil {
		zapLogger.Error("Failed to deregister service", zap.Error(err))
	}

	zap.S().Info("User service stopped")
}

func initRepository(cfg *configs.Config) service.UserRepository {
	// In a real implementation, we would initialize the actual repository with DB connections
	// For now, return a simple implementation or mock
	return nil // placeholder, will be implemented later
}

func initGRPCServer(logger *zap.Logger, userService service.UserService, authService service.AuthService) *grpc.Server {
	// Create a new gRPC server
	server := grpc.NewServer()

	// Register health service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("user-service", healthpb.HealthCheckResponse_SERVING)

	// Register user service
	userServer := grpc.NewUserServiceServer(userService, authService)
	userServer.RegisterWithServer(server)

	return server
}
