package dao

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"shop/product/internal/domain/entity"
)

// CategoryBrandDAO handles data access operations for category-brand relations
type CategoryBrandDAO struct {
	db *gorm.DB
}

// NewCategoryBrandDAO creates a new category-brand data access object
func NewCategoryBrandDAO(db *gorm.DB) *CategoryBrandDAO {
	return &CategoryBrandDAO{db: db}
}

// Create adds a new category-brand relation
func (d *CategoryBrandDAO) Create(ctx context.Context, categoryBrand *entity.CategoryBrand) error {
	return d.db.WithContext(ctx).Create(categoryBrand).Error
}

// Update modifies an existing category-brand relation
func (d *CategoryBrandDAO) Update(ctx context.Context, categoryBrand *entity.CategoryBrand) error {
	return d.db.WithContext(ctx).Updates(categoryBrand).Error
}

// Delete soft-deletes a category-brand relation by ID
func (d *CategoryBrandDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Delete(&entity.CategoryBrand{}, id).Error
}

// GetByID retrieves a category-brand relation by ID
func (d *CategoryBrandDAO) GetByID(ctx context.Context, id int64) (*entity.CategoryBrand, error) {
	var categoryBrand entity.CategoryBrand
	result := d.db.WithContext(ctx).First(&categoryBrand, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category-brand relation not found: %v", result.Error)
		}
		return nil, result.Error
	}
	return &categoryBrand, nil
}

// List retrieves category-brand relations with pagination
func (d *CategoryBrandDAO) List(ctx context.Context, page, pageSize int) ([]*entity.CategoryBrand, int64, error) {
	var categoryBrands []*entity.CategoryBrand
	var total int64
	offset := (page - 1) * pageSize

	// Get total count
	if err := d.db.WithContext(ctx).Model(&entity.CategoryBrand{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := d.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&categoryBrands).Error; err != nil {
		return nil, 0, err
	}

	// Load relations for each result
	for _, cb := range categoryBrands {
		var brand entity.Brand
		if err := d.db.WithContext(ctx).First(&brand, cb.BrandID).Error; err == nil {
			cb.Brand = &brand
		}

		var category entity.Category
		if err := d.db.WithContext(ctx).First(&category, cb.CategoryID).Error; err == nil {
			cb.Category = &category
		}
	}

	return categoryBrands, total, nil
}

// GetBrandsByCategoryID retrieves all brands for a category
func (d *CategoryBrandDAO) GetBrandsByCategoryID(ctx context.Context, categoryID int64) ([]*entity.Brand, error) {
	var categoryBrands []*entity.CategoryBrand
	if err := d.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&categoryBrands).Error; err != nil {
		return nil, err
	}

	// Extract brand IDs
	brandIDs := make([]int64, len(categoryBrands))
	for i, cb := range categoryBrands {
		brandIDs[i] = cb.BrandID
	}

	// Get brands by IDs
	var brands []*entity.Brand
	if err := d.db.WithContext(ctx).Where("id IN ?", brandIDs).Find(&brands).Error; err != nil {
		return nil, err
	}

	return brands, nil
}

// GetCategoriesByBrandID retrieves all categories for a brand
func (d *CategoryBrandDAO) GetCategoriesByBrandID(ctx context.Context, brandID int64) ([]*entity.Category, error) {
	var categoryBrands []*entity.CategoryBrand
	if err := d.db.WithContext(ctx).Where("brands_id = ?", brandID).Find(&categoryBrands).Error; err != nil {
		return nil, err
	}

	// Extract category IDs
	categoryIDs := make([]int64, len(categoryBrands))
	for i, cb := range categoryBrands {
		categoryIDs[i] = cb.CategoryID
	}

	// Get categories by IDs
	var categories []*entity.Category
	if err := d.db.WithContext(ctx).Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
