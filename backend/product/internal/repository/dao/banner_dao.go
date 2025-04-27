package dao

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"shop/product/internal/domain/entity"
)

// BannerDAO handles data access operations for banners
type BannerDAO struct {
	db *gorm.DB
}

// NewBannerDAO creates a new banner data access object
func NewBannerDAO(db *gorm.DB) *BannerDAO {
	return &BannerDAO{db: db}
}

// Create adds a new banner to the database
func (d *BannerDAO) Create(ctx context.Context, banner *entity.Banner) error {
	return d.db.WithContext(ctx).Create(banner).Error
}

// Update modifies an existing banner
func (d *BannerDAO) Update(ctx context.Context, banner *entity.Banner) error {
	return d.db.WithContext(ctx).Updates(banner).Error
}

// Delete soft-deletes a banner by ID
func (d *BannerDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Delete(&entity.Banner{}, id).Error
}

// GetByID retrieves a banner by ID
func (d *BannerDAO) GetByID(ctx context.Context, id int64) (*entity.Banner, error) {
	var banner entity.Banner
	result := d.db.WithContext(ctx).First(&banner, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("banner not found: %v", result.Error)
		}
		return nil, result.Error
	}
	return &banner, nil
}

// List retrieves all banners ordered by index
func (d *BannerDAO) List(ctx context.Context, limit int) ([]*entity.Banner, error) {
	var banners []*entity.Banner
	query := d.db.WithContext(ctx).Order("index")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&banners).Error; err != nil {
		return nil, err
	}

	return banners, nil
}
