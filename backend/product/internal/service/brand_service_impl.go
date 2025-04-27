package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shop/product/internal/domain/entity"
	"shop/product/internal/repository"
)

// BrandServiceImpl implements BrandService interface
type BrandServiceImpl struct {
	brandRepo repository.BrandRepository
}

// NewBrandService creates a new brand service
func NewBrandService(brandRepo repository.BrandRepository) BrandService {
	return &BrandServiceImpl{
		brandRepo: brandRepo,
	}
}

// GetBrandByID retrieves a brand by ID
func (s *BrandServiceImpl) GetBrandByID(ctx context.Context, id int64) (*entity.Brand, error) {
	return s.brandRepo.GetByID(ctx, id)
}

// ListBrands retrieves brands with pagination
func (s *BrandServiceImpl) ListBrands(ctx context.Context, page, pageSize int) ([]*entity.Brand, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.brandRepo.List(ctx, page, pageSize)
}

// GetAllBrands retrieves all brands
func (s *BrandServiceImpl) GetAllBrands(ctx context.Context) ([]*entity.Brand, error) {
	return s.brandRepo.ListAll(ctx)
}

// CreateBrand adds a new brand
func (s *BrandServiceImpl) CreateBrand(ctx context.Context, brand *entity.Brand) (*entity.Brand, error) {
	// Validate required fields
	if brand.Name == "" {
		return nil, errors.New("brand name is required")
	}

	// Set timestamps
	now := time.Now()
	brand.CreatedAt = now
	brand.UpdatedAt = now

	// Create brand
	if err := s.brandRepo.Create(ctx, brand); err != nil {
		return nil, err
	}

	return brand, nil
}

// UpdateBrand modifies an existing brand
func (s *BrandServiceImpl) UpdateBrand(ctx context.Context, brand *entity.Brand) error {
	// Verify brand exists
	existingBrand, err := s.brandRepo.GetByID(ctx, brand.ID)
	if err != nil {
		return fmt.Errorf("brand not found: %v", err)
	}

	// Validate fields
	if brand.Name == "" {
		brand.Name = existingBrand.Name
	}

	// If logo not provided, keep the existing one
	if brand.Logo == "" {
		brand.Logo = existingBrand.Logo
	}

	// Update timestamp
	brand.UpdatedAt = time.Now()

	return s.brandRepo.Update(ctx, brand)
}

// DeleteBrand removes a brand by ID
func (s *BrandServiceImpl) DeleteBrand(ctx context.Context, id int64) error {
	return s.brandRepo.Delete(ctx, id)
}
