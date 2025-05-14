package repository

import (
	"context"
	"errors"
	
	"shop/backend/product/internal/domain/entity"
	"shop/backend/product/internal/repository/cache"
	"shop/backend/product/internal/service"
	
	"gorm.io/gorm"
)

// ProductRepositoryImpl 商品仓储实现
type ProductRepositoryImpl struct {
	db    *gorm.DB
	cache cache.ProductCache
}

// NewProductRepository 创建商品仓储实例
func NewProductRepository(db *gorm.DB, cache cache.ProductCache) service.ProductRepository {
	return &ProductRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// GetProductByID 根据ID获取商品
func (r *ProductRepositoryImpl) GetProductByID(ctx context.Context, id int64) (*entity.Product, error) {
	// 尝试从缓存获取
	product, err := r.cache.GetProduct(ctx, id)
	if err == nil && product != nil {
		return product, nil
	}
	
	// 缓存未命中，从数据库获取
	product = &entity.Product{}
	result := r.db.WithContext(ctx).First(product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	// 将商品放入缓存
	if err := r.cache.SetProduct(ctx, product); err != nil {
		// 缓存失败只记录日志，不影响主流程
		// log.Printf("Cache product failed: %v", err)
	}
	
	return product, nil
}

// GetProductBySN 根据SN获取商品
func (r *ProductRepositoryImpl) GetProductBySN(ctx context.Context, goodsSN string) (*entity.Product, error) {
	var product entity.Product
	result := r.db.WithContext(ctx).Where("goods_sn = ?", goodsSN).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return &product, nil
}

// ListProducts 获取商品列表
func (r *ProductRepositoryImpl) ListProducts(ctx context.Context, filter service.ProductFilter) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var total int64
	
	// 构建查询条件
	query := r.db.WithContext(ctx).Model(&entity.Product{})
	
	// 应用过滤条件
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	
	if filter.CategoryID > 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	
	if filter.BrandID > 0 {
		query = query.Where("brands_id = ?", filter.BrandID)
	}
	
	if filter.OnSale != nil {
		query = query.Where("on_sale = ?", *filter.OnSale)
	}
	
	if filter.IsNew != nil {
		query = query.Where("is_new = ?", *filter.IsNew)
	}
	
	if filter.IsHot != nil {
		query = query.Where("is_hot = ?", *filter.IsHot)
	}
	
	if filter.ShipFree != nil {
		query = query.Where("ship_free = ?", *filter.ShipFree)
	}
	
	if filter.PriceMin != nil {
		query = query.Where("shop_price >= ?", *filter.PriceMin)
	}
	
	if filter.PriceMax != nil {
		query = query.Where("shop_price <= ?", *filter.PriceMax)
	}
	
	// 非删除状态
	query = query.Where("is_deleted = ?", false)
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 排序
	if filter.OrderBy != "" {
		query = query.Order(filter.OrderBy)
	} else {
		query = query.Order("updated_at DESC")
	}
	
	// 分页
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)
	
	// 执行查询
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}
	
	return products, total, nil
}

// BatchGetProducts 批量获取商品
func (r *ProductRepositoryImpl) BatchGetProducts(ctx context.Context, ids []int64) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&products).Error
	return products, err
}

// CreateProduct 创建商品
func (r *ProductRepositoryImpl) CreateProduct(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// UpdateProduct 更新商品
func (r *ProductRepositoryImpl) UpdateProduct(ctx context.Context, product *entity.Product) error {
	// 更新数据库
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return err
	}
	
	// 更新缓存
	if err := r.cache.SetProduct(ctx, product); err != nil {
		// 缓存失败只记录日志，不影响主流程
		// log.Printf("Update product cache failed: %v", err)
	}
	
	return nil
}

// DeleteProduct 删除商品
func (r *ProductRepositoryImpl) DeleteProduct(ctx context.Context, id int64) error {
	// 软删除
	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": gorm.Expr("NOW()"),
		}).Error; err != nil {
		return err
	}
	
	// 删除缓存
	if err := r.cache.DeleteProduct(ctx, id); err != nil {
		// 缓存删除失败只记录日志，不影响主流程
		// log.Printf("Delete product cache failed: %v", err)
	}
	
	return nil
}

