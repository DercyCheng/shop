package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	
	"shop/backend/product/api/proto"
	"shop/backend/product/internal/service"
)

// Server gRPC服务器
type Server struct {
	grpcServer   *grpc.Server
	healthServer *health.Server
	productHandler *ProductHandler
}

// NewServer 创建gRPC服务器
func NewServer(
	productService service.ProductService,
	categoryService service.CategoryService,
	brandService service.BrandService,
	bannerService service.BannerService,
	searchService service.SearchService,
	opts ...grpc.ServerOption,
) *Server {
	// 创建gRPC服务器
	grpcServer := grpc.NewServer(opts...)
	
	// 创建健康检查服务器
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	
	// 创建商品处理器
	productHandler := NewProductHandler(
		productService,
		categoryService,
		brandService,
		bannerService,
		searchService,
	)
	
	// 注册商品服务
	proto.RegisterProductServiceServer(grpcServer, productHandler)
	
	return &Server{
		grpcServer:   grpcServer,
		healthServer: healthServer,
		productHandler: productHandler,
	}
}

// GetGRPCServer 获取gRPC服务器实例
func (s *Server) GetGRPCServer() *grpc.Server {
	return s.grpcServer
}

// SetServingStatus 设置服务健康状态
func (s *Server) SetServingStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	s.healthServer.SetServingStatus(service, status)
}
