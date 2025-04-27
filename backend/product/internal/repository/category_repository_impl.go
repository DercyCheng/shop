package repository

import (
	"context"

	"shop/product/internal/domain/entity"
	"shop/product/internal/repository/dao"
)

// CategoryRepositoryImpl implements CategoryRepository interface
type CategoryRepositoryImpl struct {
	categoryDAO *dao.CategoryDAO
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(categoryDAO *dao.CategoryDAO) CategoryRepository {
	return &CategoryRepositoryImpl{
		categoryDAO: categoryDAO,
	}
}

// Create adds a new category
func (r *CategoryRepositoryImpl) Create(ctx context.Context, category *entity.Category) error {
	return r.categoryDAO.Create(ctx, category)
}

// Update modifies an existing category
func (r *CategoryRepositoryImpl) Update(ctx context.Context, category *entity.Category) error {
	return r.categoryDAO.Update(ctx, category)
}

// Delete removes a category by ID
func (r *CategoryRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.categoryDAO.Delete(ctx, id)
}

// GetByID retrieves a category by ID
func (r *CategoryRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	return r.categoryDAO.GetByID(ctx, id)
}

// ListAll retrieves all categories
func (r *CategoryRepositoryImpl) ListAll(ctx context.Context) ([]*entity.Category, error) {
	return r.categoryDAO.ListAll(ctx)
}

// ListByParentID retrieves categories by parent ID
func (r *CategoryRepositoryImpl) ListByParentID(ctx context.Context, parentID int64) ([]*entity.Category, error) {
	return r.categoryDAO.ListByParentID(ctx, parentID)
}

// ListByLevel retrieves categories by level
func (r *CategoryRepositoryImpl) ListByLevel(ctx context.Context, level int) ([]*entity.Category, error) {
	return r.categoryDAO.ListByLevel(ctx, level)
}

// ListSubCategories retrieves a category with all its subcategories
func (r *CategoryRepositoryImpl) ListSubCategories(ctx context.Context, parentID int64) (*entity.Category, error) {
	return r.categoryDAO.ListSubCategories(ctx, parentID)
}
