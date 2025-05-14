package repository

import (
	"context"
	"errors"
	"time"
	
	"go.uber.org/zap"
	"gorm.io/gorm"
	
	"shop/backend/inventory/internal/domain/entity"
)

// WarehouseRepositoryImpl 仓库仓储实现
type WarehouseRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewWarehouseRepository 创建仓库仓储
func NewWarehouseRepository(db *gorm.DB, logger *zap.Logger) WarehouseRepository {
	return &WarehouseRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// GetWarehouse 获取仓库信息
func (r *WarehouseRepositoryImpl) GetWarehouse(ctx context.Context, id int) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&warehouse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		r.logger.Error("Failed to get warehouse", 
			zap.Int("id", id), 
			zap.Error(err))
		return nil, err
	}
	
	return &warehouse, nil
}

// GetWarehouseByName 根据名称获取仓库信息
func (r *WarehouseRepositoryImpl) GetWarehouseByName(ctx context.Context, name string) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&warehouse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		r.logger.Error("Failed to get warehouse by name", 
			zap.String("name", name), 
			zap.Error(err))
		return nil, err
	}
	
	return &warehouse, nil
}

// ListWarehouses 获取仓库列表
func (r *WarehouseRepositoryImpl) ListWarehouses(ctx context.Context, page, pageSize int) ([]*entity.Warehouse, int64, error) {
	var warehouses []*entity.Warehouse
	var total int64
	
	// 计算总数
	err := r.db.WithContext(ctx).Model(&entity.Warehouse{}).Count(&total).Error
	if err != nil {
		r.logger.Error("Failed to count warehouses", zap.Error(err))
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * pageSize
	err = r.db.WithContext(ctx).
		Order("id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&warehouses).Error
	if err != nil {
		r.logger.Error("Failed to list warehouses", zap.Error(err))
		return nil, 0, err
	}
	
	return warehouses, total, nil
}

// CreateWarehouse 创建仓库
func (r *WarehouseRepositoryImpl) CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	// 设置时间字段
	now := time.Now()
	warehouse.CreatedAt = now
	warehouse.UpdatedAt = now
	
	err := r.db.WithContext(ctx).Create(warehouse).Error
	if err != nil {
		r.logger.Error("Failed to create warehouse", 
			zap.String("name", warehouse.Name), 
			zap.Error(err))
		return err
	}
	
	return nil
}

// UpdateWarehouse 更新仓库
func (r *WarehouseRepositoryImpl) UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	// 设置更新时间
	warehouse.UpdatedAt = time.Now()
	
	err := r.db.WithContext(ctx).Model(&entity.Warehouse{}).
		Where("id = ?", warehouse.ID).
		Updates(map[string]interface{}{
			"name":       warehouse.Name,
			"address":    warehouse.Address,
			"contact":    warehouse.Contact,
			"phone":      warehouse.Phone,
			"status":     warehouse.Status,
			"updated_at": warehouse.UpdatedAt,
		}).Error
	if err != nil {
		r.logger.Error("Failed to update warehouse", 
			zap.Int("id", warehouse.ID), 
			zap.Error(err))
		return err
	}
	
	return nil
}

// DeleteWarehouse 删除仓库
func (r *WarehouseRepositoryImpl) DeleteWarehouse(ctx context.Context, id int) error {
	// 使用软删除
	err := r.db.WithContext(ctx).Delete(&entity.Warehouse{}, id).Error
	if err != nil {
		r.logger.Error("Failed to delete warehouse", 
			zap.Int("id", id), 
			zap.Error(err))
		return err
	}
	
	return nil
}
