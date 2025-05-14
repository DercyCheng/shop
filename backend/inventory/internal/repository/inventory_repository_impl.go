package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	
	"go.uber.org/zap"
	"gorm.io/gorm"
	
	"shop/backend/inventory/internal/domain/entity"
	"shop/backend/inventory/internal/domain/valueobject"
	"shop/backend/inventory/internal/repository/cache"
)

var (
	// ErrInsufficientStock 库存不足
	ErrInsufficientStock = errors.New("insufficient stock")
	// ErrStockLocked 库存已锁定
	ErrStockLocked = errors.New("stock already locked")
	// ErrStockNotLocked 库存未锁定
	ErrStockNotLocked = errors.New("stock not locked")
	// ErrRecordNotFound 记录不存在
	ErrRecordNotFound = errors.New("record not found")
)

// InventoryRepositoryImpl 库存仓储实现
type InventoryRepositoryImpl struct {
	db    *gorm.DB
	cache cache.InventoryCache
	logger *zap.Logger
}

// NewInventoryRepository 创建库存仓储
func NewInventoryRepository(db *gorm.DB, cache cache.InventoryCache, logger *zap.Logger) InventoryRepository {
	return &InventoryRepositoryImpl{
		db:    db,
		cache: cache,
		logger: logger,
	}
}

// GetInventory 获取商品库存
func (r *InventoryRepositoryImpl) GetInventory(ctx context.Context, productID int64, warehouseID int) (*entity.Inventory, error) {
	// 先尝试从缓存获取
	inventory, err := r.cache.GetInventory(ctx, productID, warehouseID)
	if err == nil {
		return inventory, nil
	}
	
	// 缓存未命中，从数据库查询
	inventory = &entity.Inventory{}
	err = r.db.WithContext(ctx).Where("goods = ? AND warehouse_id = ?", productID, warehouseID).First(inventory).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		r.logger.Error("Failed to get inventory from database", 
			zap.Int64("product_id", productID), 
			zap.Int("warehouse_id", warehouseID), 
			zap.Error(err))
		return nil, err
	}
	
	// 更新缓存
	if err := r.cache.SetInventory(ctx, inventory); err != nil {
		r.logger.Warn("Failed to set inventory cache", 
			zap.Int64("product_id", productID), 
			zap.Int("warehouse_id", warehouseID), 
			zap.Error(err))
	}
	
	return inventory, nil
}

// BatchGetInventory 批量获取商品库存
func (r *InventoryRepositoryImpl) BatchGetInventory(ctx context.Context, productIDs []int64, warehouseID int) ([]*entity.Inventory, error) {
	if len(productIDs) == 0 {
		return []*entity.Inventory{}, nil
	}
	
	// 尝试从缓存批量获取
	inventories, missingIDs := r.cache.BatchGetInventory(ctx, productIDs, warehouseID)
	
	// 如果所有ID都在缓存中找到，则直接返回
	if len(missingIDs) == 0 {
		return inventories, nil
	}
	
	// 查询缓存未命中的记录
	var dbInventories []*entity.Inventory
	err := r.db.WithContext(ctx).Where("goods IN ? AND warehouse_id = ?", missingIDs, warehouseID).Find(&dbInventories).Error
	if err != nil {
		r.logger.Error("Failed to batch get inventory from database", 
			zap.Any("product_ids", missingIDs), 
			zap.Int("warehouse_id", warehouseID), 
			zap.Error(err))
		return inventories, err
	}
	
	// 更新缓存并合并结果
	for _, inv := range dbInventories {
		if err := r.cache.SetInventory(ctx, inv); err != nil {
			r.logger.Warn("Failed to set inventory cache", 
				zap.Int64("product_id", inv.ProductID), 
				zap.Int("warehouse_id", inv.WarehouseID), 
				zap.Error(err))
		}
		inventories = append(inventories, inv)
	}
	
	return inventories, nil
}

