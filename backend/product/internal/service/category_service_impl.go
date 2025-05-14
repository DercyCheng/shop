package service

import (
	"context"
	"errors"
	"time"
	
	"shop/backend/product/internal/domain/entity"
)

var (
	ErrCategoryNotFound        = errors.New("category not found")
	ErrInvalidCategory         = errors.New("invalid category data")
	ErrCategoryHasSubcategories = errors.New("category has subcategories")
	ErrCategoryHasProducts     = errors.New("category has associated products")
	ErrCircularReference       = errors.New("circular reference detected")
)

// CategoryServiceImpl 分类服务实现
type CategoryServiceImpl struct {
	categoryRepo CategoryRepository
	productRepo  ProductRepository
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(
	categoryRepo CategoryRepository,
	productRepo ProductRepository,
) CategoryService {
	return &CategoryServiceImpl{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryServiceImpl) GetCategoryByID(ctx context.Context, id int64) (*entity.Category, error) {
	category, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	
	return category, nil
}

// GetAllCategories 获取所有分类
func (s *CategoryServiceImpl) GetAllCategories(ctx context.Context) ([]*entity.Category, error) {
	return s.categoryRepo.ListAllCategories(ctx)
}

// GetCategoryTree 获取分类树
func (s *CategoryServiceImpl) GetCategoryTree(ctx context.Context) ([]*entity.Category, error) {
	// 获取所有分类
	allCategories, err := s.categoryRepo.ListAllCategories(ctx)
	if err != nil {
		return nil, err
	}
	
	// 按父级分类ID组织分类树
	categoryMap := make(map[int64][]*entity.Category)
	for _, category := range allCategories {
		parentID := category.ParentCategoryID
		categoryMap[parentID] = append(categoryMap[parentID], category)
	}
	
	// 获取根分类（ParentCategoryID为0的分类）
	rootCategories := categoryMap[0]
	
	// 递归构建分类树
	buildCategoryTree(rootCategories, categoryMap)
	
	return rootCategories, nil
}

// buildCategoryTree 构建分类树的递归辅助函数
func buildCategoryTree(categories []*entity.Category, categoryMap map[int64][]*entity.Category) {
	for _, category := range categories {
		// 查找子分类
		if subCategories, exists := categoryMap[category.ID]; exists {
			category.SubCategories = subCategories
			// 递归处理子分类
			buildCategoryTree(subCategories, categoryMap)
		}
	}
}

// GetSubCategories 获取子分类
func (s *CategoryServiceImpl) GetSubCategories(ctx context.Context, parentID int64) ([]*entity.Category, error) {
	return s.categoryRepo.ListCategoriesByParentID(ctx, parentID)
}

// CreateCategory 创建分类
func (s *CategoryServiceImpl) CreateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	// 基本参数验证
	if category.Name == "" {
		return nil, ErrInvalidCategory
	}
	
	// 如果设置了父分类，检查父分类是否存在
	if category.ParentCategoryID > 0 {
		parent, err := s.categoryRepo.GetCategoryByID(ctx, category.ParentCategoryID)
		if err != nil {
			return nil, err
		}
		
		if parent == nil {
			return nil, errors.New("parent category not found")
		}
		
		// 设置分类级别
		category.Level = parent.Level + 1
	} else {
		// 根分类
		category.Level = 1
	}
	
	// 设置初始值
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now
	
	// 保存分类
	if err := s.categoryRepo.CreateCategory(ctx, category); err != nil {
		return nil, err
	}
	
	return category, nil
}

// UpdateCategory 更新分类
func (s *CategoryServiceImpl) UpdateCategory(ctx context.Context, category *entity.Category) error {
	// 检查分类是否存在
	existingCategory, err := s.categoryRepo.GetCategoryByID(ctx, category.ID)
	if err != nil {
		return err
	}
	
	if existingCategory == nil {
		return ErrCategoryNotFound
	}
	
	// 如果更改了父分类，需要检查
	if category.ParentCategoryID != existingCategory.ParentCategoryID {
		// 检查循环引用：一个分类不能将自己或其子分类设为父分类
		if category.ParentCategoryID > 0 && category.ParentCategoryID == category.ID {
			return ErrCircularReference
		}
		
		// 检查新的父分类是否存在
		if category.ParentCategoryID > 0 {
			parent, err := s.categoryRepo.GetCategoryByID(ctx, category.ParentCategoryID)
			if err != nil {
				return err
			}
			
			if parent == nil {
				return errors.New("parent category not found")
			}
			
			// 检查新父分类是否为此分类的子分类（避免循环引用）
			isCircular, err := s.isChildCategory(ctx, category.ID, category.ParentCategoryID)
			if err != nil {
				return err
			}
			
			if isCircular {
				return ErrCircularReference
			}
			
			// 更新分类级别
			category.Level = parent.Level + 1
		} else {
			// 变更为根分类
			category.Level = 1
		}
	} else {
		// 保持原有级别
		category.Level = existingCategory.Level
	}
	
	// 更新时间
	category.UpdatedAt = time.Now()
	category.CreatedAt = existingCategory.CreatedAt
	
	// 保存分类
	if err := s.categoryRepo.UpdateCategory(ctx, category); err != nil {
		return err
	}
	
	return nil
}

// DeleteCategory 删除分类
func (s *CategoryServiceImpl) DeleteCategory(ctx context.Context, id int64) error {
	// 检查分类是否存在
	category, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return err
	}
	
	if category == nil {
		return ErrCategoryNotFound
	}
	
	// 检查是否有子分类
	subcategories, err := s.categoryRepo.ListCategoriesByParentID(ctx, id)
	if err != nil {
		return err
	}
	
	if len(subcategories) > 0 {
		return ErrCategoryHasSubcategories
	}
	
	// 检查分类下是否有商品
	// 此处需要ProductRepository的新方法，暂时不实现这个检查
	// 或者可以通过计数查询实现
	
	// 删除分类
	if err := s.categoryRepo.DeleteCategory(ctx, id); err != nil {
		return err
	}
	
	return nil
}

// isChildCategory 检查potentialChild是否是parentID的子分类（任意层级）
func (s *CategoryServiceImpl) isChildCategory(ctx context.Context, parentID, potentialChildID int64) (bool, error) {
	children, err := s.categoryRepo.ListCategoriesByParentID(ctx, parentID)
	if err != nil {
		return false, err
	}
	
	for _, child := range children {
		if child.ID == potentialChildID {
			return true, nil
		}
		
		// 递归检查
		isChild, err := s.isChildCategory(ctx, child.ID, potentialChildID)
		if err != nil {
			return false, err
		}
		
		if isChild {
			return true, nil
		}
	}
	
	return false, nil
}
