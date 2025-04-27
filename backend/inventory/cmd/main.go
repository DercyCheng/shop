package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"shop/inventory/configs"
	"shop/inventory/internal/domain/entity"
	"shop/inventory/internal/repository"
	"shop/inventory/internal/repository/dao"
	"shop/inventory/internal/service"
	"shop/inventory/internal/web/grpc"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	// Load configuration
	config, err := configs.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := initDatabase(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto-migrate database schema
	if err := migrateDatabase(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize DAOs
	stockDAO := dao.NewStockDAO(db)
	warehouseDAO := dao.NewWarehouseDAO(db)
	reservationDAO := dao.NewReservationDAO(db)

	// Initialize repositories
	stockRepo := repository.NewStockRepository(stockDAO)
	warehouseRepo := repository.NewWarehouseRepository(warehouseDAO)
	reservationRepo := repository.NewReservationRepository(reservationDAO)

	// Initialize services
	stockService := service.NewStockService(stockRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	reservationService := service.NewReservationService(reservationRepo, stockRepo)

	// Run expired reservations processor in a background goroutine
	go runReservationProcessor(reservationService)

	// Initialize gRPC server
	server := grpc.NewServer(
		stockService,
		warehouseService,
		reservationService,
		config.Server.Port,
	)

	// Graceful shutdown
	go handleGracefulShutdown(server)

	// Start the server (this will block until the server is stopped)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDatabase initializes the database connection
func initDatabase(config *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.ConnMaxLifetimeMinutes) * time.Minute)

	return db, nil
}

// migrateDatabase automatically migrates the database schema
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Stock{},
		&entity.Warehouse{},
		&entity.Reservation{},
		&entity.ReservationItem{},
	)
}

// runReservationProcessor starts a background task to periodically process expired reservations
func runReservationProcessor(reservationService service.ReservationService) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			count, err := reservationService.ProcessExpiredReservations(ctx)
			if err != nil {
				log.Printf("Error processing expired reservations: %v\n", err)
			} else if count > 0 {
				log.Printf("Processed %d expired reservations\n", count)
			}
			cancel()
		}
	}
}

// handleGracefulShutdown sets up signal handlers for graceful shutdown
func handleGracefulShutdown(server *grpc.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-sigCh
	log.Println("Received shutdown signal. Shutting down gracefully...")

	// Stop the server
	server.Stop()

	log.Println("Server shutdown complete")
	os.Exit(0)
}
