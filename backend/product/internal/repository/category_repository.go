package repository

import (
	"context"

	"shop/product/internal/domain/entity"
)

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	// Create adds a new category
	Create(ctx context.Context, category *entity.Category) error

	// Update modifies an existing category
	Update(ctx context.Context, category *entity.Category) error

	// Delete removes a category by ID
	Delete(ctx context.Context, id int64) error

	// GetByID retrieves a category by ID
	GetByID(ctx context.Context, id int64) (*entity.Category, error)

	// ListAll retrieves all categories
	ListAll(ctx context.Context) ([]*entity.Category, error)

	// ListByParentID retrieves categories by parent ID
	ListByParentID(ctx context.Context, parentID int64) ([]*entity.Category, error)

	// ListByLevel retrieves categories by level
	ListByLevel(ctx context.Context, level int) ([]*entity.Category, error)

	// ListSubCategories retrieves a category with all its subcategories
	ListSubCategories(ctx context.Context, parentID int64) (*entity.Category, error)
}
