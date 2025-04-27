package service

import (
	"context"

	"shop/product/internal/domain/entity"
)

// BannerService defines the interface for banner business logic
type BannerService interface {
	// GetBannerByID retrieves a banner by ID
	GetBannerByID(ctx context.Context, id int64) (*entity.Banner, error)

	// ListBanners retrieves all banners ordered by index
	ListBanners(ctx context.Context, limit int) ([]*entity.Banner, error)

	// CreateBanner adds a new banner
	CreateBanner(ctx context.Context, banner *entity.Banner) (*entity.Banner, error)

	// UpdateBanner modifies an existing banner
	UpdateBanner(ctx context.Context, banner *entity.Banner) error

	// DeleteBanner removes a banner by ID
	DeleteBanner(ctx context.Context, id int64) error
}