// GetSKUByID 根据ID获取SKU
func (r *ProductRepositoryImpl) GetSKUByID(ctx context.Context, id int64) (*entity.ProductSKU, error) {
	var sku entity.ProductSKU
	result := r.db.WithContext(ctx).First(&sku, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return &sku, nil
}

// GetSKUsByProductID 获取商品的SKU列表
func (r *ProductRepositoryImpl) GetSKUsByProductID(ctx context.Context, productID int64) ([]*entity.ProductSKU, error) {
	var skus []*entity.ProductSKU
	err := r.db.WithContext(ctx).Where("goods = ?", productID).Find(&skus).Error
	return skus, err
}

// CreateSKU 创建SKU
func (r *ProductRepositoryImpl) CreateSKU(ctx context.Context, sku *entity.ProductSKU) error {
	return r.db.WithContext(ctx).Create(sku).Error
}

// UpdateSKU 更新SKU
func (r *ProductRepositoryImpl) UpdateSKU(ctx context.Context, sku *entity.ProductSKU) error {
	return r.db.WithContext(ctx).Save(sku).Error
}

// DeleteSKU 删除SKU
func (r *ProductRepositoryImpl) DeleteSKU(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.ProductSKU{}, id).Error
}

// GetAttributesByProductID 获取商品属性
func (r *ProductRepositoryImpl) GetAttributesByProductID(ctx context.Context, productID int64) ([]*entity.ProductAttribute, error) {
	var attrs []*entity.ProductAttribute
	err := r.db.WithContext(ctx).Where("goods = ?", productID).Order("attr_sort").Find(&attrs).Error
	return attrs, err
}

// SaveAttributes 保存商品属性
func (r *ProductRepositoryImpl) SaveAttributes(ctx context.Context, attributes []*entity.ProductAttribute) error {
	return r.db.WithContext(ctx).Create(&attributes).Error
}

// DeleteAttributeByProductID 删除商品属性
func (r *ProductRepositoryImpl) DeleteAttributeByProductID(ctx context.Context, productID int64) error {
	return r.db.WithContext(ctx).Where("goods = ?", productID).Delete(&entity.ProductAttribute{}).Error
}

// GetImagesByProductID 获取商品图片
func (r *ProductRepositoryImpl) GetImagesByProductID(ctx context.Context, productID int64) ([]*entity.ProductImage, error) {
	var images []*entity.ProductImage
	err := r.db.WithContext(ctx).Where("goods = ?", productID).Order("sort").Find(&images).Error
	return images, err
}

// SaveImages 保存商品图片
func (r *ProductRepositoryImpl) SaveImages(ctx context.Context, images []*entity.ProductImage) error {
	return r.db.WithContext(ctx).Create(&images).Error
}

// DeleteImageByProductID 删除商品图片
func (r *ProductRepositoryImpl) DeleteImageByProductID(ctx context.Context, productID int64) error {
	return r.db.WithContext(ctx).Where("goods = ?", productID).Delete(&entity.ProductImage{}).Error
}

// GetSpecsByProductID 获取商品规格
func (r *ProductRepositoryImpl) GetSpecsByProductID(ctx context.Context, productID int64) ([]*entity.ProductSpec, error) {
	var specs []*entity.ProductSpec
	err := r.db.WithContext(ctx).Where("goods = ?", productID).Find(&specs).Error
	return specs, err
}

// SaveSpecs 保存商品规格
func (r *ProductRepositoryImpl) SaveSpecs(ctx context.Context, specs []*entity.ProductSpec) error {
	return r.db.WithContext(ctx).Create(&specs).Error
}

// DeleteSpecByProductID 删除商品规格
func (r *ProductRepositoryImpl) DeleteSpecByProductID(ctx context.Context, productID int64) error {
	return r.db.WithContext(ctx).Where("goods = ?", productID).Delete(&entity.ProductSpec{}).Error
}
