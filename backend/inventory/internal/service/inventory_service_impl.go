package service

import (
	"context"
	"errors"
	"time"
	
	"go.uber.org/zap"
	
	"shop/backend/inventory/internal/domain/entity"
	"shop/backend/inventory/internal/repository"
)

// 定义错误
var (
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrStockNotFound      = errors.New("stock not found")
	ErrOperationFailed    = errors.New("operation failed")
)

// InventoryServiceImpl 库存服务实现
type InventoryServiceImpl struct {
	repo   repository.InventoryRepository
	logger *zap.Logger
}

// NewInventoryService 创建库存服务实例
func NewInventoryService(repo repository.InventoryRepository, logger *zap.Logger) InventoryService {
	return &InventoryServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

// GetInventory 获取库存信息
func (s *InventoryServiceImpl) GetInventory(ctx context.Context, productID int64) (*entity.Inventory, error) {
	// 默认使用仓库ID为1
	const defaultWarehouseID = 1
	
	inventory, err := s.repo.GetInventory(ctx, productID, defaultWarehouseID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrStockNotFound
		}
		s.logger.Error("Failed to get inventory",
			zap.Int64("product_id", productID),
			zap.Error(err))
		return nil, err
	}
	
	return inventory, nil
}

// BatchGetInventory 批量获取库存信息
func (s *InventoryServiceImpl) BatchGetInventory(ctx context.Context, productIDs []int64) ([]*entity.Inventory, error) {
	// 默认使用仓库ID为1
	const defaultWarehouseID = 1
	
	if len(productIDs) == 0 {
		return []*entity.Inventory{}, nil
	}
	
	inventories, err := s.repo.BatchGetInventory(ctx, productIDs, defaultWarehouseID)
	if err != nil {
		s.logger.Error("Failed to batch get inventory",
			zap.Any("product_ids", productIDs),
			zap.Error(err))
		return nil, err
	}
	
	return inventories, nil
}

// GetAvailableStock 获取可用库存数量
func (s *InventoryServiceImpl) GetAvailableStock(ctx context.Context, productID int64) (int, error) {
	inventory, err := s.GetInventory(ctx, productID)
	if err != nil {
		return 0, err
	}
	
	return inventory.AvailableStock(), nil
}

// CheckStockAvailable 检查库存是否充足
func (s *InventoryServiceImpl) CheckStockAvailable(ctx context.Context, productID int64, quantity int) (bool, error) {
	if quantity <= 0 {
		return false, ErrInvalidArgument
	}
	
	inventory, err := s.GetInventory(ctx, productID)
	if err != nil {
		return false, err
	}
	
	return inventory.IsAvailable(quantity), nil
}

// SetInventory 设置商品库存
func (s *InventoryServiceImpl) SetInventory(ctx context.Context, productID int64, stock int, operator string) error {
	if productID <= 0 || stock < 0 {
		return ErrInvalidArgument
	}
	
	const defaultWarehouseID = 1
	
	// 先查询是否存在
	inventory, err := s.repo.GetInventory(ctx, productID, defaultWarehouseID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			// 不存在则创建新记录
			now := time.Now()
			inventory = &entity.Inventory{
				ProductID:   productID,
				Stock:       stock,
				WarehouseID: defaultWarehouseID,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
		} else {
			s.logger.Error("Failed to get inventory when setting",
				zap.Int64("product_id", productID),
				zap.Error(err))
			return err
		}
	} else {
		// 记录变更前的数量
		oldStock := inventory.Stock
		
		// 存在则更新
		inventory.Stock = stock
		inventory.UpdatedAt = time.Now()
		
		// 记录历史
		historyRecord := &entity.InventoryHistory{
			ProductID:   productID,
			WarehouseID: defaultWarehouseID,
			Quantity:    stock - oldStock, // 正数表示增加，负数表示减少
			Operation:   entity.OperationAdjust,
			Operator:    operator,
			Remark:      "手动设置库存",
			CreatedAt:   time.Now(),
		}
		
		if err := s.repo.RecordInventoryHistory(ctx, historyRecord); err != nil {
			s.logger.Warn("Failed to record inventory history",
				zap.Int64("product_id", productID),
				zap.Error(err))
			// 不影响主流程，继续执行
		}
	}
	
	// 保存或更新库存
	if err := s.repo.SetInventory(ctx, inventory); err != nil {
		s.logger.Error("Failed to set inventory",
			zap.Int64("product_id", productID),
			zap.Int("stock", stock),
			zap.Error(err))
		return ErrOperationFailed
	}
	
	return nil
}

// AddStock 增加库存
func (s *InventoryServiceImpl) AddStock(ctx context.Context, productID int64, quantity int, remark string) error {
	if productID <= 0 || quantity <= 0 {
		return ErrInvalidArgument
	}
	
	const defaultWarehouseID = 1
	
	err := s.repo.IncreaseStock(ctx, productID, defaultWarehouseID, quantity, remark)
	if err != nil {
		s.logger.Error("Failed to increase stock",
			zap.Int64("product_id", productID),
			zap.Int("quantity", quantity),
			zap.Error(err))
		return ErrOperationFailed
	}
	
	return nil
}

// AdjustStock 调整库存（库存盘点）
func (s *InventoryServiceImpl) AdjustStock(ctx context.Context, productID int64, newStock int, operator string, remark string) error {
	if productID <= 0 || newStock < 0 {
		return ErrInvalidArgument
	}
	
	const defaultWarehouseID = 1
	
	err := s.repo.AdjustStock(ctx, productID, defaultWarehouseID, newStock, operator, remark)
	if err != nil {
		s.logger.Error("Failed to adjust stock",
			zap.Int64("product_id", productID),
			zap.Int("new_stock", newStock),
			zap.Error(err))
		return ErrOperationFailed
	}
	
	return nil
}

// GetInventoryHistory 获取库存历史记录
func (s *InventoryServiceImpl) GetInventoryHistory(ctx context.Context, productID int64, page, pageSize int) ([]*entity.InventoryHistory, int64, error) {
	if productID <= 0 || page <= 0 || pageSize <= 0 {
		return nil, 0, ErrInvalidArgument
	}
	
	return s.repo.GetInventoryHistoryByProductID(ctx, productID, page, pageSize)
}
