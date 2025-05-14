package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	
	"shop/backend/pkg/consul"
	zapLogger "shop/backend/pkg/logger"
	"shop/backend/product/api/proto"
	"shop/backend/product/configs"
	"shop/backend/product/internal/repository"
	"shop/backend/product/internal/repository/cache"
	"shop/backend/product/internal/service"
	"shop/backend/product/internal/web/grpc"
)

func main() {
	// 1. 初始化配置
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// 2. 初始化日志
	zapConfig := zapLogger.NewConfig(cfg.LogLevel, cfg.LogFile)
	zapLogger.Init(zapConfig)
	defer zap.L().Sync()
	
	// 使用zap的sugar logger
	sugar := zap.L().Sugar()
	sugar.Infow("Starting product service", "host", cfg.Server.Host, "port", cfg.Server.Port)
	
	// 3. 初始化数据库连接
	db, err := initDB(cfg)
	if err != nil {
		sugar.Fatalw("Failed to initialize database", "error", err)
	}
	
	// 4. 初始化Redis
	redisClient, err := initRedis(cfg)
	if err != nil {
		sugar.Fatalw("Failed to initialize Redis", "error", err)
	}
	defer redisClient.Close()
	
	// 5. 初始化ElasticSearch
	esClient, err := initElasticsearch(cfg)
	if err != nil {
		sugar.Fatalw("Failed to initialize Elasticsearch", "error", err)
	}
	
	// 6. 初始化仓储层
	productCache := cache.NewRedisProductCache(redisClient)
	productRepo := repository.NewProductRepository(db, productCache)
	categoryRepo := repository.NewCategoryRepository(db, productCache)
	brandRepo := repository.NewBrandRepository(db, productCache)
	bannerRepo := repository.NewBannerRepository(db, productCache)
	searchRepo := service.NewElasticSearchRepository(esClient, productRepo, categoryRepo, brandRepo)
	
	// 初始化Elasticsearch索引
	if err := searchRepo.Init(context.Background()); err != nil {
		sugar.Warnw("Failed to initialize Elasticsearch index, will retry later", "error", err)
	}
	
	// 7. 初始化服务层
	productService := service.NewProductService(productRepo, categoryRepo, brandRepo, searchRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	brandService := service.NewBrandService(brandRepo)
	bannerService := service.NewBannerService(bannerRepo)
	searchService := service.NewSearchService(searchRepo, productRepo)
		// 8. 创建gRPC服务器
	grpcServer := grpc.NewServer(
		productService,
		categoryService,
		brandService,
		bannerService,
		searchService,
	)
	
	// 设置健康检查状态
	grpcServer.SetServingStatus("product.ProductService", healthpb.HealthCheckResponse_SERVING)
	
	// 9. 注册服务到Consul
	consulClient, err := api.NewClient(&api.Config{
		Address: cfg.Consul.Address,
	})
	if err != nil {
		sugar.Fatalw("Failed to create Consul client", "error", err)
	}
	
	serviceID := fmt.Sprintf("product-service-%s-%d", cfg.Server.Host, cfg.Server.Port)
	err = consul.RegisterService(consulClient, &consul.ServiceInstance{
		ID:      serviceID,
		Name:    "product-service",
		Address: cfg.Server.Host,
		Port:    cfg.Server.Port,
		Check: consul.ServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
		Tags: []string{"product", "goods", "api"},
	})
	if err != nil {
		sugar.Fatalw("Failed to register service", "error", err)
	}
		// 10. 启动gRPC服务器
	go func() {
		listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
		if err != nil {
			sugar.Fatalw("Failed to listen", "error", err)
		}
		
		sugar.Infof("Product service is running on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := grpcServer.GetGRPCServer().Serve(listen); err != nil {
			sugar.Fatalw("Failed to serve", "error", err)
		}
	}()
	
	// 11. 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	sugar.Info("Shutting down product service...")
	
	// 注销服务
	if err := consulClient.Agent().ServiceDeregister(serviceID); err != nil {
		sugar.Errorw("Failed to deregister service", "error", err)
	}
	
	// 优雅地关闭GRPC服务器
	grpcServer.GetGRPCServer().GracefulStop()
	
	sugar.Info("Product service stopped")
}

// 初始化数据库
func initDB(cfg *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)
	
	// 配置GORM
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: logger.Default.LogMode(logger.Silent),
	}
	
	return gorm.Open(mysql.Open(dsn), config)
}

// 初始化Redis
func initRedis(cfg *configs.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	
	// 检查连接
	_, err := redisClient.Ping(context.Background()).Result()
	return redisClient, err
}

// 初始化Elasticsearch
func initElasticsearch(cfg *configs.Config) (*elasticsearch.Client, error) {
	esConfig := elasticsearch.Config{
		Addresses: cfg.ElasticSearch.Addresses,
	}
	
	// 如果配置了认证信息
	if cfg.ElasticSearch.Username != "" && cfg.ElasticSearch.Password != "" {
		esConfig.Username = cfg.ElasticSearch.Username
		esConfig.Password = cfg.ElasticSearch.Password
	}
	
	return elasticsearch.NewClient(esConfig)
}
