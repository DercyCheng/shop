package repository

import (
	"context"

	"shop/inventory/internal/domain/entity"
	"shop/inventory/internal/repository/dao"
)

// WarehouseRepositoryImpl implements the WarehouseRepository interface
type WarehouseRepositoryImpl struct {
	warehouseDAO *dao.WarehouseDAO
}

// NewWarehouseRepository creates a new WarehouseRepositoryImpl
func NewWarehouseRepository(warehouseDAO *dao.WarehouseDAO) WarehouseRepository {
	return &WarehouseRepositoryImpl{warehouseDAO: warehouseDAO}
}

// Create creates a new warehouse record
func (r *WarehouseRepositoryImpl) Create(ctx context.Context, warehouse *entity.Warehouse) error {
	return r.warehouseDAO.Create(ctx, warehouse)
}

// GetByID retrieves a warehouse record by its ID
func (r *WarehouseRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Warehouse, error) {
	return r.warehouseDAO.GetByID(ctx, id)
}

// List retrieves warehouses with pagination
func (r *WarehouseRepositoryImpl) List(ctx context.Context, page, pageSize int, activeOnly bool) ([]*entity.Warehouse, int64, error) {
	return r.warehouseDAO.List(ctx, page, pageSize, activeOnly)
}

// Update updates a warehouse record
func (r *WarehouseRepositoryImpl) Update(ctx context.Context, warehouse *entity.Warehouse) error {
	return r.warehouseDAO.Update(ctx, warehouse)
}

// Delete deletes a warehouse record
func (r *WarehouseRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.warehouseDAO.Delete(ctx, id)
}
