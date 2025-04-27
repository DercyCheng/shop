package repository

import (
	"context"

	"shop/product/internal/domain/entity"
	"shop/product/internal/repository/dao"
)

// BrandRepositoryImpl implements BrandRepository interface
type BrandRepositoryImpl struct {
	brandDAO *dao.BrandDAO
}

// NewBrandRepository creates a new brand repository
func NewBrandRepository(brandDAO *dao.BrandDAO) BrandRepository {
	return &BrandRepositoryImpl{
		brandDAO: brandDAO,
	}
}

// Create adds a new brand
func (r *BrandRepositoryImpl) Create(ctx context.Context, brand *entity.Brand) error {
	return r.brandDAO.Create(ctx, brand)
}

// Update modifies an existing brand
func (r *BrandRepositoryImpl) Update(ctx context.Context, brand *entity.Brand) error {
	return r.brandDAO.Update(ctx, brand)
}

// Delete removes a brand by ID
func (r *BrandRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.brandDAO.Delete(ctx, id)
}

// GetByID retrieves a brand by ID
func (r *BrandRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Brand, error) {
	return r.brandDAO.GetByID(ctx, id)
}

// List retrieves brands with pagination
func (r *BrandRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*entity.Brand, int64, error) {
	return r.brandDAO.List(ctx, page, pageSize)
}

// ListAll retrieves all brands
func (r *BrandRepositoryImpl) ListAll(ctx context.Context) ([]*entity.Brand, error) {
	return r.brandDAO.ListAll(ctx)
}

// ListByIDs retrieves brands by a slice of IDs
func (r *BrandRepositoryImpl) ListByIDs(ctx context.Context, ids []int64) ([]*entity.Brand, error) {
	return r.brandDAO.ListByIDs(ctx, ids)
}
