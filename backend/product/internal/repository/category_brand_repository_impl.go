package repository

import (
	"context"

	"shop/product/internal/domain/entity"
	"shop/product/internal/repository/dao"
)

// CategoryBrandRepositoryImpl implements CategoryBrandRepository interface
type CategoryBrandRepositoryImpl struct {
	categoryBrandDAO *dao.CategoryBrandDAO
}

// NewCategoryBrandRepository creates a new category-brand repository
func NewCategoryBrandRepository(categoryBrandDAO *dao.CategoryBrandDAO) CategoryBrandRepository {
	return &CategoryBrandRepositoryImpl{
		categoryBrandDAO: categoryBrandDAO,
	}
}

// Create adds a new category-brand relation
func (r *CategoryBrandRepositoryImpl) Create(ctx context.Context, categoryBrand *entity.CategoryBrand) error {
	return r.categoryBrandDAO.Create(ctx, categoryBrand)
}

// Update modifies an existing category-brand relation
func (r *CategoryBrandRepositoryImpl) Update(ctx context.Context, categoryBrand *entity.CategoryBrand) error {
	return r.categoryBrandDAO.Update(ctx, categoryBrand)
}

// Delete removes a category-brand relation by ID
func (r *CategoryBrandRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.categoryBrandDAO.Delete(ctx, id)
}

// GetByID retrieves a category-brand relation by ID
func (r *CategoryBrandRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.CategoryBrand, error) {
	return r.categoryBrandDAO.GetByID(ctx, id)
}

// List retrieves category-brand relations with pagination
func (r *CategoryBrandRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.CategoryBrand, int64, error) {
	return r.categoryBrandDAO.List(ctx, page, pageSize)
}

// GetBrandsByCategoryID retrieves all brands for a category
func (r *CategoryBrandRepositoryImpl) GetBrandsByCategoryID(ctx context.Context, categoryID int64) ([]*entity.Brand, error) {
	return r.categoryBrandDAO.GetBrandsByCategoryID(ctx, categoryID)
}

// GetCategoriesByBrandID retrieves all categories for a brand
func (r *CategoryBrandRepositoryImpl) GetCategoriesByBrandID(ctx context.Context, brandID int64) ([]*entity.Category, error) {
	return r.categoryBrandDAO.GetCategoriesByBrandID(ctx, brandID)
}
