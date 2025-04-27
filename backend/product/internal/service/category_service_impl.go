package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shop/product/internal/domain/entity"
	"shop/product/internal/repository"
)

// CategoryServiceImpl implements CategoryService interface
type CategoryServiceImpl struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{
		categoryRepo: categoryRepo,
	}
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryServiceImpl) GetCategoryByID(ctx context.Context, id int64) (*entity.Category, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

// GetAllCategories retrieves all categories
func (s *CategoryServiceImpl) GetAllCategories(ctx context.Context) ([]*entity.Category, error) {
	return s.categoryRepo.ListAll(ctx)
}

// GetCategoriesByParentID retrieves categories by parent ID
func (s *CategoryServiceImpl) GetCategoriesByParentID(ctx context.Context, parentID int64) ([]*entity.Category, error) {
	return s.categoryRepo.ListByParentID(ctx, parentID)
}

// GetCategoriesByLevel retrieves categories by level
func (s *CategoryServiceImpl) GetCategoriesByLevel(ctx context.Context, level int) ([]*entity.Category, error) {
	return s.categoryRepo.ListByLevel(ctx, level)
}

// GetCategoryTree retrieves a category with all its subcategories
func (s *CategoryServiceImpl) GetCategoryTree(ctx context.Context, parentID int64) (*entity.Category, error) {
	return s.categoryRepo.ListSubCategories(ctx, parentID)
}

// CreateCategory adds a new category
func (s *CategoryServiceImpl) CreateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	// Validate required fields
	if category.Name == "" {
		return nil, errors.New("category name is required")
	}

	// If parent category is specified, verify it exists
	if category.ParentCategoryID > 0 {
		parent, err := s.categoryRepo.GetByID(ctx, category.ParentCategoryID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %v", err)
		}

		// Set the level based on parent's level
		category.Level = parent.Level + 1
	} else {
		// Top-level category
		category.ParentCategoryID = 0
		category.Level = 1
	}

	// Set timestamps
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	// Create category
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory modifies an existing category
func (s *CategoryServiceImpl) UpdateCategory(ctx context.Context, category *entity.Category) error {
	// Verify category exists
	existingCategory, err := s.categoryRepo.GetByID(ctx, category.ID)
	if err != nil {
		return fmt.Errorf("category not found: %v", err)
	}

	// Validate name if provided
	if category.Name == "" {
		category.Name = existingCategory.Name
	}

	// Check if parent category changed
	if category.ParentCategoryID > 0 && category.ParentCategoryID != existingCategory.ParentCategoryID {
		// Verify new parent exists
		parent, err := s.categoryRepo.GetByID(ctx, category.ParentCategoryID)
		if err != nil {
			return fmt.Errorf("parent category not found: %v", err)
		}

		// Prevent circular reference
		if category.ID == category.ParentCategoryID {
			return errors.New("category cannot be its own parent")
		}

		// Set level based on new parent
		category.Level = parent.Level + 1
	} else {
		// Keep existing parent and level if not changing
		if category.ParentCategoryID <= 0 {
			category.ParentCategoryID = existingCategory.ParentCategoryID
		}
		if category.Level <= 0 {
			category.Level = existingCategory.Level
		}
	}

	// Update timestamp
	category.UpdatedAt = time.Now()

	return s.categoryRepo.Update(ctx, category)
}

// DeleteCategory removes a category by ID
func (s *CategoryServiceImpl) DeleteCategory(ctx context.Context, id int64) error {
	// Check if there are subcategories
	subcategories, err := s.categoryRepo.ListByParentID(ctx, id)
	if err != nil {
		return err
	}

	// Prevent deleting categories with subcategories
	if len(subcategories) > 0 {
		return errors.New("cannot delete category with subcategories")
	}

	// Delete the category
	return s.categoryRepo.Delete(ctx, id)
}
