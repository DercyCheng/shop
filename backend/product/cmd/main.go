package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	pb "shop/product/api/proto"
	"shop/product/configs"
	"shop/product/internal/domain/entity"
	"shop/product/internal/repository"
	"shop/product/internal/repository/dao"
	"shop/product/internal/service"
	grpchandler "shop/product/internal/web/grpc"
)

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// Default config path if not specified
		configPath = filepath.Join("configs", "config.yaml")
	}

	config, err := configs.LoadConfig(configPath)
	if err != nil {
		log.Printf("Warning: Failed to load configuration file: %v, using defaults and environment variables", err)
	}

	// Database configuration
	dsn := config.Database.DSN

	// Initialize database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database schemas
	if err := autoMigrateDB(db); err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	// Initialize DAOs
	productDAO := dao.NewProductDAO(db)
	categoryDAO := dao.NewCategoryDAO(db)
	brandDAO := dao.NewBrandDAO(db)
	bannerDAO := dao.NewBannerDAO(db)
	categoryBrandDAO := dao.NewCategoryBrandDAO(db)

	// Initialize repositories
	productRepo := repository.NewProductRepository(productDAO)
	categoryRepo := repository.NewCategoryRepository(categoryDAO)
	brandRepo := repository.NewBrandRepository(brandDAO)
	bannerRepo := repository.NewBannerRepository(bannerDAO)
	categoryBrandRepo := repository.NewCategoryBrandRepository(categoryBrandDAO)

	// Initialize services
	productService := service.NewProductService(productRepo, categoryRepo, brandRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	brandService := service.NewBrandService(brandRepo)
	bannerService := service.NewBannerService(bannerRepo)
	categoryBrandService := service.NewCategoryBrandService(categoryBrandRepo, categoryRepo, brandRepo)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()

	// Register product service
	productGRPCHandler := grpchandler.NewProductGRPCServer(
		productService,
		categoryService,
		brandService,
		bannerService,
		categoryBrandService,
	)
	pb.RegisterProductServiceServer(grpcServer, productGRPCHandler)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("product", grpc_health_v1.HealthCheckResponse_SERVING)

	// Start gRPC server
	port := config.Server.Port

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Product service started on port %d", port)

	// Handle graceful shutdown
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down product service...")
	grpcServer.GracefulStop()
	log.Println("Product service shutdown complete")
}

// autoMigrateDB automatically migrates database schemas
func autoMigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Product{},
		&entity.ProductImage{},
		&entity.Category{},
		&entity.Brand{},
		&entity.Banner{},
		&entity.CategoryBrand{},
	)
}
