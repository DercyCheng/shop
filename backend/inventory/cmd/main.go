package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	
	"shop/backend/inventory/api/proto"
	"shop/backend/inventory/configs"
	"shop/backend/inventory/internal/domain/entity"
	"shop/backend/inventory/internal/repository"
	"shop/backend/inventory/internal/repository/cache"
	"shop/backend/inventory/internal/service"
	grpcServer "shop/backend/inventory/internal/web/grpc"
	"shop/backend/pkg/logger/zaplogger"
)

func main() {
	// 加载配置
	config, err := configs.LoadConfig("")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	
	// 初始化日志
	log := setupLogger(config.Logger)
	defer log.Sync()
	
	log.Info("Starting inventory service...")
	
	// 初始化数据库
	db := setupDatabase(config.Database, log)
	
	// 自动迁移数据表
	autoMigrateTables(db, log)
	
	// 初始化Redis连接
	redisClient := setupRedis(config.Redis, log)
	defer redisClient.Close()
	
	// 初始化缓存
	inventoryCache := cache.NewRedisInventoryCache(
		redisClient, 
		log, 
		config.Inventory.CacheTTL,
	)
	
	// 初始化仓储层
	inventoryRepo := repository.NewInventoryRepository(db, inventoryCache, log)
	warehouseRepo := repository.NewWarehouseRepository(db, log)
	
	// 初始化服务层
	inventoryService := service.NewInventoryService(inventoryRepo, log)
	inventoryLockService := service.NewInventoryLockService(inventoryRepo, log)
	warehouseService := service.NewWarehouseService(warehouseRepo, log)
	
	// 启动HTTP服务
	httpServer := setupHTTPServer(config.Server, log)
	
	// 启动gRPC服务
	grpcListener, grpcServer := setupGRPCServer(
		config.Server, 
		log, 
		inventoryService, 
		inventoryLockService, 
		warehouseService,
	)
	
	// 启动服务
	go func() {
		log.Info("Starting gRPC server", zap.Int("port", config.Server.GRPC.Port))
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal("Failed to start gRPC server", zap.Error(err))
		}
	}()
	
	go func() {
		log.Info("Starting HTTP server", zap.Int("port", config.Server.Port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()
	
	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Info("Shutting down server...")
	
	// 关闭HTTP服务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("HTTP server forced to shutdown", zap.Error(err))
	}
	
	// 关闭gRPC服务
	grpcServer.GracefulStop()
	
	log.Info("Server exited")
}

// 设置日志记录器
func setupLogger(config configs.LoggerConfig) *zap.Logger {
	logLevel := zaplogger.ParseLevel(config.Level)
	logger := zaplogger.New(zaplogger.Config{
		Level:       logLevel,
		Format:      config.Format,
		OutputPaths: []string{config.OutputFile, "stdout"},
		MaxSize:     config.MaxSize,
		MaxBackups:  config.MaxBackups,
		MaxAge:      config.MaxAge,
		Compress:    config.Compress,
	})
	
	return logger
}

// 设置数据库连接
func setupDatabase(config configs.DatabaseConfig, log *zap.Logger) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.Charset)
	
	gormLogger := logger.New(
		&zaplogger.GormLogWriter{Logger: log},
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection", zap.Error(err))
	}
	
	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	return db
}

// 自动迁移数据表
func autoMigrateTables(db *gorm.DB, log *zap.Logger) {
	log.Info("Auto-migrating database tables...")
	err := db.AutoMigrate(
		&entity.Inventory{},
		&entity.StockSellDetail{},
		&entity.Warehouse{},
		&entity.InventoryHistory{},
	)
	if err != nil {
		log.Fatal("Failed to auto-migrate tables", zap.Error(err))
	}
	log.Info("Database migration completed.")
}

// 设置Redis连接
func setupRedis(config configs.RedisConfig, log *zap.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	
	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	
	return client
}

// 设置HTTP服务器
func setupHTTPServer(config configs.ServerConfig, log *zap.Logger) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// 中间件
	router.Use(gin.Recovery())
	
	// 健康检查接口
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"name":   config.Name,
		})
	})
	
	// 创建HTTP服务器
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
}

// 设置gRPC服务器
func setupGRPCServer(
	config configs.ServerConfig,
	log *zap.Logger,
	inventoryService service.InventoryService,
	inventoryLockService service.InventoryLockService,
	warehouseService service.WarehouseService,
) (net.Listener, *grpc.Server) {
	// 创建gRPC服务器
	server := grpc.NewServer()
	
	// 注册服务
	proto.RegisterInventoryServiceServer(
		server,
		grpcServer.NewInventoryServer(
			inventoryService,
			inventoryLockService,
			warehouseService,
			log,
		),
	)
	
	// 注册反射服务，便于grpcurl等工具调试
	reflection.Register(server)
	
	// 监听端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GRPC.Port))
	if err != nil {
		log.Fatal("Failed to listen for gRPC server", zap.Error(err))
	}
	
	return lis, server
}
