package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"shop/backend/user/api/proto"
	"shop/backend/user/configs"
	"shop/backend/user/internal/repository"
	"shop/backend/user/internal/repository/cache"
	"shop/backend/user/internal/repository/dao"
	"shop/backend/user/internal/service"
	grpcHandler "shop/backend/user/internal/web/grpc"
	httpHandler "shop/backend/user/internal/web/http"
	"shop/backend/user/pkg/jwt"
)

func main() {
	// 加载配置
	cfg := loadConfig()

	// 初始化日志
	logger := initLogger(cfg.Log)
	defer logger.Sync()

	// 初始化数据库连接
	db := initDatabase(cfg.Database, logger)

	// 初始化Redis连接
	redisClient := initRedis(cfg.Redis, logger)

	// 初始化MongoDB连接
	mongoClient := initMongoDB(cfg.MongoDB, logger)
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.Error("Failed to disconnect MongoDB", zap.Error(err))
		}
	}()

	// 创建依赖项
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewRedisUserCache(redisClient, "shop:user", 30*time.Minute)
	userRepo := repository.NewUserRepository(userDAO, userCache)

	// 创建JWT工具
	jwtUtil := jwt.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)

	// 创建认证服务
	authService := service.NewAuthService(
		userRepo,
		jwtUtil,
		logger,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		5,           // 最大登录失败次数
		mongoClient, // 添加MongoDB客户端
	)

	// 创建用户服务
	userService := service.NewUserService(userRepo, authService, logger, mongoClient)

	// 启动gRPC服务器
	go startGRPCServer(cfg.Server.GRPCPort, authService, userService, logger)

	// 启动HTTP服务器
	go startHTTPServer(cfg.Server.HTTPPort, authService, userService, logger)

	// 等待终止信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// 优雅关闭
	logger.Info("Server exited")
}

// 加载配置
func loadConfig() *configs.Config {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}

// 初始化日志
func initLogger(cfg configs.LogConfig) *zap.Logger {
	var zapLogger *zap.Logger
	var err error

	// 根据环境配置日志
	if os.Getenv("APP_ENV") == "production" {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	return zapLogger
}

// 初始化数据库
func initDatabase(cfg configs.DatabaseConfig, logger *zap.Logger) *gorm.DB {
	// 配置GORM日志
	gormLogger := logger.With(zap.String("component", "gorm"))

	gormConfig := &gorm.Config{
		Logger: logger.Sugar().Level(&logger.Sugar().Desugar().Core()),
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.DSN), gormConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database connection", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	logger.Info("Database connection established")
	return db
}

// 初始化Redis
func initRedis(cfg configs.RedisConfig, logger *zap.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	logger.Info("Redis connection established")
	return client
}

// 初始化MongoDB
func initMongoDB(cfg configs.MongoDBConfig, logger *zap.Logger) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.URI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal("Failed to ping MongoDB", zap.Error(err))
	}

	logger.Info("MongoDB connection established")
	return client
}

// 启动gRPC服务器
func startGRPCServer(port string, authService service.AuthService, userService service.UserService, logger *zap.Logger) {
	addr := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpcHandler.NewUserGRPCServer(authService, userService)
	proto.RegisterUserServiceServer(grpcServer, grpcHandler.NewUserGRPCServer(authService, userService))

	logger.Info("Starting gRPC server", zap.String("addr", addr))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve gRPC", zap.Error(err))
	}
}

// 启动HTTP服务器
func startHTTPServer(port string, authService service.AuthService, userService service.UserService, logger *zap.Logger) {
	addr := fmt.Sprintf(":%s", port)
	httpServer := httpHandler.NewHTTPServer(addr, authService, userService, logger)

	logger.Info("Starting HTTP server", zap.String("addr", addr))
	if err := httpServer.Start(); err != nil {
		logger.Fatal("Failed to serve HTTP", zap.Error(err))
	}
}
