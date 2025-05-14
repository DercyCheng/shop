package grpc

import (
	"context"
	"errors"
	"time"
	
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	
	pb "shop/backend/inventory/api/proto"
	"shop/backend/inventory/internal/domain/entity"
	"shop/backend/inventory/internal/service"
)

// InventoryServer 库存服务gRPC实现
type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer
	inventoryService    service.InventoryService
	inventoryLockService service.InventoryLockService
	warehouseService    service.WarehouseService
	logger             *zap.Logger
}

// NewInventoryServer 创建库存服务gRPC实例
func NewInventoryServer(
	inventoryService service.InventoryService,
	inventoryLockService service.InventoryLockService,
	warehouseService service.WarehouseService,
	logger *zap.Logger,
) *InventoryServer {
	return &InventoryServer{
		inventoryService:    inventoryService,
		inventoryLockService: inventoryLockService,
		warehouseService:    warehouseService,
		logger:             logger,
	}
}

// SetInv 设置商品库存
func (s *InventoryServer) SetInv(ctx context.Context, req *pb.GoodsInvInfo) (*emptypb.Empty, error) {
	if req.GoodsId <= 0 || req.Stock < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid goods_id or stock")
	}
	
	err := s.inventoryService.SetInventory(ctx, req.GoodsId, int(req.Stock), req.Operator)
	if err != nil {
		s.logger.Error("Failed to set inventory",
			zap.Int64("goods_id", req.GoodsId),
			zap.Int32("stock", req.Stock),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to set inventory: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// InvDetail 获取商品库存详情
func (s *InventoryServer) InvDetail(ctx context.Context, req *pb.GoodsInvInfo) (*pb.GoodsInvInfo, error) {
	if req.GoodsId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid goods_id")
	}
	
	inventory, err := s.inventoryService.GetInventory(ctx, req.GoodsId)
	if err != nil {
		if errors.Is(err, service.ErrStockNotFound) {
			return nil, status.Errorf(codes.NotFound, "inventory not found")
		}
		s.logger.Error("Failed to get inventory detail",
			zap.Int64("goods_id", req.GoodsId),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get inventory: %v", err)
	}
	
	return &pb.GoodsInvInfo{
		GoodsId:       inventory.ProductID,
		Stock:         int32(inventory.Stock),
		LockStock:     int32(inventory.LockStock),
		WarehouseId:   int32(inventory.WarehouseID),
		AlertThreshold: int32(inventory.AlertThreshold),
	}, nil
}

// BatchInvDetail 批量获取商品库存详情
func (s *InventoryServer) BatchInvDetail(ctx context.Context, req *pb.BatchGoodsInvInfo) (*pb.BatchGoodsInvInfo, error) {
	if len(req.GoodsList) == 0 {
		return &pb.BatchGoodsInvInfo{GoodsList: []*pb.GoodsInvInfo{}}, nil
	}
	
	// 提取商品ID列表
	var productIDs []int64
	for _, item := range req.GoodsList {
		productIDs = append(productIDs, item.GoodsId)
	}
	
	// 批量查询库存
	inventories, err := s.inventoryService.BatchGetInventory(ctx, productIDs)
	if err != nil {
		s.logger.Error("Failed to batch get inventory",
			zap.Any("product_ids", productIDs),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to batch get inventory: %v", err)
	}
	
	// 转换为proto格式
	result := &pb.BatchGoodsInvInfo{
		GoodsList: make([]*pb.GoodsInvInfo, 0, len(inventories)),
	}
	
	for _, inventory := range inventories {
		result.GoodsList = append(result.GoodsList, &pb.GoodsInvInfo{
			GoodsId:       inventory.ProductID,
			Stock:         int32(inventory.Stock),
			LockStock:     int32(inventory.LockStock),
			WarehouseId:   int32(inventory.WarehouseID),
			AlertThreshold: int32(inventory.AlertThreshold),
		})
	}
	
	return result, nil
}

// AddStock 添加库存
func (s *InventoryServer) AddStock(ctx context.Context, req *pb.AddStockInfo) (*emptypb.Empty, error) {
	if req.GoodsId <= 0 || req.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid goods_id or quantity")
	}
	
	err := s.inventoryService.AddStock(ctx, req.GoodsId, int(req.Quantity), req.Remark)
	if err != nil {
		s.logger.Error("Failed to add stock",
			zap.Int64("goods_id", req.GoodsId),
			zap.Int32("quantity", req.Quantity),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to add stock: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// AdjustStock 调整库存
func (s *InventoryServer) AdjustStock(ctx context.Context, req *pb.AdjustStockInfo) (*emptypb.Empty, error) {
	if req.GoodsId <= 0 || req.Stock < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid goods_id or stock")
	}
	
	err := s.inventoryService.AdjustStock(ctx, req.GoodsId, int(req.Stock), req.Operator, req.Remark)
	if err != nil {
		s.logger.Error("Failed to adjust stock",
			zap.Int64("goods_id", req.GoodsId),
			zap.Int32("stock", req.Stock),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to adjust stock: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// GetInventoryHistory 获取库存历史记录
func (s *InventoryServer) GetInventoryHistory(ctx context.Context, req *pb.InventoryHistoryRequest) (*pb.InventoryHistoryResponse, error) {
	if req.GoodsId <= 0 || req.Page <= 0 || req.PageSize <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters")
	}
	
	// 查询历史记录
	histories, total, err := s.inventoryService.GetInventoryHistory(ctx, req.GoodsId, int(req.Page), int(req.PageSize))
	if err != nil {
		s.logger.Error("Failed to get inventory history",
			zap.Int64("goods_id", req.GoodsId),
			zap.Int32("page", req.Page),
			zap.Int32("page_size", req.PageSize),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get inventory history: %v", err)
	}
	
	// 转换为proto格式
	response := &pb.InventoryHistoryResponse{
		Total: total,
		Items: make([]*pb.InventoryHistoryItem, 0, len(histories)),
	}
	
	for _, history := range histories {
		response.Items = append(response.Items, &pb.InventoryHistoryItem{
			Id:         history.ID,
			GoodsId:    history.ProductID,
			WarehouseId: int32(history.WarehouseID),
			Quantity:    int32(history.Quantity),
			Operation:   string(history.Operation),
			Operator:    history.Operator,
			OrderSn:     history.OrderSN,
			Remark:      history.Remark,
			CreatedAt:   timestamppb.New(history.CreatedAt),
		})
	}
	
	return response, nil
}

// Lock 锁定商品库存
func (s *InventoryServer) Lock(ctx context.Context, req *pb.SellInfo) (*pb.LockResponse, error) {
	if req.OrderSn == "" || len(req.GoodsList) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_sn or empty goods list")
	}
	
	// 转换请求为服务层需要的格式
	lockItems := make([]service.LockItem, 0, len(req.GoodsList))
	for _, item := range req.GoodsList {
		lockItems = append(lockItems, service.LockItem{
			ProductID: item.GoodsId,
			Quantity:  int(item.Quantity),
		})
	}
	
	// 设置默认超时时间
	timeoutSeconds := 30 * 60 // 默认30分钟
	if req.TimeoutSeconds > 0 {
		timeoutSeconds = int(req.TimeoutSeconds)
	}
	
	// 调用锁定库存服务
	result, err := s.inventoryLockService.LockInventory(ctx, req.OrderSn, lockItems, timeoutSeconds)
	if err != nil {
		s.logger.Error("Failed to lock inventory",
			zap.String("order_sn", req.OrderSn),
			zap.Any("items", lockItems),
			zap.Error(err))
		
		// 即使出错，也需要返回哪些商品锁定失败
		response := &pb.LockResponse{
			Success: false,
			Message: err.Error(),
		}
		
		if result != nil && len(result.FailItems) > 0 {
			for _, item := range result.FailItems {
				response.FailItems = append(response.FailItems, &pb.LockFailItem{
					GoodsId:   item.ProductID,
					Quantity:  int32(item.Quantity),
					Available: int32(item.Available),
					Reason:    item.Reason,
				})
			}
		}
		
		return response, nil
	}
	
	// 成功锁定
	response := &pb.LockResponse{
		Success: true,
		Message: "Lock inventory success",
	}
	
	// 如果有失败项（部分成功）
	if result != nil && len(result.FailItems) > 0 {
		response.Success = false
		response.Message = "Partially locked inventory"
		
		for _, item := range result.FailItems {
			response.FailItems = append(response.FailItems, &pb.LockFailItem{
				GoodsId:   item.ProductID,
				Quantity:  int32(item.Quantity),
				Available: int32(item.Available),
				Reason:    item.Reason,
			})
		}
	}
	
	return response, nil
}

// Sell 确认销售并扣减库存
func (s *InventoryServer) Sell(ctx context.Context, req *pb.SellInfo) (*emptypb.Empty, error) {
	if req.OrderSn == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_sn")
	}
	
	err := s.inventoryLockService.ConfirmReduce(ctx, req.OrderSn)
	if err != nil {
		s.logger.Error("Failed to confirm inventory reduction",
			zap.String("order_sn", req.OrderSn),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to confirm inventory reduction: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// Reback 归还库存
func (s *InventoryServer) Reback(ctx context.Context, req *pb.SellInfo) (*emptypb.Empty, error) {
	if req.OrderSn == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_sn")
	}
	
	err := s.inventoryLockService.UnlockInventory(ctx, req.OrderSn)
	if err != nil {
		s.logger.Error("Failed to unlock inventory",
			zap.String("order_sn", req.OrderSn),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to unlock inventory: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// GetReservationStatus 查询库存预定状态
func (s *InventoryServer) GetReservationStatus(ctx context.Context, req *pb.OrderSn) (*pb.ReservationStatus, error) {
	if req.OrderSn == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order_sn")
	}
	
	detail, err := s.inventoryLockService.GetLockDetail(ctx, req.OrderSn)
	if err != nil {
		if errors.Is(err, service.ErrLockNotFound) {
			return &pb.ReservationStatus{Status: pb.ReservationStatus_NOT_FOUND}, nil
		}
		s.logger.Error("Failed to get reservation status",
			zap.String("order_sn", req.OrderSn),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get reservation status: %v", err)
	}
	
	// 转换状态
	pbStatus := pb.ReservationStatus_UNKNOWN
	switch detail.Status {
	case entity.StockLocked:
		pbStatus = pb.ReservationStatus_LOCKED
	case entity.StockReduced:
		pbStatus = pb.ReservationStatus_REDUCED
	case entity.StockReturned:
		pbStatus = pb.ReservationStatus_RETURNED
	}
	
	// 提取商品明细
	goodsList := make([]*pb.GoodsSellInfo, 0)
	if len(detail.DetailItems) > 0 {
		for _, item := range detail.DetailItems {
			goodsList = append(goodsList, &pb.GoodsSellInfo{
				GoodsId:     item.ProductID,
				Quantity:    int32(item.Quantity),
				WarehouseId: int32(item.WarehouseID),
			})
		}
	}
	
	// 转换时间
	var lockTime, confirmTime *timestamppb.Timestamp
	if detail.LockTime != nil {
		lockTime = timestamppb.New(*detail.LockTime)
	}
	if detail.ConfirmTime != nil {
		confirmTime = timestamppb.New(*detail.ConfirmTime)
	}
	
	return &pb.ReservationStatus{
		Status:      pbStatus,
		LockTime:    lockTime,
		ConfirmTime: confirmTime,
		GoodsList:   goodsList,
	}, nil
}

// CreateWarehouse 创建仓库
func (s *InventoryServer) CreateWarehouse(ctx context.Context, req *pb.WarehouseInfo) (*pb.WarehouseInfo, error) {
	if req.Name == "" || req.Address == "" {
		return nil, status.Errorf(codes.InvalidArgument, "warehouse name and address are required")
	}
	
	// 转换为实体
	warehouse := &entity.Warehouse{
		Name:    req.Name,
		Address: req.Address,
		Contact: req.Contact,
		Phone:   req.Phone,
		Status:  int8(req.Status),
	}
	
	// 创建仓库
	err := s.warehouseService.CreateWarehouse(ctx, warehouse)
	if err != nil {
		if errors.Is(err, service.ErrWarehouseNameExists) {
			return nil, status.Errorf(codes.AlreadyExists, "warehouse with this name already exists")
		}
		s.logger.Error("Failed to create warehouse", 
			zap.String("name", req.Name),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create warehouse: %v", err)
	}
	
	// 返回创建后的仓库信息
	return &pb.WarehouseInfo{
		Id:        int32(warehouse.ID),
		Name:      warehouse.Name,
		Address:   warehouse.Address,
		Contact:   warehouse.Contact,
		Phone:     warehouse.Phone,
		Status:    int32(warehouse.Status),
		CreatedAt: timestamppb.New(warehouse.CreatedAt),
		UpdatedAt: timestamppb.New(warehouse.UpdatedAt),
	}, nil
}

// UpdateWarehouse 更新仓库
func (s *InventoryServer) UpdateWarehouse(ctx context.Context, req *pb.WarehouseInfo) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse id")
	}
	
	// 获取原始仓库信息
	existingWarehouse, err := s.warehouseService.GetWarehouse(ctx, int(req.Id))
	if err != nil {
		if errors.Is(err, service.ErrWarehouseNotFound) {
			return nil, status.Errorf(codes.NotFound, "warehouse not found")
		}
		s.logger.Error("Failed to get existing warehouse", 
			zap.Int32("id", req.Id),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get existing warehouse: %v", err)
	}
	
	// 更新字段
	if req.Name != "" {
		existingWarehouse.Name = req.Name
	}
	if req.Address != "" {
		existingWarehouse.Address = req.Address
	}
	if req.Contact != "" {
		existingWarehouse.Contact = req.Contact
	}
	if req.Phone != "" {
		existingWarehouse.Phone = req.Phone
	}
	if req.Status != 0 {
		existingWarehouse.Status = int8(req.Status)
	}
	
	// 更新仓库
	err = s.warehouseService.UpdateWarehouse(ctx, existingWarehouse)
	if err != nil {
		if errors.Is(err, service.ErrWarehouseNameExists) {
			return nil, status.Errorf(codes.AlreadyExists, "warehouse with this name already exists")
		}
		s.logger.Error("Failed to update warehouse", 
			zap.Int32("id", req.Id),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update warehouse: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// GetWarehouseList 获取仓库列表
func (s *InventoryServer) GetWarehouseList(ctx context.Context, req *pb.WarehouseQuery) (*pb.WarehouseListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	
	// 获取仓库列表
	warehouses, total, err := s.warehouseService.ListWarehouses(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		s.logger.Error("Failed to get warehouse list", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get warehouse list: %v", err)
	}
	
	// 转换为proto格式
	response := &pb.WarehouseListResponse{
		Total:      total,
		Warehouses: make([]*pb.WarehouseInfo, 0, len(warehouses)),
	}
	
	for _, warehouse := range warehouses {
		response.Warehouses = append(response.Warehouses, &pb.WarehouseInfo{
			Id:        int32(warehouse.ID),
			Name:      warehouse.Name,
			Address:   warehouse.Address,
			Contact:   warehouse.Contact,
			Phone:     warehouse.Phone,
			Status:    int32(warehouse.Status),
			CreatedAt: timestamppb.New(warehouse.CreatedAt),
			UpdatedAt: timestamppb.New(warehouse.UpdatedAt),
		})
	}
	
	return response, nil
}

// GetWarehouseDetail 获取仓库详情
func (s *InventoryServer) GetWarehouseDetail(ctx context.Context, req *pb.WarehouseID) (*pb.WarehouseInfo, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse id")
	}
	
	warehouse, err := s.warehouseService.GetWarehouse(ctx, int(req.Id))
	if err != nil {
		if errors.Is(err, service.ErrWarehouseNotFound) {
			return nil, status.Errorf(codes.NotFound, "warehouse not found")
		}
		s.logger.Error("Failed to get warehouse detail", 
			zap.Int32("id", req.Id),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get warehouse detail: %v", err)
	}
	
	return &pb.WarehouseInfo{
		Id:        int32(warehouse.ID),
		Name:      warehouse.Name,
		Address:   warehouse.Address,
		Contact:   warehouse.Contact,
		Phone:     warehouse.Phone,
		Status:    int32(warehouse.Status),
		CreatedAt: timestamppb.New(warehouse.CreatedAt),
		UpdatedAt: timestamppb.New(warehouse.UpdatedAt),
	}, nil
}

// DeleteWarehouse 删除仓库
func (s *InventoryServer) DeleteWarehouse(ctx context.Context, req *pb.WarehouseID) (*emptypb.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid warehouse id")
	}
	
	err := s.warehouseService.DeleteWarehouse(ctx, int(req.Id))
	if err != nil {
		if errors.Is(err, service.ErrWarehouseNotFound) {
			return nil, status.Errorf(codes.NotFound, "warehouse not found")
		}
		s.logger.Error("Failed to delete warehouse", 
			zap.Int32("id", req.Id),
			zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete warehouse: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}
