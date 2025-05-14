package repository

import (
	"context"
	"errors"
	
	"shop/backend/product/internal/domain/entity"
	"shop/backend/product/internal/repository/cache"
	"shop/backend/product/internal/service"
	
	"gorm.io/gorm"
)

// BrandRepositoryImpl 品牌仓储实现
type BrandRepositoryImpl struct {
	db    *gorm.DB
	cache cache.ProductCache
}

// NewBrandRepository 创建品牌仓储实例
func NewBrandRepository(db *gorm.DB, cache cache.ProductCache) service.BrandRepository {
	return &BrandRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// GetBrandByID 根据ID获取品牌
func (r *BrandRepositoryImpl) GetBrandByID(ctx context.Context, id int64) (*entity.Brand, error) {
	// 尝试从缓存获取
	brand, err := r.cache.GetBrand(ctx, id)
	if err == nil && brand != nil {
		return brand, nil
	}
	
	// 缓存未命中，从数据库获取
	brand = &entity.Brand{}
	result := r.db.WithContext(ctx).First(brand, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	// 将品牌放入缓存
	if err := r.cache.SetBrand(ctx, brand); err != nil {
		// 缓存失败只记录日志，不影响主流程
		// log.Printf("Cache brand failed: %v", err)
	}
	
	return brand, nil
}

// ListBrands 获取品牌列表
func (r *BrandRepositoryImpl) ListBrands(ctx context.Context, filter service.BrandFilter) ([]*entity.Brand, int64, error) {
	var brands []*entity.Brand
	var total int64
	
	// 构建查询
	query := r.db.WithContext(ctx).Model(&entity.Brand{})
	
	// 应用过滤条件
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页
	query = query.Offset((filter.Page - 1) * filter.PageSize).Limit(filter.PageSize)
	
	// 排序
	query = query.Order("id")
	
	// 执行查询
	if err := query.Find(&brands).Error; err != nil {
		return nil, 0, err
	}
	
	return brands, total, nil
}

// GetBrandsByCategory 获取分类下的品牌列表
func (r *BrandRepositoryImpl) GetBrandsByCategory(ctx context.Context, categoryID int64) ([]*entity.Brand, error) {
	var brands []*entity.Brand
	
	// 通过关联表查询该分类下的品牌
	err := r.db.WithContext(ctx).
		Table("brands").
		Joins("JOIN goods_category_brand ON brands.id = goods_category_brand.brands_id").
		Where("goods_category_brand.category_id = ?", categoryID).
		Find(&brands).Error
	
	if err != nil {
		return nil, err
	}
	
	return brands, nil
}

// CreateBrand 创建品牌
func (r *BrandRepositoryImpl) CreateBrand(ctx context.Context, brand *entity.Brand) error {
	return r.db.WithContext(ctx).Create(brand).Error
}

// UpdateBrand 更新品牌
func (r *BrandRepositoryImpl) UpdateBrand(ctx context.Context, brand *entity.Brand) error {
	// 更新数据库
	if err := r.db.WithContext(ctx).Save(brand).Error; err != nil {
		return err
	}
	
	// 更新缓存
	if err := r.cache.SetBrand(ctx, brand); err != nil {
		// 缓存失败只记录日志，不影响主流程
		// log.Printf("Update brand cache failed: %v", err)
	}
	
	return nil
}

// DeleteBrand 删除品牌
func (r *BrandRepositoryImpl) DeleteBrand(ctx context.Context, id int64) error {
	// 删除品牌
	if err := r.db.WithContext(ctx).Delete(&entity.Brand{}, id).Error; err != nil {
		return err
	}
	
	// 删除缓存
	if err := r.cache.DeleteBrand(ctx, id); err != nil {
		// 缓存删除失败只记录日志，不影响主流程
		// log.Printf("Delete brand cache failed: %v", err)
	}
	
	return nil
}

// AddCategoryBrand 添加品牌分类关系
func (r *BrandRepositoryImpl) AddCategoryBrand(ctx context.Context, categoryID, brandID int64) error {
	// 检查关系是否已存在
	var count int64
	err := r.db.WithContext(ctx).
		Table("goods_category_brand").
		Where("category_id = ? AND brands_id = ?", categoryID, brandID).
		Count(&count).Error
	
	if err != nil {
		return err
	}
	
	// 关系已存在
	if count > 0 {
		return nil
	}
	
	// 添加关系
	return r.db.WithContext(ctx).
		Table("goods_category_brand").
		Create(map[string]interface{}{
			"category_id": categoryID,
			"brands_id":   brandID,
		}).Error
}

// RemoveCategoryBrand 删除品牌分类关系
func (r *BrandRepositoryImpl) RemoveCategoryBrand(ctx context.Context, categoryID, brandID int64) error {
	return r.db.WithContext(ctx).
		Table("goods_category_brand").
		Where("category_id = ? AND brands_id = ?", categoryID, brandID).
		Delete(nil).Error
}
