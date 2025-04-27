package dao

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"shop/product/internal/domain/entity"
)

// CategoryDAO handles data access operations for categories
type CategoryDAO struct {
	db *gorm.DB
}

// NewCategoryDAO creates a new category data access object
func NewCategoryDAO(db *gorm.DB) *CategoryDAO {
	return &CategoryDAO{db: db}
}

// Create adds a new category to the database
func (d *CategoryDAO) Create(ctx context.Context, category *entity.Category) error {
	return d.db.WithContext(ctx).Create(category).Error
}

// Update modifies an existing category
func (d *CategoryDAO) Update(ctx context.Context, category *entity.Category) error {
	return d.db.WithContext(ctx).Updates(category).Error
}

// Delete soft-deletes a category by ID
func (d *CategoryDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Delete(&entity.Category{}, id).Error
}

// GetByID retrieves a category by ID
func (d *CategoryDAO) GetByID(ctx context.Context, id int64) (*entity.Category, error) {
	var category entity.Category
	result := d.db.WithContext(ctx).First(&category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category not found: %v", result.Error)
		}
		return nil, result.Error
	}
	return &category, nil
}

// ListAll retrieves all categories
func (d *CategoryDAO) ListAll(ctx context.Context) ([]*entity.Category, error) {
	var categories []*entity.Category
	if err := d.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// ListByParentID retrieves categories by parent ID
func (d *CategoryDAO) ListByParentID(ctx context.Context, parentID int64) ([]*entity.Category, error) {
	var categories []*entity.Category
	if err := d.db.WithContext(ctx).Where("parent_category_id = ?", parentID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// ListByLevel retrieves categories by level
func (d *CategoryDAO) ListByLevel(ctx context.Context, level int) ([]*entity.Category, error) {
	var categories []*entity.Category
	if err := d.db.WithContext(ctx).Where("level = ?", level).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// ListSubCategories recursively retrieves a category with all its subcategories
func (d *CategoryDAO) ListSubCategories(ctx context.Context, parentID int64) (*entity.Category, error) {
	// Get parent category
	var parent entity.Category
	if parentID == 0 {
		// If parent ID is 0, return all top-level categories
		var topLevelCategories []*entity.Category
		if err := d.db.WithContext(ctx).Where("parent_category_id = 0").Find(&topLevelCategories).Error; err != nil {
			return nil, err
		}

		// Create a virtual root category
		parent = entity.Category{
			ID:               0,
			Name:             "Root",
			ParentCategoryID: 0,
			Level:            0,
			SubCategories:    topLevelCategories,
		}
	} else {
		// Get the actual parent category
		if err := d.db.WithContext(ctx).First(&parent, parentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("parent category not found: %v", err)
			}
			return nil, err
		}

		// Get immediate children
		var subCategories []*entity.Category
		if err := d.db.WithContext(ctx).Where("parent_category_id = ?", parentID).Find(&subCategories).Error; err != nil {
			return nil, err
		}
		parent.SubCategories = subCategories
	}

	// Recursively get sub-categories for each child
	for _, subCategory := range parent.SubCategories {
		subWithChildren, err := d.ListSubCategories(ctx, subCategory.ID)
		if err != nil {
			return nil, err
		}
		subCategory.SubCategories = subWithChildren.SubCategories
	}

	return &parent, nil
}
