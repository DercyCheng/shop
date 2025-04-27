package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"shop/inventory/internal/domain/entity"
)

// WarehouseDAO provides data access methods for warehouse entity
type WarehouseDAO struct {
	db *gorm.DB
}

// NewWarehouseDAO creates a new WarehouseDAO
func NewWarehouseDAO(db *gorm.DB) *WarehouseDAO {
	return &WarehouseDAO{db: db}
}

// Create inserts a new warehouse
func (d *WarehouseDAO) Create(ctx context.Context, warehouse *entity.Warehouse) error {
	return d.db.WithContext(ctx).Create(warehouse).Error
}

// GetByID retrieves a warehouse by ID
func (d *WarehouseDAO) GetByID(ctx context.Context, id int64) (*entity.Warehouse, error) {
	var warehouse entity.Warehouse
	if err := d.db.WithContext(ctx).First(&warehouse, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("warehouse not found")
		}
		return nil, err
	}
	return &warehouse, nil
}

// List retrieves warehouses with pagination
func (d *WarehouseDAO) List(ctx context.Context, page, pageSize int, activeOnly bool) ([]*entity.Warehouse, int64, error) {
	var warehouses []*entity.Warehouse
	var total int64

	query := d.db.WithContext(ctx)

	// Apply active filter if requested
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	// Count total before applying pagination
	countQuery := query
	if err := countQuery.Model(&entity.Warehouse{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&warehouses).Error; err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}

// Update updates a warehouse
func (d *WarehouseDAO) Update(ctx context.Context, warehouse *entity.Warehouse) error {
	return d.db.WithContext(ctx).Save(warehouse).Error
}

// Delete deletes a warehouse
func (d *WarehouseDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Delete(&entity.Warehouse{}, id).Error
}

// SetActive sets the active status of a warehouse
func (d *WarehouseDAO) SetActive(ctx context.Context, id int64, active bool) error {
	return d.db.WithContext(ctx).Model(&entity.Warehouse{}).
		Where("id = ?", id).
		Update("is_active", active).
		Error
}
