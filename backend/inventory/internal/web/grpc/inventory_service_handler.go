package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "shop/inventory/api/proto"
	"shop/inventory/internal/domain/entity"
	"shop/inventory/internal/service"
)

// InventoryServiceHandler implements the gRPC InventoryService interface
type InventoryServiceHandler struct {
	pb.UnimplementedInventoryServiceServer
	stockService       service.StockService
	warehouseService   service.WarehouseService
	reservationService service.ReservationService
}

// NewInventoryServiceHandler creates a new InventoryServiceHandler
func NewInventoryServiceHandler(
	stockService service.StockService,
	warehouseService service.WarehouseService,
	reservationService service.ReservationService,
) *InventoryServiceHandler {
	return &InventoryServiceHandler{
		stockService:       stockService,
		warehouseService:   warehouseService,
		reservationService: reservationService,
	}
}

// SetStock sets or creates stock for a product in a warehouse
func (h *InventoryServiceHandler) SetStock(ctx context.Context, req *pb.SetStockRequest) (*pb.StockInfo, error) {
	// Validate request
	if req.ProductId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID")
	}
	if req.WarehouseId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse ID")
	}
	if req.Quantity < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "quantity cannot be negative")
	}

	// Check if stock already exists for this product and warehouse
	existingStock, err := h.stockService.GetStockByProductAndWarehouse(ctx, req.ProductId, req.WarehouseId)
	if err == nil && existingStock != nil {
		// Update existing stock
		existingStock.Quantity = int(req.Quantity)
		if req.LowStockThreshold > 0 {
			existingStock.LowStockThreshold = int(req.LowStockThreshold)
		}

		err = h.stockService.UpdateStock(ctx, existingStock)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update stock: %v", err)
		}

		return convertStockToProto(existingStock), nil
	}

	// Create new stock record
	newStock := &entity.Stock{
		ProductID:         req.ProductId,
		WarehouseID:       req.WarehouseId,
		Quantity:          int(req.Quantity),
		LowStockThreshold: int(req.LowStockThreshold),
	}

	err = h.stockService.CreateStock(ctx, newStock)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create stock: %v", err)
	}

	// Get the created stock to ensure we have the complete data
	createdStock, err := h.stockService.GetStockByProductAndWarehouse(ctx, req.ProductId, req.WarehouseId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "stock created but failed to retrieve: %v", err)
	}

	return convertStockToProto(createdStock), nil
}

// GetStock retrieves stock information for a product in a warehouse
func (h *InventoryServiceHandler) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.StockInfo, error) {
	// Validate request
	if req.ProductId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID")
	}
	if req.WarehouseId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse ID")
	}

	// Get stock information
	stock, err := h.stockService.GetStockByProductAndWarehouse(ctx, req.ProductId, req.WarehouseId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "stock not found: %v", err)
	}

	return convertStockToProto(stock), nil
}

// ListStocks retrieves a paginated list of stocks
func (h *InventoryServiceHandler) ListStocks(ctx context.Context, req *pb.ListStocksRequest) (*pb.ListStocksResponse, error) {
	// Set default pagination if not provided
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// Get stocks
	stocks, total, err := h.stockService.ListStocks(ctx, page, pageSize, req.WarehouseId, req.LowStockOnly, req.OutOfStockOnly)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list stocks: %v", err)
	}

	// Convert to proto message
	response := &pb.ListStocksResponse{
		Total: int32(total),
	}

	for _, stock := range stocks {
		response.Stocks = append(response.Stocks, convertStockToProto(stock))
	}

	return response, nil
}

// IncrementStock increases the stock quantity for a product in a warehouse
func (h *InventoryServiceHandler) IncrementStock(ctx context.Context, req *pb.IncrementStockRequest) (*pb.StockInfo, error) {
	// Validate request
	if req.ProductId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID")
	}
	if req.WarehouseId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse ID")
	}
	if req.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "quantity must be greater than zero")
	}

	// Add stock
	stock, err := h.stockService.AddStock(ctx, req.ProductId, req.WarehouseId, int(req.Quantity))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to increment stock: %v", err)
	}

	return convertStockToProto(stock), nil
}