// SetInventory 设置商品库存
func (r *InventoryRepositoryImpl) SetInventory(ctx context.Context, inventory *entity.Inventory) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing entity.Inventory
		
		// 检查记录是否存在
		err := tx.Where("goods = ? AND warehouse_id = ?", inventory.ProductID, inventory.WarehouseID).First(&existing).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 不存在则创建
				inventory.CreatedAt = time.Now()
				inventory.UpdatedAt = time.Now()
				if err := tx.Create(inventory).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// 存在则更新
			updates := map[string]interface{}{
				"stocks":           inventory.Stock,
				"lock_stocks":      inventory.LockStock,
				"alert_threshold":  inventory.AlertThreshold,
				"version":          existing.Version + 1,
				"updated_at":       time.Now(),
			}
			
			if err := tx.Model(&entity.Inventory{}).
				Where("id = ? AND version = ?", existing.ID, existing.Version).
				Updates(updates).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
	
	if err != nil {
		r.logger.Error("Failed to set inventory", 
			zap.Int64("product_id", inventory.ProductID), 
			zap.Int("warehouse_id", inventory.WarehouseID), 
			zap.Error(err))
		return err
	}
	
	// 更新缓存
	if err := r.cache.DeleteInventory(ctx, inventory.ProductID, inventory.WarehouseID); err != nil {
		r.logger.Warn("Failed to delete inventory cache", 
			zap.Int64("product_id", inventory.ProductID), 
			zap.Int("warehouse_id", inventory.WarehouseID), 
			zap.Error(err))
	}
	
	return nil
}

// LockStock 锁定库存
func (r *InventoryRepositoryImpl) LockStock(ctx context.Context, orderSN string, items []*valueobject.StockOperation) (*valueobject.LockResult, error) {
	if len(items) == 0 {
		return &valueobject.LockResult{
			Success: true,
			Message: "No items to lock",
		}, nil
	}
	
	// 检查订单是否已经锁定过库存
	existingDetail, err := r.GetStockSellDetail(ctx, orderSN)
	if err == nil && existingDetail != nil {
		// 已经处理过的请求，返回之前的结果
		return &valueobject.LockResult{
			Success: true,
			Message: "Stock already locked for this order",
		}, nil
	}
	
	// 开启事务处理
	result := &valueobject.LockResult{
		Success: true,
		FailItems: []*valueobject.LockFailItem{},
	}
	
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 处理每个商品的库存锁定
		detailItems := make([]*entity.StockDetail, 0, len(items))
		
		for _, item := range items {
			// 使用乐观锁更新库存
			var inv entity.Inventory
			
			// 获取当前库存
			if err := tx.Where("goods = ? AND warehouse_id = ?", item.ProductID, item.WarehouseID).First(&inv).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 库存不存在，添加到失败项
					result.FailItems = append(result.FailItems, &valueobject.LockFailItem{
						ProductID: item.ProductID,
						Quantity:  item.Quantity,
						Available: 0,
						Reason:    "Inventory not found",
					})
					continue
				}
				return err
			}
			
			// 检查库存是否足够
			if inv.Stock - inv.LockStock < item.Quantity {
				// 库存不足，添加到失败项
				result.FailItems = append(result.FailItems, &valueobject.LockFailItem{
					ProductID: item.ProductID,
					Quantity:  item.Quantity,
					Available: inv.Stock - inv.LockStock,
					Reason:    "Insufficient stock",
				})
				continue
			}
			
			// 更新锁定库存
			if err := tx.Model(&entity.Inventory{}).
				Where("id = ? AND version = ? AND stocks - lock_stocks >= ?", inv.ID, inv.Version, item.Quantity).
				Updates(map[string]interface{}{
					"lock_stocks": gorm.Expr("lock_stocks + ?", item.Quantity),
					"version":     inv.Version + 1,
					"updated_at":  time.Now(),
				}).Error; err != nil {
				return err
			}
			
			// 添加到详情列表
			detailItems = append(detailItems, &entity.StockDetail{
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				WarehouseID: item.WarehouseID,
			})
			
			// 记录库存历史
			history := &entity.InventoryHistory{
				ProductID:   item.ProductID,
				WarehouseID: item.WarehouseID,
				Quantity:    item.Quantity,
				Operation:   entity.OperationLock,
				OrderSN:     orderSN,
				Operator:    item.Operator,
				Remark:      item.Remark,
				CreatedAt:   time.Now(),
			}
			
			if err := tx.Create(history).Error; err != nil {
				r.logger.Error("Failed to record inventory history", 
					zap.Error(err),
					zap.String("order_sn", orderSN),
					zap.Int64("product_id", item.ProductID))
				// 不中断主流程
			}
			
			// 更新缓存
			if err := r.cache.DeleteInventory(ctx, item.ProductID, item.WarehouseID); err != nil {
				r.logger.Warn("Failed to delete inventory cache", 
					zap.Int64("product_id", item.ProductID), 
					zap.Int("warehouse_id", item.WarehouseID), 
					zap.Error(err))
			}
		}
		
		// 如果有失败项且要求全部成功，则回滚事务
		if len(result.FailItems) > 0 {
			result.Success = false
			result.Message = fmt.Sprintf("%d items failed to lock", len(result.FailItems))
			return ErrInsufficientStock
		}
		
		// 创建库存锁定记录
		now := time.Now()
		sellDetail := &entity.StockSellDetail{
			OrderSN:     orderSN,
			Status:      entity.StockLocked,
			DetailItems: detailItems,
			LockTime:    &now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		
		if err := tx.Create(sellDetail).Error; err != nil {
			return err
		}
		
		result.Message = "Stock locked successfully"
		return nil
	})
	
	if err != nil {
		if !errors.Is(err, ErrInsufficientStock) {
			r.logger.Error("Failed to lock stock", 
				zap.Error(err),
				zap.String("order_sn", orderSN))
			return nil, err
		}
		// ErrInsufficientStock 已经被处理，返回result
	}
	
	return result, nil
}

