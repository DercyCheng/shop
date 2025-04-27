package repository

import (
	"context"

	"shop/product/internal/domain/entity"
)

// BannerRepository defines the interface for banner data operations
type BannerRepository interface {
	// Create adds a new banner
	Create(ctx context.Context, banner *entity.Banner) error

	// Update modifies an existing banner
	Update(ctx context.Context, banner *entity.Banner) error

	// Delete removes a banner by ID
	Delete(ctx context.Context, id int64) error

	// GetByID retrieves a banner by ID
	GetByID(ctx context.Context, id int64) (*entity.Banner, error)

	// List retrieves all banners ordered by index
	List(ctx context.Context, limit int) ([]*entity.Banner, error)
}
