package service

import (
	"context"

	"shop/product/internal/domain/entity"
)

// CategoryService defines the interface for category business logic
type CategoryService interface {
	// GetCategoryByID retrieves a category by ID
	GetCategoryByID(ctx context.Context, id int64) (*entity.Category, error)

	// GetAllCategories retrieves all categories
	GetAllCategories(ctx context.Context) ([]*entity.Category, error)

	// GetCategoriesByParentID retrieves categories by parent ID
	GetCategoriesByParentID(ctx context.Context, parentID int64) ([]*entity.Category, error)

	// GetCategoriesByLevel retrieves categories by level
	GetCategoriesByLevel(ctx context.Context, level int) ([]*entity.Category, error)

	// GetCategoryTree retrieves a category with all its subcategories
	GetCategoryTree(ctx context.Context, parentID int64) (*entity.Category, error)

	// CreateCategory adds a new category
	CreateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error)

	// UpdateCategory modifies an existing category
	UpdateCategory(ctx context.Context, category *entity.Category) error

	// DeleteCategory removes a category by ID
	DeleteCategory(ctx context.Context, id int64) error
}