// UnlockStock 解锁库存
func (r *InventoryRepositoryImpl) UnlockStock(ctx context.Context, orderSN string) error {
	// 获取库存锁定详情
	detail, err := r.GetStockSellDetail(ctx, orderSN)
	if err != nil {
		return err
	}
	
	// 检查状态是否为已锁定
	if detail.Status != entity.StockLocked {
		r.logger.Warn("Cannot unlock stock with invalid status",
			zap.String("order_sn", orderSN),
			zap.Any("current_status", detail.Status))
		return ErrStockNotLocked
	}
	
	// 开启事务
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 解锁每个商品的库存
		for _, item := range detail.DetailItems {
			// 更新库存记录
			if err := tx.Model(&entity.Inventory{}).
				Where("goods = ? AND warehouse_id = ?", item.ProductID, item.WarehouseID).
				Updates(map[string]interface{}{
					"lock_stocks": gorm.Expr("lock_stocks - ?", item.Quantity),
					"version":     gorm.Expr("version + 1"),
					"updated_at":  time.Now(),
				}).Error; err != nil {
				return err
			}
			
			// 记录库存历史
			history := &entity.InventoryHistory{
				ProductID:   item.ProductID,
				WarehouseID: item.WarehouseID,
				Quantity:    -item.Quantity, // 负数表示解锁
				Operation:   entity.OperationUnlock,
				OrderSN:     orderSN,
				Remark:      "Order cancelled or timeout",
				CreatedAt:   time.Now(),
			}
			
			if err := tx.Create(history).Error; err != nil {
				r.logger.Error("Failed to record inventory history", 
					zap.Error(err),
					zap.String("order_sn", orderSN),
					zap.Int64("product_id", item.ProductID))
				// 不中断主流程
			}
			
			// 更新缓存
			if err := r.cache.DeleteInventory(ctx, item.ProductID, item.WarehouseID); err != nil {
				r.logger.Warn("Failed to delete inventory cache", 
					zap.Int64("product_id", item.ProductID), 
					zap.Int("warehouse_id", item.WarehouseID), 
					zap.Error(err))
			}
		}
		
		// 更新库存锁定记录状态
		return r.UpdateStockSellDetailStatus(ctx, orderSN, entity.StockReturned)
	})
	
	if err != nil {
		r.logger.Error("Failed to unlock stock", 
			zap.Error(err),
			zap.String("order_sn", orderSN))
		return err
	}
	
	return nil
}

