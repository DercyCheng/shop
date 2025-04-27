package service

import (
	"context"

	"shop/product/internal/domain/entity"
)

// CategoryBrandService defines the interface for category-brand relation business logic
type CategoryBrandService interface {
	// GetCategoryBrandByID retrieves a category-brand relation by ID
	GetCategoryBrandByID(ctx context.Context, id int64) (*entity.CategoryBrand, error)

	// ListCategoryBrands retrieves category-brand relations with pagination
	ListCategoryBrands(ctx context.Context, page, pageSize int) ([]*entity.CategoryBrand, int64, error)

	// GetBrandsByCategoryID retrieves all brands for a category
	GetBrandsByCategoryID(ctx context.Context, categoryID int64) ([]*entity.Brand, error)

	// GetCategoriesByBrandID retrieves all categories for a brand
	GetCategoriesByBrandID(ctx context.Context, brandID int64) ([]*entity.Category, error)

	// CreateCategoryBrand adds a new category-brand relation
	CreateCategoryBrand(ctx context.Context, categoryBrand *entity.CategoryBrand) (*entity.CategoryBrand, error)

	// UpdateCategoryBrand modifies an existing category-brand relation
	UpdateCategoryBrand(ctx context.Context, categoryBrand *entity.CategoryBrand) error

	// DeleteCategoryBrand removes a category-brand relation by ID
	DeleteCategoryBrand(ctx context.Context, id int64) error
}
