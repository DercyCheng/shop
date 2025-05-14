package repository

import (
	"context"
	"errors"
	
	"shop/backend/product/internal/domain/entity"
	"shop/backend/product/internal/repository/cache"
	"shop/backend/product/internal/service"
	
	"gorm.io/gorm"
)

// CategoryRepositoryImpl 分类仓储实现
type CategoryRepositoryImpl struct {
	db    *gorm.DB
	cache cache.ProductCache
}

// NewCategoryRepository 创建分类仓储实例
func NewCategoryRepository(db *gorm.DB, cache cache.ProductCache) service.CategoryRepository {
	return &CategoryRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// GetCategoryByID 根据ID获取分类
func (r *CategoryRepositoryImpl) GetCategoryByID(ctx context.Context, id int64) (*entity.Category, error) {
	// 尝试从缓存获取
	category, err := r.cache.GetCategory(ctx, id)
	if err == nil && category != nil {
		return category, nil
	}
	
	// 缓存未命中，从数据库获取
	category = &entity.Category{}
	result := r.db.WithContext(ctx).First(category, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	// 将分类放入缓存
	if err := r.cache.SetCategory(ctx, category); err != nil {
		// 缓存失败只记录日志，不影响主流程
		// log.Printf("Cache category failed: %v", err)
	}
	
	return category, nil
}

// GetAllCategories 获取所有分类
func (r *CategoryRepositoryImpl) GetAllCategories(ctx context.Context) ([]*entity.Category, error) {
	// 尝试从缓存获取分类树
	categories, err := r.cache.GetCategoryTree(ctx)
	if err == nil && len(categories) > 0 {
		return categories, nil
	}
	
	// 从数据库获取
	var dbCategories []*entity.Category
	if err := r.db.WithContext(ctx).Order("level, parent_category_id, id").Find(&dbCategories).Error; err != nil {
		return nil, err
	}
	
	// 将分类树放入缓存
	if err := r.cache.SetCategoryTree(ctx, dbCategories); err != nil {
		// 缓存失败只记录日志，不影响主流程
		// log.Printf("Cache category tree failed: %v", err)
	}
	
	return dbCategories, nil
}

// GetSubCategories 获取子分类
func (r *CategoryRepositoryImpl) GetSubCategories(ctx context.Context, parentID int64) ([]*entity.Category, error) {
	var categories []*entity.Category
	if err := r.db.WithContext(ctx).Where("parent_category_id = ?", parentID).Find(&categories).Error; err != nil {
		return nil, err
	}
	
	return categories, nil
}

// GetCategoriesByLevel 获取指定级别的分类
func (r *CategoryRepositoryImpl) GetCategoriesByLevel(ctx context.Context, level int) ([]*entity.Category, error) {
	var categories []*entity.Category
	if err := r.db.WithContext(ctx).Where("level = ?", level).Find(&categories).Error; err != nil {
		return nil, err
	}
	
	return categories, nil
}

// CreateCategory 创建分类
func (r *CategoryRepositoryImpl) CreateCategory(ctx context.Context, category *entity.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return err
	}
	
	// 清除分类树缓存，强制重新加载
	if err := r.cache.DeleteCategoryTree(ctx); err != nil {
		// 缓存删除失败只记录日志，不影响主流程
		// log.Printf("Delete category tree cache failed: %v", err)
	}
	
	return nil
}

// UpdateCategory 更新分类
func (r *CategoryRepositoryImpl) UpdateCategory(ctx context.Context, category *entity.Category) error {
	// 更新数据库
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		return err
	}
	
	// 更新缓存
	if err := r.cache.SetCategory(ctx, category); err != nil {
		// 缓存失败只记录日志，不影响主流程
		// log.Printf("Update category cache failed: %v", err)
	}
	
	// 清除分类树缓存，强制重新加载
	if err := r.cache.DeleteCategoryTree(ctx); err != nil {
		// 缓存删除失败只记录日志，不影响主流程
		// log.Printf("Delete category tree cache failed: %v", err)
	}
	
	return nil
}

// DeleteCategory 删除分类
func (r *CategoryRepositoryImpl) DeleteCategory(ctx context.Context, id int64) error {
	// 删除分类
	if err := r.db.WithContext(ctx).Delete(&entity.Category{}, id).Error; err != nil {
		return err
	}
	
	// 删除缓存
	if err := r.cache.DeleteCategory(ctx, id); err != nil {
		// 缓存删除失败只记录日志，不影响主流程
		// log.Printf("Delete category cache failed: %v", err)
	}
	
	// 清除分类树缓存，强制重新加载
	if err := r.cache.DeleteCategoryTree(ctx); err != nil {
		// 缓存删除失败只记录日志，不影响主流程
		// log.Printf("Delete category tree cache failed: %v", err)
	}
	
	return nil
}