// ReduceStock 扣减库存（确认扣减）
func (r *InventoryRepositoryImpl) ReduceStock(ctx context.Context, orderSN string) error {
	// 获取库存锁定详情
	detail, err := r.GetStockSellDetail(ctx, orderSN)
	if err != nil {
		return err
	}
	
	// 检查状态是否为已锁定
	if detail.Status != entity.StockLocked {
		r.logger.Warn("Cannot reduce stock with invalid status",
			zap.String("order_sn", orderSN),
			zap.Any("current_status", detail.Status))
		return ErrStockNotLocked
	}
	
	// 开启事务
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 扣减每个商品的库存
		for _, item := range detail.DetailItems {
			// 更新库存记录
			if err := tx.Model(&entity.Inventory{}).
				Where("goods = ? AND warehouse_id = ?", item.ProductID, item.WarehouseID).
				Updates(map[string]interface{}{
					"stocks":      gorm.Expr("stocks - ?", item.Quantity),
					"lock_stocks": gorm.Expr("lock_stocks - ?", item.Quantity),
					"version":     gorm.Expr("version + 1"),
					"updated_at":  time.Now(),
				}).Error; err != nil {
				return err
			}
			
			// 记录库存历史
			history := &entity.InventoryHistory{
				ProductID:   item.ProductID,
				WarehouseID: item.WarehouseID,
				Quantity:    -item.Quantity, // 负数表示扣减
				Operation:   entity.OperationDecrease,
				OrderSN:     orderSN,
				Remark:      "Order confirmed",
				CreatedAt:   time.Now(),
			}
			
			if err := tx.Create(history).Error; err != nil {
				r.logger.Error("Failed to record inventory history", 
					zap.Error(err),
					zap.String("order_sn", orderSN),
					zap.Int64("product_id", item.ProductID))
				// 不中断主流程
			}
			
			// 更新缓存
			if err := r.cache.DeleteInventory(ctx, item.ProductID, item.WarehouseID); err != nil {
				r.logger.Warn("Failed to delete inventory cache", 
					zap.Int64("product_id", item.ProductID), 
					zap.Int("warehouse_id", item.WarehouseID), 
					zap.Error(err))
			}
		}
		
		// 更新确认时间
		now := time.Now()
		if err := tx.Model(&entity.StockSellDetail{}).
			Where("order_sn = ?", orderSN).
			Updates(map[string]interface{}{
				"status":       entity.StockReduced,
				"confirm_time": now,
				"updated_at":   now,
			}).Error; err != nil {
			return err
		}
		
		return nil
	})
	
	if err != nil {
		r.logger.Error("Failed to reduce stock", 
			zap.Error(err),
			zap.String("order_sn", orderSN))
		return err
	}
	
	return nil
}

// IncreaseStock 增加库存
func (r *InventoryRepositoryImpl) IncreaseStock(ctx context.Context, productID int64, warehouseID int, quantity int, remark string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新库存
		var inv entity.Inventory
		err := tx.Where("goods = ? AND warehouse_id = ?", productID, warehouseID).First(&inv).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 创建新的库存记录
				now := time.Now()
				inv = entity.Inventory{
					ProductID:      productID,
					WarehouseID:    warehouseID,
					Stock:          quantity,
					Version:        0,
					LockStock:      0,
					AlertThreshold: 10, // 默认预警阈值
					CreatedAt:      now,
					UpdatedAt:      now,
				}
				return tx.Create(&inv).Error
			}
			return err
		}
		
		// 增加库存
		if err := tx.Model(&entity.Inventory{}).
			Where("id = ? AND version = ?", inv.ID, inv.Version).
			Updates(map[string]interface{}{
				"stocks":     gorm.Expr("stocks + ?", quantity),
				"version":    gorm.Expr("version + 1"),
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		
		// 记录库存历史
		history := &entity.InventoryHistory{
			ProductID:   productID,
			WarehouseID: warehouseID,
			Quantity:    quantity,
			Operation:   entity.OperationIncrease,
			Remark:      remark,
			CreatedAt:   time.Now(),
		}
		
		if err := tx.Create(history).Error; err != nil {
			r.logger.Error("Failed to record inventory history", 
				zap.Error(err),
				zap.Int64("product_id", productID))
			// 不中断主流程
		}
		
		return nil
	})
	
	if err != nil {
		r.logger.Error("Failed to increase stock", 
			zap.Error(err),
			zap.Int64("product_id", productID),
			zap.Int("warehouse_id", warehouseID))
		return err
	}
	
	// 更新缓存
	if err := r.cache.DeleteInventory(ctx, productID, warehouseID); err != nil {
		r.logger.Warn("Failed to delete inventory cache", 
			zap.Int64("product_id", productID), 
			zap.Int("warehouse_id", warehouseID), 
			zap.Error(err))
	}
	
	return nil
}

