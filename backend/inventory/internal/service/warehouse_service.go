package service

import (
	"context"

	"shop/inventory/internal/domain/entity"
)

// WarehouseService defines methods for warehouse-related operations
type WarehouseService interface {
	// GetWarehouseByID retrieves a warehouse by ID
	GetWarehouseByID(ctx context.Context, id int64) (*entity.Warehouse, error)

	// ListWarehouses retrieves warehouses with pagination
	ListWarehouses(ctx context.Context, page, pageSize int, activeOnly bool) ([]*entity.Warehouse, int64, error)

	// CreateWarehouse creates a new warehouse
	CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error

	// UpdateWarehouse updates a warehouse
	UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error

	// DeleteWarehouse deletes a warehouse
	DeleteWarehouse(ctx context.Context, id int64) error

	// ActivateWarehouse activates a warehouse
	ActivateWarehouse(ctx context.Context, id int64) error

	// DeactivateWarehouse deactivates a warehouse
	DeactivateWarehouse(ctx context.Context, id int64) error
}
