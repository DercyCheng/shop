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

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"shop/backend/profile/api/proto"
	"shop/backend/profile/configs"
	"shop/backend/profile/internal/repository"
	"shop/backend/profile/internal/service"
	"shop/backend/profile/internal/web/grpc"
	"shop/backend/profile/pkg/client"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := initLogger(config)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize database connections
	db, err := initDatabase(config)
	if err != nil {
		logger.Fatal("Failed to initialize database connections", zap.Error(err))
	}
	defer db.Close()

	mongoClient, err := initMongoDB(config)
	if err != nil {
		logger.Fatal("Failed to initialize MongoDB connection", zap.Error(err))
	}
	defer mongoClient.Disconnect(context.Background())

	// Initialize repositories
	userFavRepo := repository.NewUserFavRepository(db, logger)
	addressRepo := repository.NewAddressRepository(db, logger)
	messageRepo := repository.NewMessageRepository(mongoClient, config.MongoDB.Database, "messages", logger)

	// Initialize clients
	productClient, err := client.NewGRPCProductClient(config)
	if err != nil {
		logger.Fatal("Failed to initialize Product client", zap.Error(err))
	}
	defer productClient.Close()

	// Initialize services
	userFavService := service.NewUserFavService(userFavRepo, productClient, logger)
	addressService := service.NewAddressService(addressRepo, logger)
	messageService := service.NewMessageService(messageRepo, logger)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	profileServer := grpc.NewProfileGRPCServer(userFavService, addressService, messageService, logger)
	proto.RegisterProfileServiceServer(grpcServer, profileServer)

	// Register reflection service on gRPC server for grpcurl and grpc_cli
	reflection.Register(grpcServer)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Start gRPC server
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen", zap.String("addr", addr), zap.Error(err))
	}

	// Register service with Consul
	if err := registerServiceWithConsul(config); err != nil {
		logger.Warn("Failed to register service with Consul", zap.Error(err))
	}

	// Handle shutdown gracefully
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		logger.Info("Shutting down gRPC server")
		grpcServer.GracefulStop()
	}()

	logger.Info("Starting gRPC server", zap.String("addr", addr))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}

// initLogger initializes the logger
func initLogger(config *configs.Config) (*zap.Logger, error) {
	// Initialize logger configuration
	var cfg zap.Config

	if config.Log.Format == "json" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// Set log level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.Log.Level)); err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}
	cfg.Level.SetLevel(level)

	// Create logger
	logger, err := cfg.Build(
		zap.Fields(
			zap.String("service", config.Server.Name),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return logger, nil
}

// initDatabase initializes the database connection
func initDatabase(config *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MySQL.Username,
		config.MySQL.Password,
		config.MySQL.Host,
		config.MySQL.Port,
		config.MySQL.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(config.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.MySQL.ConnMaxLifetime) * time.Second)

	return db, nil
}

// initMongoDB initializes the MongoDB connection
func initMongoDB(config *configs.Config) (*mongo.Client, error) {
	// Create MongoDB client options
	clientOptions := options.Client().
		ApplyURI(config.MongoDB.URI).
		SetMaxPoolSize(config.MongoDB.MaxPoolSize).
		SetMinPoolSize(config.MongoDB.MinPoolSize)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the MongoDB server
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// registerServiceWithConsul registers the service with Consul
func registerServiceWithConsul(config *configs.Config) error {
	// Create a new client
	consulClient, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%d", config.Consul.Host, config.Consul.Port),
	})
	if err != nil {
		return fmt.Errorf("failed to create Consul client: %w", err)
	}

	// Register service
	serviceID := fmt.Sprintf("%s-%s-%d", config.Server.Name, config.Server.Host, config.Server.Port)
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    config.Server.Name,
		Address: config.Server.Host,
		Port:    config.Server.Port,
		Tags:    []string{"profile", "api", "grpc"},
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", config.Server.Host, config.Server.Port, ""),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	if err := consulClient.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service with Consul: %w", err)
	}

	// Deregister service on shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		if err := consulClient.Agent().ServiceDeregister(serviceID); err != nil {
			log.Printf("Failed to deregister service with Consul: %v", err)
		}
	}()

	return nil
}
