package service

import (
	"context"
	"errors"
	"time"
	
	"go.uber.org/zap"
	
	"shop/backend/inventory/internal/domain/entity"
	"shop/backend/inventory/internal/domain/valueobject"
	"shop/backend/inventory/internal/repository"
)

// 定义错误
var (
	ErrLockFailed      = errors.New("lock inventory failed")
	ErrUnlockFailed    = errors.New("unlock inventory failed")
	ErrConfirmFailed   = errors.New("confirm inventory deduction failed")
	ErrLockNotFound    = errors.New("lock record not found")
	ErrInvalidLockKey  = errors.New("invalid lock key")
)

// InventoryLockServiceImpl 库存锁定服务实现
type InventoryLockServiceImpl struct {
	repo   repository.InventoryRepository
	logger *zap.Logger
}

// NewInventoryLockService 创建库存锁定服务实例
func NewInventoryLockService(repo repository.InventoryRepository, logger *zap.Logger) InventoryLockService {
	return &InventoryLockServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

// LockInventory 锁定库存
func (s *InventoryLockServiceImpl) LockInventory(ctx context.Context, lockKey string, items []LockItem, timeoutSeconds int) (*LockResult, error) {
	if lockKey == "" || len(items) == 0 {
		return nil, ErrInvalidArgument
	}
	
	// 转换为仓储层需要的格式
	stockOps := make([]*valueobject.StockOperation, 0, len(items))
	for _, item := range items {
		stockOps = append(stockOps, &valueobject.StockOperation{
			ProductID:   item.ProductID,
			WarehouseID: 1, // 默认仓库ID
			Quantity:    item.Quantity,
			OrderSN:     lockKey,
		})
	}
	
	// 调用仓储层锁定库存
	result, err := s.repo.LockStock(ctx, lockKey, stockOps)
	if err != nil {
		s.logger.Error("Failed to lock inventory",
			zap.String("lock_key", lockKey),
			zap.Any("items", items),
			zap.Error(err))
			
		// 将仓储层的失败项转换为服务层的失败项
		failItems := make([]*LockFailItem, 0)
		if result != nil && len(result.FailItems) > 0 {
			for _, item := range result.FailItems {
				failItems = append(failItems, &LockFailItem{
					ProductID: item.ProductID,
					Quantity:  item.Quantity,
					Available: item.Available,
					Reason:    item.Reason,
				})
			}
		}
		
		return &LockResult{
			Success:   false,
			Message:   err.Error(),
			FailItems: failItems,
		}, ErrLockFailed
	}
	
	// 成功锁定
	if result == nil {
		return &LockResult{
			Success: true,
			Message: "Lock inventory success",
		}, nil
	}
	
	// 将仓储层的结果转换为服务层的结果
	serviceResult := &LockResult{
		Success: result.Success,
		Message: result.Message,
	}
	
	if len(result.FailItems) > 0 {
		serviceResult.FailItems = make([]*LockFailItem, 0, len(result.FailItems))
		for _, item := range result.FailItems {
			serviceResult.FailItems = append(serviceResult.FailItems, &LockFailItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Available: item.Available,
				Reason:    item.Reason,
			})
		}
	}
	
	return serviceResult, nil
}

// UnlockInventory 解锁库存
func (s *InventoryLockServiceImpl) UnlockInventory(ctx context.Context, lockKey string) error {
	if lockKey == "" {
		return ErrInvalidLockKey
	}
	
	err := s.repo.UnlockStock(ctx, lockKey)
	if err != nil {
		s.logger.Error("Failed to unlock inventory",
			zap.String("lock_key", lockKey),
			zap.Error(err))
		return ErrUnlockFailed
	}
	
	return nil
}

// ConfirmReduce 确认扣减库存
func (s *InventoryLockServiceImpl) ConfirmReduce(ctx context.Context, lockKey string) error {
	if lockKey == "" {
		return ErrInvalidLockKey
	}
	
	err := s.repo.ReduceStock(ctx, lockKey)
	if err != nil {
		s.logger.Error("Failed to confirm reduce stock",
			zap.String("lock_key", lockKey),
			zap.Error(err))
		return ErrConfirmFailed
	}
	
	return nil
}

// GetLockDetail 获取锁定记录详情
func (s *InventoryLockServiceImpl) GetLockDetail(ctx context.Context, lockKey string) (*entity.StockSellDetail, error) {
	if lockKey == "" {
		return nil, ErrInvalidLockKey
	}
	
	detail, err := s.repo.GetStockSellDetail(ctx, lockKey)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrLockNotFound
		}
		s.logger.Error("Failed to get lock detail",
			zap.String("lock_key", lockKey),
			zap.Error(err))
		return nil, err
	}
	
	return detail, nil
}

// GetExpiredLocks 获取过期未处理的锁定记录
func (s *InventoryLockServiceImpl) GetExpiredLocks(ctx context.Context, before time.Time) ([]string, error) {
	// 这里需要扩展仓储接口添加获取过期锁定记录的方法
	// 简单实现，假设我们已经有这样的方法
	// return s.repo.GetExpiredLocks(ctx, before)
	
	s.logger.Warn("GetExpiredLocks method not implemented")
	return []string{}, nil
}