// DecrementStock decreases the stock quantity for a product in a warehouse
func (h *InventoryServiceHandler) DecrementStock(ctx context.Context, req *pb.DecrementStockRequest) (*pb.StockInfo, error) {
	// Validate request
	if req.ProductId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID")
	}
	if req.WarehouseId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse ID")
	}
	if req.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "quantity must be greater than zero")
	}

	// Check if there's enough stock available
	available, err := h.stockService.CheckStockAvailability(ctx, req.ProductId, req.WarehouseId, int(req.Quantity))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check stock availability: %v", err)
	}

	if !available {
		return nil, status.Errorf(codes.FailedPrecondition, "insufficient stock available")
	}

	// Remove stock
	stock, err := h.stockService.RemoveStock(ctx, req.ProductId, req.WarehouseId, int(req.Quantity))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decrement stock: %v", err)
	}

	return convertStockToProto(stock), nil
}

// BatchGetStocks retrieves stock information for multiple products in a warehouse
func (h *InventoryServiceHandler) BatchGetStocks(ctx context.Context, req *pb.BatchGetStocksRequest) (*pb.BatchGetStocksResponse, error) {
	// Validate request
	if len(req.ProductIds) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "product IDs cannot be empty")
	}
	if req.WarehouseId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse ID")
	}

	// Get stocks for the products
	stocks, err := h.stockService.GetStocksByProductIDs(ctx, req.ProductIds, req.WarehouseId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get stocks: %v", err)
	}

	// Convert to proto message
	response := &pb.BatchGetStocksResponse{}
	for _, stock := range stocks {
		response.Stocks = append(response.Stocks, convertStockToProto(stock))
	}

	return response, nil
}

// ReserveStock creates a temporary reservation of stock
func (h *InventoryServiceHandler) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReservationResponse, error) {
	// Validate request
	if req.OrderId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "order ID cannot be empty")
	}
	if len(req.Items) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "items cannot be empty")
	}

	// Convert items to domain model
	items := make([]*entity.ReservationItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, &entity.ReservationItem{
			ProductID:   item.ProductId,
			WarehouseID: item.WarehouseId,
			Quantity:    int(item.Quantity),
		})
	}

	// Set expiration minutes (use default if not provided)
	expirationMinutes := int(req.ExpirationMinutes)
	if expirationMinutes <= 0 {
		expirationMinutes = 30 // Default: 30 minutes
	}

	// Create reservation
	reservation, err := h.reservationService.CreateReservation(ctx, req.OrderId, items, expirationMinutes)
	if err != nil {
		// Handle error cases separately
		if errors.Is(err, errors.New("reservation already exists for this order")) {
			return nil, status.Errorf(codes.AlreadyExists, "reservation already exists for order %s", req.OrderId)
		}

		// Check if it's a stock availability error
		if errors.Is(err, errors.New("insufficient stock for reservation")) {
			return &pb.ReservationResponse{
				Success: false,
				Errors: []*pb.ReservationError{
					{
						Message: err.Error(),
					},
				},
			}, nil
		}

		return nil, status.Errorf(codes.Internal, "failed to create reservation: %v", err)
	}

	// Convert to proto message
	return &pb.ReservationResponse{
		Success:     true,
		Reservation: convertReservationToProto(reservation),
	}, nil
}

