package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shop/product/internal/domain/entity"
	"shop/product/internal/repository"
)

// CategoryBrandServiceImpl implements CategoryBrandService interface
type CategoryBrandServiceImpl struct {
	categoryBrandRepo repository.CategoryBrandRepository
	categoryRepo      repository.CategoryRepository
	brandRepo         repository.BrandRepository
}

// NewCategoryBrandService creates a new category-brand service
func NewCategoryBrandService(
	categoryBrandRepo repository.CategoryBrandRepository,
	categoryRepo repository.CategoryRepository,
	brandRepo repository.BrandRepository,
) CategoryBrandService {
	return &CategoryBrandServiceImpl{
		categoryBrandRepo: categoryBrandRepo,
		categoryRepo:      categoryRepo,
		brandRepo:         brandRepo,
	}
}

// GetCategoryBrandByID retrieves a category-brand relation by ID
func (s *CategoryBrandServiceImpl) GetCategoryBrandByID(ctx context.Context, id int64) (*entity.CategoryBrand, error) {
	return s.categoryBrandRepo.GetByID(ctx, id)
}

// ListCategoryBrands retrieves category-brand relations with pagination
func (s *CategoryBrandServiceImpl) ListCategoryBrands(ctx context.Context, page, pageSize int) ([]*entity.CategoryBrand, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.categoryBrandRepo.List(ctx, page, pageSize)
}

// GetBrandsByCategoryID retrieves all brands for a category
func (s *CategoryBrandServiceImpl) GetBrandsByCategoryID(ctx context.Context, categoryID int64) ([]*entity.Brand, error) {
	// Verify category exists
	_, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %v", err)
	}
	return s.categoryBrandRepo.GetBrandsByCategoryID(ctx, categoryID)
}

// GetCategoriesByBrandID retrieves all categories for a brand
func (s *CategoryBrandServiceImpl) GetCategoriesByBrandID(ctx context.Context, brandID int64) ([]*entity.Category, error) {
	// Verify brand exists
	_, err := s.brandRepo.GetByID(ctx, brandID)
	if err != nil {
		return nil, fmt.Errorf("brand not found: %v", err)
	}
	return s.categoryBrandRepo.GetCategoriesByBrandID(ctx, brandID)
}

// CreateCategoryBrand adds a new category-brand relation
func (s *CategoryBrandServiceImpl) CreateCategoryBrand(ctx context.Context, categoryBrand *entity.CategoryBrand) (*entity.CategoryBrand, error) {
	// Validate required fields
	if categoryBrand.CategoryID <= 0 {
		return nil, errors.New("category ID is required")
	}
	if categoryBrand.BrandID <= 0 {
		return nil, errors.New("brand ID is required")
	}

	// Verify category exists
	_, err := s.categoryRepo.GetByID(ctx, categoryBrand.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %v", err)
	}

	// Verify brand exists
	_, err = s.brandRepo.GetByID(ctx, categoryBrand.BrandID)
	if err != nil {
		return nil, fmt.Errorf("brand not found: %v", err)
	}

	// Set timestamps
	now := time.Now()
	categoryBrand.CreatedAt = now
	categoryBrand.UpdatedAt = now

	// Create the relationship
	if err := s.categoryBrandRepo.Create(ctx, categoryBrand); err != nil {
		return nil, err
	}

	// Get full details with relationships
	return s.categoryBrandRepo.GetByID(ctx, categoryBrand.ID)
}

// UpdateCategoryBrand modifies an existing category-brand relation
func (s *CategoryBrandServiceImpl) UpdateCategoryBrand(ctx context.Context, categoryBrand *entity.CategoryBrand) error {
	// Verify relationship exists
	existing, err := s.categoryBrandRepo.GetByID(ctx, categoryBrand.ID)
	if err != nil {
		return fmt.Errorf("category-brand relation not found: %v", err)
	}

	// Validate and check category if changed
	if categoryBrand.CategoryID > 0 && categoryBrand.CategoryID != existing.CategoryID {
		_, err := s.categoryRepo.GetByID(ctx, categoryBrand.CategoryID)
		if err != nil {
			return fmt.Errorf("category not found: %v", err)
		}
	} else {
		categoryBrand.CategoryID = existing.CategoryID
	}

	// Validate and check brand if changed
	if categoryBrand.BrandID > 0 && categoryBrand.BrandID != existing.BrandID {
		_, err := s.brandRepo.GetByID(ctx, categoryBrand.BrandID)
		if err != nil {
			return fmt.Errorf("brand not found: %v", err)
		}
	} else {
		categoryBrand.BrandID = existing.BrandID
	}

	// Update timestamp
	categoryBrand.UpdatedAt = time.Now()

	return s.categoryBrandRepo.Update(ctx, categoryBrand)
}

// DeleteCategoryBrand removes a category-brand relation by ID
func (s *CategoryBrandServiceImpl) DeleteCategoryBrand(ctx context.Context, id int64) error {
	return s.categoryBrandRepo.Delete(ctx, id)
}
