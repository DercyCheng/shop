package service

import (
	"context"
	"errors"

	"shop/inventory/internal/domain/entity"
	"shop/inventory/internal/repository"
)

// WarehouseServiceImpl implements the WarehouseService interface
type WarehouseServiceImpl struct {
	warehouseRepo repository.WarehouseRepository
}

// NewWarehouseService creates a new WarehouseServiceImpl
func NewWarehouseService(warehouseRepo repository.WarehouseRepository) WarehouseService {
	return &WarehouseServiceImpl{warehouseRepo: warehouseRepo}
}

// GetWarehouseByID retrieves a warehouse by ID
func (s *WarehouseServiceImpl) GetWarehouseByID(ctx context.Context, id int64) (*entity.Warehouse, error) {
	if id <= 0 {
		return nil, errors.New("invalid warehouse ID")
	}
	return s.warehouseRepo.GetByID(ctx, id)
}

// ListWarehouses retrieves warehouses with pagination
func (s *WarehouseServiceImpl) ListWarehouses(ctx context.Context, page, pageSize int, activeOnly bool) ([]*entity.Warehouse, int64, error) {
	// Ensure minimum page and pageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	return s.warehouseRepo.List(ctx, page, pageSize, activeOnly)
}

// CreateWarehouse creates a new warehouse
func (s *WarehouseServiceImpl) CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	// Validate the warehouse data
	if warehouse.Name == "" {
		return errors.New("warehouse name cannot be empty")
	}

	return s.warehouseRepo.Create(ctx, warehouse)
}

// UpdateWarehouse updates a warehouse
func (s *WarehouseServiceImpl) UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error {
	// Validate the warehouse data
	if warehouse.ID <= 0 {
		return errors.New("invalid warehouse ID")
	}
	if warehouse.Name == "" {
		return errors.New("warehouse name cannot be empty")
	}

	// Check if warehouse exists
	existingWarehouse, err := s.warehouseRepo.GetByID(ctx, warehouse.ID)
	if err != nil {
		return err
	}
	if existingWarehouse == nil {
		return errors.New("warehouse not found")
	}

	return s.warehouseRepo.Update(ctx, warehouse)
}

// DeleteWarehouse deletes a warehouse
func (s *WarehouseServiceImpl) DeleteWarehouse(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid warehouse ID")
	}

	// Check if warehouse exists
	existingWarehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingWarehouse == nil {
		return errors.New("warehouse not found")
	}

	return s.warehouseRepo.Delete(ctx, id)
}

// ActivateWarehouse activates a warehouse
func (s *WarehouseServiceImpl) ActivateWarehouse(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid warehouse ID")
	}

	// Check if warehouse exists
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if warehouse == nil {
		return errors.New("warehouse not found")
	}

	// Update if not already active
	if !warehouse.IsActive {
		warehouse.IsActive = true
		return s.warehouseRepo.Update(ctx, warehouse)
	}

	return nil
}

// DeactivateWarehouse deactivates a warehouse
func (s *WarehouseServiceImpl) DeactivateWarehouse(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid warehouse ID")
	}

	// Check if warehouse exists
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if warehouse == nil {
		return errors.New("warehouse not found")
	}

	// Update if currently active
	if warehouse.IsActive {
		warehouse.IsActive = false
		return s.warehouseRepo.Update(ctx, warehouse)
	}

	return nil
}
