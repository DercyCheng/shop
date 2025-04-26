package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	pb "shop/backend/user/api/proto/user"
	"shop/backend/user/configs"
	"shop/backend/user/internal/repository"
	"shop/backend/user/internal/repository/cache"
	"shop/backend/user/internal/repository/dao"
	"shop/backend/user/internal/service"
	grpcHandler "shop/backend/user/internal/web/grpc"
	httpHandler "shop/backend/user/internal/web/http"
)

func main() {
	// 加载配置
	cfg := loadConfig()

	// 初始化数据库连接
	db := initDatabase(cfg.Database)

	// 初始化Redis连接
	redisClient := initRedis(cfg.Redis)

	// 创建依赖项
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewRedisUserCache(redisClient, "shop:user", 30*time.Minute)
	userRepo := repository.NewUserRepository(userDAO, userCache)
	authService := service.NewAuthService(
		userRepo,
		redisClient,
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.VerificationTTL,
	)
	userService := service.NewUserService(userRepo, userCache, authService)

	// 启动gRPC服务
	go startGRPCServer(cfg.Server.GRPC, userService, authService)

	// 启动HTTP服务
	startHTTPServer(cfg.Server.HTTP, userService, authService)
}

// 加载配置
func loadConfig() *configs.Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var cfg configs.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}

	return &cfg
}

// 初始化数据库连接
func initDatabase(cfg configs.DatabaseConfig) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 自动迁移数据库表
	if err := db.AutoMigrate(&dao.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db
}

// 初始化Redis连接
func initRedis(cfg configs.RedisConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return redisClient
}

// 启动gRPC服务
func startGRPCServer(cfg configs.GRPCConfig, userService service.UserService, authService service.AuthService) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userGRPCServer := grpcHandler.NewUserGRPCServer(userService, authService)
	pb.RegisterUserServiceServer(grpcServer, userGRPCServer)

	log.Printf("gRPC server starting on port %d", cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// 启动HTTP服务
func startHTTPServer(cfg configs.HTTPConfig, userService service.UserService, authService service.AuthService) {
	router := gin.Default()

	// 注册HTTP路由
	userHandler := httpHandler.NewUserHandler(userService, authService)
	userHandler.RegisterRoutes(router)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// 优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		log.Println("Server exiting")
	}()

	log.Printf("HTTP server starting on port %d", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
