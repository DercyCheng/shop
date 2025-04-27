package dao

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"shop/product/internal/domain/entity"
)

// BrandDAO handles data access operations for brands
type BrandDAO struct {
	db *gorm.DB
}

// NewBrandDAO creates a new brand data access object
func NewBrandDAO(db *gorm.DB) *BrandDAO {
	return &BrandDAO{db: db}
}

// Create adds a new brand to the database
func (d *BrandDAO) Create(ctx context.Context, brand *entity.Brand) error {
	return d.db.WithContext(ctx).Create(brand).Error
}

// Update modifies an existing brand
func (d *BrandDAO) Update(ctx context.Context, brand *entity.Brand) error {
	return d.db.WithContext(ctx).Updates(brand).Error
}

// Delete soft-deletes a brand by ID
func (d *BrandDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Delete(&entity.Brand{}, id).Error
}

// GetByID retrieves a brand by ID
func (d *BrandDAO) GetByID(ctx context.Context, id int64) (*entity.Brand, error) {
	var brand entity.Brand
	result := d.db.WithContext(ctx).First(&brand, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("brand not found: %v", result.Error)
		}
		return nil, result.Error
	}
	return &brand, nil
}

// List retrieves brands with pagination
func (d *BrandDAO) List(ctx context.Context, page, pageSize int) ([]*entity.Brand, int64, error) {
	var brands []*entity.Brand
	var total int64
	offset := (page - 1) * pageSize

	// Get total count
	if err := d.db.WithContext(ctx).Model(&entity.Brand{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := d.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&brands).Error; err != nil {
		return nil, 0, err
	}

	return brands, total, nil
}

// ListAll retrieves all brands
func (d *BrandDAO) ListAll(ctx context.Context) ([]*entity.Brand, error) {
	var brands []*entity.Brand
	if err := d.db.WithContext(ctx).Find(&brands).Error; err != nil {
		return nil, err
	}
	return brands, nil
}

// ListByIDs retrieves brands by a slice of IDs
func (d *BrandDAO) ListByIDs(ctx context.Context, ids []int64) ([]*entity.Brand, error) {
	var brands []*entity.Brand
	if err := d.db.WithContext(ctx).Where("id IN ?", ids).Find(&brands).Error; err != nil {
		return nil, err
	}
	return brands, nil
}
