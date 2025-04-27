package repository

import (
	"context"

	"shop/product/internal/domain/entity"
)

// CategoryBrandRepository defines the interface for category-brand relation data operations
type CategoryBrandRepository interface {
	// Create adds a new category-brand relation
	Create(ctx context.Context, categoryBrand *entity.CategoryBrand) error

	// Update modifies an existing category-brand relation
	Update(ctx context.Context, categoryBrand *entity.CategoryBrand) error

	// Delete removes a category-brand relation by ID
	Delete(ctx context.Context, id int64) error

	// GetByID retrieves a category-brand relation by ID
	GetByID(ctx context.Context, id int64) (*entity.CategoryBrand, error)

	// List retrieves category-brand relations with pagination
	List(ctx context.Context, page, pageSize int) ([]*entity.CategoryBrand, int64, error)

	// GetBrandsByCategoryID retrieves all brands for a category
	GetBrandsByCategoryID(ctx context.Context, categoryID int64) ([]*entity.Brand, error)

	// GetCategoriesByBrandID retrieves all categories for a brand
	GetCategoriesByBrandID(ctx context.Context, brandID int64) ([]*entity.Category, error)
}