// CommitReservation confirms a reservation and permanently reduces stock
func (h *InventoryServiceHandler) CommitReservation(ctx context.Context, req *pb.CommitReservationRequest) (*emptypb.Empty, error) {
	// Validate request
	if req.ReservationId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid reservation ID")
	}

	// Commit the reservation
	err := h.reservationService.CommitReservation(ctx, req.ReservationId)
	if err != nil {
		if errors.Is(err, errors.New("reservation cannot be committed: it is either expired or not in pending state")) {
			return nil, status.Errorf(codes.FailedPrecondition, "reservation cannot be committed: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to commit reservation: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// CancelReservation cancels a reservation and returns stock
func (h *InventoryServiceHandler) CancelReservation(ctx context.Context, req *pb.CancelReservationRequest) (*emptypb.Empty, error) {
	// Validate request
	if req.ReservationId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid reservation ID")
	}

	// Cancel the reservation
	err := h.reservationService.CancelReservation(ctx, req.ReservationId)
	if err != nil {
		if errors.Is(err, errors.New("reservation cannot be cancelled: it is not in pending state")) {
			return nil, status.Errorf(codes.FailedPrecondition, "reservation cannot be cancelled: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to cancel reservation: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetReservation retrieves reservation information
func (h *InventoryServiceHandler) GetReservation(ctx context.Context, req *pb.GetReservationRequest) (*pb.ReservationInfo, error) {
	// Validate request
	if req.ReservationId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid reservation ID")
	}

	// Get the reservation
	reservation, err := h.reservationService.GetReservationByID(ctx, req.ReservationId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "reservation not found: %v", err)
	}

	return convertReservationToProto(reservation), nil
}

// CreateWarehouse creates a new warehouse
func (h *InventoryServiceHandler) CreateWarehouse(ctx context.Context, req *pb.CreateWarehouseRequest) (*pb.WarehouseInfo, error) {
	// Validate request
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "warehouse name cannot be empty")
	}

	// Create warehouse entity
	warehouse := &entity.Warehouse{
		Name:     req.Name,
		Address:  req.Address,
		IsActive: req.IsActive,
	}

	// Create warehouse
	err := h.warehouseService.CreateWarehouse(ctx, warehouse)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create warehouse: %v", err)
	}

	// Get the warehouse to ensure we have the complete data with ID
	// Note: This is a bit of a workaround since we don't have direct access to the ID after creation
	// In a real application, you might want to return the ID directly from the CreateWarehouse method
	// or have a better way to get the newly created warehouse
	warehouses, _, err := h.warehouseService.ListWarehouses(ctx, 1, 100, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "warehouse created but failed to retrieve: %v", err)
	}

	var createdWarehouse *entity.Warehouse
	for _, w := range warehouses {
		if w.Name == req.Name && w.Address == req.Address {
			createdWarehouse = w
			break
		}
	}

	if createdWarehouse == nil {
		return nil, status.Errorf(codes.Internal, "warehouse created but not found in list")
	}

	return convertWarehouseToProto(createdWarehouse), nil
}

// UpdateWarehouse updates an existing warehouse
func (h *InventoryServiceHandler) UpdateWarehouse(ctx context.Context, req *pb.UpdateWarehouseRequest) (*pb.WarehouseInfo, error) {
	// Validate request
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse ID")
	}
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "warehouse name cannot be empty")
	}

	// Get existing warehouse to verify it exists
	existingWarehouse, err := h.warehouseService.GetWarehouseByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "warehouse not found: %v", err)
	}

	// Update warehouse entity
	existingWarehouse.Name = req.Name
	existingWarehouse.Address = req.Address
	existingWarehouse.IsActive = req.IsActive

	// Update warehouse
	err = h.warehouseService.UpdateWarehouse(ctx, existingWarehouse)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update warehouse: %v", err)
	}

	// Get updated warehouse
	updatedWarehouse, err := h.warehouseService.GetWarehouseByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "warehouse updated but failed to retrieve: %v", err)
	}

	return convertWarehouseToProto(updatedWarehouse), nil
}