// DecreaseStock 减少库存（非锁定方式直接减少）
func (r *InventoryRepositoryImpl) DecreaseStock(ctx context.Context, productID int64, warehouseID int, quantity int, remark string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查询库存
		var inv entity.Inventory
		if err := tx.Where("goods = ? AND warehouse_id = ?", productID, warehouseID).First(&inv).Error; err != nil {
			return err
		}
		
		// 检查库存是否足够
		if inv.Stock < quantity {
			return ErrInsufficientStock
		}
		
		// 减少库存
		if err := tx.Model(&entity.Inventory{}).
			Where("id = ? AND version = ? AND stocks >= ?", inv.ID, inv.Version, quantity).
			Updates(map[string]interface{}{
				"stocks":     gorm.Expr("stocks - ?", quantity),
				"version":    gorm.Expr("version + 1"),
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		
		// 记录库存历史
		history := &entity.InventoryHistory{
			ProductID:   productID,
			WarehouseID: warehouseID,
			Quantity:    -quantity, // 负数表示减少
			Operation:   entity.OperationDecrease,
			Remark:      remark,
			CreatedAt:   time.Now(),
		}
		
		if err := tx.Create(history).Error; err != nil {
			r.logger.Error("Failed to record inventory history", 
				zap.Error(err),
				zap.Int64("product_id", productID))
			// 不中断主流程
		}
		
		return nil
	})
	
	if err != nil {
		r.logger.Error("Failed to decrease stock", 
			zap.Error(err),
			zap.Int64("product_id", productID),
			zap.Int("warehouse_id", warehouseID))
		return err
	}
	
	// 更新缓存
	if err := r.cache.DeleteInventory(ctx, productID, warehouseID); err != nil {
		r.logger.Warn("Failed to delete inventory cache", 
			zap.Int64("product_id", productID), 
			zap.Int("warehouse_id", warehouseID), 
			zap.Error(err))
	}
	
	return nil
}

// AdjustStock 调整库存（库存盘点）
func (r *InventoryRepositoryImpl) AdjustStock(ctx context.Context, productID int64, warehouseID int, newStock int, operator string, remark string) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查询库存
		var inv entity.Inventory
		err := tx.Where("goods = ? AND warehouse_id = ?", productID, warehouseID).First(&inv).Error
		
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 创建新的库存记录
				now := time.Now()
				inv = entity.Inventory{
					ProductID:      productID,
					WarehouseID:    warehouseID,
					Stock:          newStock,
					Version:        0,
					LockStock:      0,
					AlertThreshold: 10, // 默认预警阈值
					CreatedAt:      now,
					UpdatedAt:      now,
				}
				return tx.Create(&inv).Error
			}
			return err
		}
		
		// 记录调整前的库存
		oldStock := inv.Stock
		
		// 调整库存
		if err := tx.Model(&entity.Inventory{}).
			Where("id = ? AND version = ?", inv.ID, inv.Version).
			Updates(map[string]interface{}{
				"stocks":     newStock,
				"version":    gorm.Expr("version + 1"),
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		
		// 记录库存历史
		history := &entity.InventoryHistory{
			ProductID:   productID,
			WarehouseID: warehouseID,
			Quantity:    newStock - oldStock, // 可能为正数或负数
			Operation:   entity.OperationAdjust,
			Operator:    operator,
			Remark:      remark,
			CreatedAt:   time.Now(),
		}
		
		if err := tx.Create(history).Error; err != nil {
			r.logger.Error("Failed to record inventory history", 
				zap.Error(err),
				zap.Int64("product_id", productID))
			// 不中断主流程
		}
		
		return nil
	})
	
	if err != nil {
		r.logger.Error("Failed to adjust stock", 
			zap.Error(err),
			zap.Int64("product_id", productID),
			zap.Int("warehouse_id", warehouseID))
		return err
	}
	
	// 更新缓存
	if err := r.cache.DeleteInventory(ctx, productID, warehouseID); err != nil {
		r.logger.Warn("Failed to delete inventory cache", 
			zap.Int64("product_id", productID), 
			zap.Int("warehouse_id", warehouseID), 
			zap.Error(err))
	}
	
	return nil
}

