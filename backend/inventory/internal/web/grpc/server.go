package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "shop/inventory/api/proto"
	"shop/inventory/internal/service"
)

// Server represents the gRPC server for the inventory service
type Server struct {
	stockService       service.StockService
	warehouseService   service.WarehouseService
	reservationService service.ReservationService
	grpcServer         *grpc.Server
	port               int
}

// NewServer creates a new gRPC server instance
func NewServer(
	stockService service.StockService,
	warehouseService service.WarehouseService,
	reservationService service.ReservationService,
	port int,
) *Server {
	return &Server{
		stockService:       stockService,
		warehouseService:   warehouseService,
		reservationService: reservationService,
		port:               port,
	}
}

// Start initializes and starts the gRPC server
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %v", s.port, err)
	}

	// Create a new gRPC server instance
	s.grpcServer = grpc.NewServer()

	// Create the inventory service handler
	handler := NewInventoryServiceHandler(
		s.stockService,
		s.warehouseService,
		s.reservationService,
	)

	// Register the inventory service with the gRPC server
	pb.RegisterInventoryServiceServer(s.grpcServer, handler)

	// Register reflection service for gRPC tools like grpcurl
	reflection.Register(s.grpcServer)

	// Start serving gRPC requests
	fmt.Printf("Starting inventory gRPC server on port %d\n", s.port)
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	if s.grpcServer != nil {
		fmt.Println("Stopping inventory gRPC server...")
		s.grpcServer.GracefulStop()
		fmt.Println("Inventory gRPC server stopped")
	}
}