// DeleteWarehouse deletes a warehouse
func (h *InventoryServiceHandler) DeleteWarehouse(ctx context.Context, req *pb.DeleteWarehouseRequest) (*emptypb.Empty, error) {
	// Validate request
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse ID")
	}

	// Delete warehouse
	err := h.warehouseService.DeleteWarehouse(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete warehouse: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetWarehouse retrieves warehouse information
func (h *InventoryServiceHandler) GetWarehouse(ctx context.Context, req *pb.GetWarehouseRequest) (*pb.WarehouseInfo, error) {
	// Validate request
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse ID")
	}

	// Get warehouse
	warehouse, err := h.warehouseService.GetWarehouseByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "warehouse not found: %v", err)
	}

	return convertWarehouseToProto(warehouse), nil
}

// ListWarehouses retrieves a paginated list of warehouses
func (h *InventoryServiceHandler) ListWarehouses(ctx context.Context, req *pb.ListWarehousesRequest) (*pb.ListWarehousesResponse, error) {
	// Set default pagination if not provided
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// Get warehouses
	warehouses, total, err := h.warehouseService.ListWarehouses(ctx, page, pageSize, req.ActiveOnly)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list warehouses: %v", err)
	}

	// Convert to proto message
	response := &pb.ListWarehousesResponse{
		Total: int32(total),
	}

	for _, warehouse := range warehouses {
		response.Warehouses = append(response.Warehouses, convertWarehouseToProto(warehouse))
	}

	return response, nil
}

// Helper functions to convert between domain models and protobuf messages

// convertStockToProto converts a Stock entity to a StockInfo protobuf message
func convertStockToProto(stock *entity.Stock) *pb.StockInfo {
	if stock == nil {
		return nil
	}

	return &pb.StockInfo{
		Id:                stock.ID,
		ProductId:         stock.ProductID,
		WarehouseId:       stock.WarehouseID,
		Quantity:          int32(stock.Quantity),
		Reserved:          int32(stock.Reserved),
		Available:         int32(stock.Available()),
		LowStockThreshold: int32(stock.LowStockThreshold),
		InStock:           stock.IsInStock(),
		CreatedAt:         timestamppb.New(stock.CreatedAt),
		UpdatedAt:         timestamppb.New(stock.UpdatedAt),
	}
}

// convertReservationToProto converts a Reservation entity to a ReservationInfo protobuf message
func convertReservationToProto(reservation *entity.Reservation) *pb.ReservationInfo {
	if reservation == nil {
		return nil
	}

	// Convert status
	var status pb.ReservationStatus
	switch reservation.Status {
	case entity.ReservationPending:
		status = pb.ReservationStatus_PENDING
	case entity.ReservationCommitted:
		status = pb.ReservationStatus_COMMITTED
	case entity.ReservationCancelled:
		status = pb.ReservationStatus_CANCELLED
	case entity.ReservationExpired:
		status = pb.ReservationStatus_EXPIRED
	}

	// Convert items
	items := make([]*pb.ReservationItemInfo, 0, len(reservation.Items))
	for _, item := range reservation.Items {
		items = append(items, &pb.ReservationItemInfo{
			ProductId:   item.ProductID,
			WarehouseId: item.WarehouseID,
			Quantity:    int32(item.Quantity),
		})
	}

	return &pb.ReservationInfo{
		Id:        reservation.ID,
		OrderId:   reservation.OrderID,
		Status:    status,
		Items:     items,
		CreatedAt: timestamppb.New(reservation.CreatedAt),
		UpdatedAt: timestamppb.New(reservation.UpdatedAt),
		ExpiresAt: timestamppb.New(reservation.ExpiresAt),
	}
}

// convertWarehouseToProto converts a Warehouse entity to a WarehouseInfo protobuf message
func convertWarehouseToProto(warehouse *entity.Warehouse) *pb.WarehouseInfo {
	if warehouse == nil {
		return nil
	}

	return &pb.WarehouseInfo{
		Id:        warehouse.ID,
		Name:      warehouse.Name,
		Address:   warehouse.Address,
		IsActive:  warehouse.IsActive,
		CreatedAt: timestamppb.New(warehouse.CreatedAt),
		UpdatedAt: timestamppb.New(warehouse.UpdatedAt),
	}
}
