package repository

import (
	"context"

	"shop/inventory/internal/domain/entity"
)

// WarehouseRepository defines the interface for warehouse data access
type WarehouseRepository interface {
	// Create creates a new warehouse record
	Create(ctx context.Context, warehouse *entity.Warehouse) error

	// GetByID retrieves a warehouse record by its ID
	GetByID(ctx context.Context, id int64) (*entity.Warehouse, error)

	// List retrieves warehouses with pagination
	List(ctx context.Context, page, pageSize int, activeOnly bool) ([]*entity.Warehouse, int64, error)

	// Update updates a warehouse record
	Update(ctx context.Context, warehouse *entity.Warehouse) error

	// Delete deletes a warehouse record
	Delete(ctx context.Context, id int64) error
}
