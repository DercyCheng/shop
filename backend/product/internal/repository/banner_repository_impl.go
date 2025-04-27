package repository

import (
	"context"

	"shop/product/internal/domain/entity"
	"shop/product/internal/repository/dao"
)

// BannerRepositoryImpl implements BannerRepository interface
type BannerRepositoryImpl struct {
	bannerDAO *dao.BannerDAO
}

// NewBannerRepository creates a new banner repository
func NewBannerRepository(bannerDAO *dao.BannerDAO) BannerRepository {
	return &BannerRepositoryImpl{
		bannerDAO: bannerDAO,
	}
}

// Create adds a new banner
func (r *BannerRepositoryImpl) Create(ctx context.Context, banner *entity.Banner) error {
	return r.bannerDAO.Create(ctx, banner)
}

// Update modifies an existing banner
func (r *BannerRepositoryImpl) Update(ctx context.Context, banner *entity.Banner) error {
	return r.bannerDAO.Update(ctx, banner)
}

// Delete removes a banner by ID
func (r *BannerRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.bannerDAO.Delete(ctx, id)
}

// GetByID retrieves a banner by ID
func (r *BannerRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Banner, error) {
	return r.bannerDAO.GetByID(ctx, id)
}

// List retrieves all banners ordered by index
func (r *BannerRepositoryImpl) List(ctx context.Context, limit int) ([]*entity.Banner, error) {
	return r.bannerDAO.List(ctx, limit)
}
