package service

import (
	"context"

	"shop/product/internal/domain/entity"
)

// BrandService defines the interface for brand business logic
type BrandService interface {
	// GetBrandByID retrieves a brand by ID
	GetBrandByID(ctx context.Context, id int64) (*entity.Brand, error)

	// ListBrands retrieves brands with pagination
	ListBrands(ctx context.Context, page, pageSize int) ([]*entity.Brand, int64, error)

	// GetAllBrands retrieves all brands
	GetAllBrands(ctx context.Context) ([]*entity.Brand, error)

	// CreateBrand adds a new brand
	CreateBrand(ctx context.Context, brand *entity.Brand) (*entity.Brand, error)

	// UpdateBrand modifies an existing brand
	UpdateBrand(ctx context.Context, brand *entity.Brand) error

	// DeleteBrand removes a brand by ID
	DeleteBrand(ctx context.Context, id int64) error
}
