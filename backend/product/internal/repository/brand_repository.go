package repository

import (
	"context"

	"shop/product/internal/domain/entity"
)

// BrandRepository defines the interface for brand data operations
type BrandRepository interface {
	// Create adds a new brand
	Create(ctx context.Context, brand *entity.Brand) error

	// Update modifies an existing brand
	Update(ctx context.Context, brand *entity.Brand) error

	// Delete removes a brand by ID
	Delete(ctx context.Context, id int64) error

	// GetByID retrieves a brand by ID
	GetByID(ctx context.Context, id int64) (*entity.Brand, error)

	// List retrieves brands with pagination
	List(ctx context.Context, page, pageSize int) ([]*entity.Brand, int64, error)

	// ListAll retrieves all brands
	ListAll(ctx context.Context) ([]*entity.Brand, error)

	// ListByIDs retrieves brands by a slice of IDs
	ListByIDs(ctx context.Context, ids []int64) ([]*entity.Brand, error)
}