// GetStockSellDetail 获取库存锁定记录
func (r *InventoryRepositoryImpl) GetStockSellDetail(ctx context.Context, orderSN string) (*entity.StockSellDetail, error) {
	var detail entity.StockSellDetail
	err := r.db.WithContext(ctx).Where("order_sn = ?", orderSN).First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		r.logger.Error("Failed to get stock sell detail", 
			zap.Error(err),
			zap.String("order_sn", orderSN))
		return nil, err
	}
	
	return &detail, nil
}

// UpdateStockSellDetailStatus 更新库存锁定记录状态
func (r *InventoryRepositoryImpl) UpdateStockSellDetailStatus(ctx context.Context, orderSN string, status entity.StockStatus) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}
	
	// 如果是确认状态，更新确认时间
	if status == entity.StockReduced {
		updates["confirm_time"] = now
	}
	
	err := r.db.WithContext(ctx).Model(&entity.StockSellDetail{}).
		Where("order_sn = ?", orderSN).
		Updates(updates).Error
		
	if err != nil {
		r.logger.Error("Failed to update stock sell detail status", 
			zap.Error(err),
			zap.String("order_sn", orderSN),
			zap.Any("status", status))
		return err
	}
	
	return nil
}

// CreateStockSellDetail 创建库存锁定记录
func (r *InventoryRepositoryImpl) CreateStockSellDetail(ctx context.Context, detail *entity.StockSellDetail) error {
	if detail.CreatedAt.IsZero() {
		detail.CreatedAt = time.Now()
	}
	if detail.UpdatedAt.IsZero() {
		detail.UpdatedAt = time.Now()
	}
	
	err := r.db.WithContext(ctx).Create(detail).Error
	if err != nil {
		r.logger.Error("Failed to create stock sell detail", 
			zap.Error(err),
			zap.String("order_sn", detail.OrderSN))
		return err
	}
	
	return nil
}

// RecordInventoryHistory 记录库存变更历史
func (r *InventoryRepositoryImpl) RecordInventoryHistory(ctx context.Context, history *entity.InventoryHistory) error {
	if history.CreatedAt.IsZero() {
		history.CreatedAt = time.Now()
	}
	
	err := r.db.WithContext(ctx).Create(history).Error
	if err != nil {
		r.logger.Error("Failed to record inventory history", 
			zap.Error(err),
			zap.Int64("product_id", history.ProductID),
			zap.String("operation", string(history.Operation)))
		return err
	}
	
	return nil
}

// GetInventoryHistoryByProductID 获取商品库存变更历史
func (r *InventoryRepositoryImpl) GetInventoryHistoryByProductID(ctx context.Context, productID int64, page, pageSize int) ([]*entity.InventoryHistory, int64, error) {
	var histories []*entity.InventoryHistory
	var total int64
	
	// 计算总数
	err := r.db.WithContext(ctx).Model(&entity.InventoryHistory{}).
		Where("goods = ?", productID).
		Count(&total).Error
	if err != nil {
		r.logger.Error("Failed to count inventory history", 
			zap.Error(err),
			zap.Int64("product_id", productID))
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err = r.db.WithContext(ctx).
		Where("goods = ?", productID).
		Order("id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&histories).Error
	if err != nil {
		r.logger.Error("Failed to get inventory history", 
			zap.Error(err),
			zap.Int64("product_id", productID))
		return nil, 0, err
	}
	
	return histories, total, nil
}
