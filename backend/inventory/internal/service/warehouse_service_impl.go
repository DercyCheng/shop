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
	ErrWarehouseNotFound    = errors.New("warehouse not found")
	ErrWarehouseNameExists  = errors.New("warehouse name already exists")
	ErrWarehouseCreateFailed = errors.New("failed to create warehouse")
	ErrWarehouseUpdateFailed = errors.New("failed to update warehouse")
	ErrWarehouseDeleteFailed = errors.New("failed to delete warehouse")
)

// WarehouseServiceImpl 仓库管理服务实现
type WarehouseServiceImpl struct {
	repo   repository.WarehouseRepository
	logger *zap.Logger
}

// NewWarehouseService 创建仓库服务实例
func NewWarehouseService(repo repository.WarehouseRepository, logger *zap.Logger) WarehouseService {
	return &WarehouseServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

// GetWarehouse 获取仓库信息
func (s *WarehouseServiceImpl) GetWarehouse(ctx context.Context, id int) (*entity.Warehouse, error) {
	if id <= 0 {
		return nil, ErrInvalidArgument
	}
	
	warehouse, err := s.repo.GetWarehouse(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrWarehouseNotFound
		}
		s.logger.Error("Failed to get warehouse",
			zap.Int("id", id),
			zap.Error(err))
		return nil, err
	}
	
	return warehouse, nil
}

// GetWarehouseByName 根据名称获取仓库
func (s *WarehouseServiceImpl) GetWarehouseByName(ctx context.Context, name string) (*entity.Warehouse, error) {
	if name == "" {
		return nil, ErrInvalidArgument
	}
	
	warehouse, err := s.repo.GetWarehouseByName(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrWarehouseNotFound
		}
		s.logger.Error("Failed to get warehouse by name",
			zap.String("name", name),
			zap.Error(err))
		return nil, err
	}
	
	return warehouse, nil
}

// ListWarehouses 获取仓库列表
func (s *WarehouseServiceImpl) ListWarehouses(ctx context.Context, page, pageSize int) ([]*entity.Warehouse, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, ErrInvalidArgument
	}
	
	// 查询仓库列表
	warehouses, total, err := s.repo.ListWarehouses(ctx, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to list warehouses",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err))
		return nil, 0, err
	}
	
	return warehouses, total, nil
}

// CreateWarehouse 创建仓库
func (s *WarehouseServiceImpl) CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	if warehouse == nil || warehouse.Name == "" || warehouse.Address == "" {
		return ErrInvalidArgument
	}
	
	// 检查同名仓库是否已存在
	existingWarehouse, err := s.repo.GetWarehouseByName(ctx, warehouse.Name)
	if err == nil && existingWarehouse != nil {
		return ErrWarehouseNameExists
	} else if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		s.logger.Error("Failed to check existing warehouse",
			zap.String("name", warehouse.Name),
			zap.Error(err))
		return err
	}
	
	// 设置创建时间和更新时间
	now := time.Now()
	warehouse.CreatedAt = now
	warehouse.UpdatedAt = now
	
	// 默认启用状态
	if warehouse.Status == 0 {
		warehouse.Status = 1
	}
	
	// 创建仓库
	if err := s.repo.CreateWarehouse(ctx, warehouse); err != nil {
		s.logger.Error("Failed to create warehouse",
			zap.String("name", warehouse.Name),
			zap.Error(err))
		return ErrWarehouseCreateFailed
	}
	
	return nil
}

// UpdateWarehouse 更新仓库
func (s *WarehouseServiceImpl) UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	if warehouse == nil || warehouse.ID <= 0 {
		return ErrInvalidArgument
	}
	
	// 检查仓库是否存在
	existingWarehouse, err := s.repo.GetWarehouse(ctx, warehouse.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return ErrWarehouseNotFound
		}
		s.logger.Error("Failed to check existing warehouse",
			zap.Int("id", warehouse.ID),
			zap.Error(err))
		return err
	}
	
	// 如果更改了名称，检查新名称是否冲突
	if warehouse.Name != existingWarehouse.Name {
		checkWarehouse, err := s.repo.GetWarehouseByName(ctx, warehouse.Name)
		if err == nil && checkWarehouse != nil && checkWarehouse.ID != warehouse.ID {
			return ErrWarehouseNameExists
		} else if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
			s.logger.Error("Failed to check existing warehouse by name",
				zap.String("name", warehouse.Name),
				zap.Error(err))
			return err
		}
	}
	
	// 更新时间
	warehouse.UpdatedAt = time.Now()
	
	// 保留创建时间
	warehouse.CreatedAt = existingWarehouse.CreatedAt
	
	// 更新仓库
	if err := s.repo.UpdateWarehouse(ctx, warehouse); err != nil {
		s.logger.Error("Failed to update warehouse",
			zap.Int("id", warehouse.ID),
			zap.Error(err))
		return ErrWarehouseUpdateFailed
	}
	
	return nil
}

// DeleteWarehouse 删除仓库
func (s *WarehouseServiceImpl) DeleteWarehouse(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidArgument
	}
	
	// 检查仓库是否存在
	_, err := s.repo.GetWarehouse(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return ErrWarehouseNotFound
		}
		s.logger.Error("Failed to check existing warehouse",
			zap.Int("id", id),
			zap.Error(err))
		return err
	}
	
	// 删除仓库（软删除）
	if err := s.repo.DeleteWarehouse(ctx, id); err != nil {
		s.logger.Error("Failed to delete warehouse",
			zap.Int("id", id),
			zap.Error(err))
		return ErrWarehouseDeleteFailed
	}
	
	return nil
}
