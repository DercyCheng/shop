package service

import (
	"context"
	"errors"
	"time"
	
	"shop/backend/product/internal/domain/entity"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product data")
	ErrDuplicateSKU    = errors.New("duplicate SKU code")
	ErrSKUNotFound     = errors.New("SKU not found")
)

// ProductServiceImpl 商品服务实现
type ProductServiceImpl struct {
	productRepo ProductRepository
	categoryRepo CategoryRepository
	brandRepo    BrandRepository
	searchRepo   SearchRepository
}

// NewProductService 创建商品服务实例
func NewProductService(
	productRepo ProductRepository,
	categoryRepo CategoryRepository,
	brandRepo BrandRepository,
	searchRepo SearchRepository,
) ProductService {
	return &ProductServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
		searchRepo:   searchRepo,
	}
}

// GetProductByID 根据ID获取商品
func (s *ProductServiceImpl) GetProductByID(ctx context.Context, id int64) (*entity.Product, error) {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if product == nil {
		return nil, ErrProductNotFound
	}
	
	// 获取商品类别
	if product.CategoryID > 0 {
		category, err := s.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
		if err == nil && category != nil {
			product.Category = category
		}
	}
	
	// 获取商品品牌
	if product.BrandsID > 0 {
		brand, err := s.brandRepo.GetBrandByID(ctx, product.BrandsID)
		if err == nil && brand != nil {
			product.Brand = brand
		}
	}
	
	// 获取商品SKU列表
	skus, err := s.productRepo.GetSKUsByProductID(ctx, id)
	if err == nil {
		product.SkuList = skus
	}
	
	// 增加点击次数
	go s.incrementClickCount(context.Background(), id)
	
	return product, nil
}

// incrementClickCount 异步增加点击次数
func (s *ProductServiceImpl) incrementClickCount(ctx context.Context, id int64) {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil || product == nil {
		return
	}
	
	product.IncreaseClickNum()
	s.productRepo.UpdateProduct(ctx, product)
	
	// 更新搜索索引
	s.searchRepo.IndexProduct(ctx, product)
}

// GetProductBySN 根据商品编号获取商品
func (s *ProductServiceImpl) GetProductBySN(ctx context.Context, goodsSN string) (*entity.Product, error) {
	product, err := s.productRepo.GetProductBySN(ctx, goodsSN)
	if err != nil {
		return nil, err
	}
	
	if product == nil {
		return nil, ErrProductNotFound
	}
	
	return product, nil
}

// ListProducts 获取商品列表
func (s *ProductServiceImpl) ListProducts(ctx context.Context, filter ProductFilter) ([]*entity.Product, int64, error) {
	return s.productRepo.ListProducts(ctx, filter)
}

// BatchGetProducts 批量获取商品
func (s *ProductServiceImpl) BatchGetProducts(ctx context.Context, ids []int64) ([]*entity.Product, error) {
	return s.productRepo.BatchGetProducts(ctx, ids)
}

// CreateProduct 创建商品
func (s *ProductServiceImpl) CreateProduct(
	ctx context.Context,
	product *entity.Product,
	skus []*entity.ProductSKU,
	attrs []*entity.ProductAttribute,
	specs []*entity.ProductSpec,
	images []string,
) (*entity.Product, error) {
	// 基本参数验证
	if product.Name == "" || product.CategoryID <= 0 || product.BrandsID <= 0 {
		return nil, ErrInvalidProduct
	}
	
	// 检查分类是否存在
	category, err := s.categoryRepo.GetCategoryByID(ctx, product.CategoryID)
	if err != nil || category == nil {
		return nil, errors.New("category not found")
	}
	
	// 检查品牌是否存在
	brand, err := s.brandRepo.GetBrandByID(ctx, product.BrandsID)
	if err != nil || brand == nil {
		return nil, errors.New("brand not found")
	}
	
	// 设置初始值
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now
	
	// 保存商品基本信息
	if err := s.productRepo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}
	
	// 保存SKU信息
	if len(skus) > 0 {
		for _, sku := range skus {
			sku.ProductID = product.ID
			sku.CreatedAt = now
			sku.UpdatedAt = now
			if err := s.productRepo.CreateSKU(ctx, sku); err != nil {
				return nil, err
			}
		}
	}
	
	// 保存属性信息
	if len(attrs) > 0 {
		for i := range attrs {
			attrs[i].ProductID = product.ID
		}
		if err := s.productRepo.SaveAttributes(ctx, attrs); err != nil {
			return nil, err
		}
	}
	
	// 保存规格信息
	if len(specs) > 0 {
		for i := range specs {
			specs[i].ProductID = product.ID
		}
		if err := s.productRepo.SaveSpecs(ctx, specs); err != nil {
			return nil, err
		}
	}
	
	// 保存图片信息
	if len(images) > 0 {
		productImages := make([]*entity.ProductImage, len(images))
		for i, image := range images {
			productImages[i] = &entity.ProductImage{
				ProductID: product.ID,
				ImageURL:  image,
				IsMain:    i == 0,
				Sort:      i,
				CreatedAt: now,
				UpdatedAt: now,
			}
		}
		if err := s.productRepo.SaveImages(ctx, productImages); err != nil {
			return nil, err
		}
	}
	
	// 同步到搜索引擎
	go s.searchRepo.IndexProduct(context.Background(), product)
	
	return product, nil
}

// UpdateProduct 更新商品
func (s *ProductServiceImpl) UpdateProduct(
	ctx context.Context,
	product *entity.Product,
	skus []*entity.ProductSKU,
	attrs []*entity.ProductAttribute,
	specs []*entity.ProductSpec,
	images []string,
) error {
	// 检查商品是否存在
	existingProduct, err := s.productRepo.GetProductByID(ctx, product.ID)
	if err != nil {
		return err
	}
	
	if existingProduct == nil {
		return ErrProductNotFound
	}
	
	// 更新时间
	product.UpdatedAt = time.Now()
	product.CreatedAt = existingProduct.CreatedAt
	
	// 保存商品基本信息
	if err := s.productRepo.UpdateProduct(ctx, product); err != nil {
		return err
	}
	
	// 更新SKU信息
	if len(skus) > 0 {
		// 删除旧SKU
		oldSkus, err := s.productRepo.GetSKUsByProductID(ctx, product.ID)
		if err == nil {
			for _, oldSku := range oldSkus {
				// 检查是否保留
				found := false
				for _, newSku := range skus {
					if newSku.ID == oldSku.ID {
						found = true
						break
					}
				}
				if !found {
					s.productRepo.DeleteSKU(ctx, oldSku.ID)
				}
			}
		}
		
		// 添加或更新SKU
		now := time.Now()
		for _, sku := range skus {
			sku.ProductID = product.ID
			sku.UpdatedAt = now
			
			if sku.ID > 0 {
				// 更新现有SKU
				s.productRepo.UpdateSKU(ctx, sku)
			} else {
				// 添加新SKU
				sku.CreatedAt = now
				s.productRepo.CreateSKU(ctx, sku)
			}
		}
	}
	
	// 更新属性
	if len(attrs) > 0 {
		s.productRepo.DeleteAttributeByProductID(ctx, product.ID)
		for i := range attrs {
			attrs[i].ProductID = product.ID
		}
		s.productRepo.SaveAttributes(ctx, attrs)
	}
	
	// 更新规格
	if len(specs) > 0 {
		s.productRepo.DeleteSpecByProductID(ctx, product.ID)
		for i := range specs {
			specs[i].ProductID = product.ID
		}
		s.productRepo.SaveSpecs(ctx, specs)
	}
	
	// 更新图片
	if len(images) > 0 {
		s.productRepo.DeleteImageByProductID(ctx, product.ID)
		now := time.Now()
		productImages := make([]*entity.ProductImage, len(images))
		for i, image := range images {
			productImages[i] = &entity.ProductImage{
				ProductID: product.ID,
				ImageURL:  image,
				IsMain:    i == 0,
				Sort:      i,
				CreatedAt: now,
				UpdatedAt: now,
			}
		}
		if err := s.productRepo.SaveImages(ctx, productImages); err != nil {
			return err
		}
	}
	
	// 同步到搜索引擎
	go s.searchRepo.IndexProduct(context.Background(), product)
	
	return nil
}

// DeleteProduct 删除商品
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id int64) error {
	// 检查商品是否存在
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return err
	}
	
	if product == nil {
		return ErrProductNotFound
	}
	
	// 删除商品
	if err := s.productRepo.DeleteProduct(ctx, id); err != nil {
		return err
	}
	
	// 删除关联的SKU、属性、规格、图片等
	// 这些可能在数据库层通过外键cascade删除，但仍然保留这些调用以确保数据完整性
	s.productRepo.DeleteAttributeByProductID(ctx, id)
	s.productRepo.DeleteSpecByProductID(ctx, id)
	s.productRepo.DeleteImageByProductID(ctx, id)
	
	// 从搜索引擎删除
	go s.searchRepo.DeleteProductIndex(context.Background(), id)
	
	return nil
}

// GetSKUByID 根据ID获取SKU
func (s *ProductServiceImpl) GetSKUByID(ctx context.Context, id int64) (*entity.ProductSKU, error) {
	return s.productRepo.GetSKUByID(ctx, id)
}

// GetSKUsByProductID 获取商品的SKU列表
func (s *ProductServiceImpl) GetSKUsByProductID(ctx context.Context, productID int64) ([]*entity.ProductSKU, error) {
	return s.productRepo.GetSKUsByProductID(ctx, productID)
}

// UpdateSKU 更新SKU信息
func (s *ProductServiceImpl) UpdateSKU(ctx context.Context, sku *entity.ProductSKU) error {
	return s.productRepo.UpdateSKU(ctx, sku)
}

// SetOnSale 设置商品上下架状态
func (s *ProductServiceImpl) SetOnSale(ctx context.Context, id int64, onSale bool) error {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return err
	}
	
	if product == nil {
		return ErrProductNotFound
	}
	
	product.SetOnSale(onSale)
	if err := s.productRepo.UpdateProduct(ctx, product); err != nil {
		return err
	}
	
	// 更新搜索索引
	go s.searchRepo.IndexProduct(context.Background(), product)
	
	return nil
}

// RecordClick 记录商品点击
func (s *ProductServiceImpl) RecordClick(ctx context.Context, id int64) error {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return err
	}
	
	if product == nil {
		return ErrProductNotFound
	}
	
	product.IncreaseClickNum()
	if err := s.productRepo.UpdateProduct(ctx, product); err != nil {
		return err
	}
	
	// 更新搜索索引
	go s.searchRepo.IndexProduct(context.Background(), product)
	
	return nil
}

// UpdateSoldCount 更新销售数量
func (s *ProductServiceImpl) UpdateSoldCount(ctx context.Context, id int64, count int) error {
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return err
	}
	
	if product == nil {
		return ErrProductNotFound
	}
	
	product.IncreaseSoldNum(count)
	if err := s.productRepo.UpdateProduct(ctx, product); err != nil {
		return err
	}
	
	// 更新搜索索引
	go s.searchRepo.IndexProduct(context.Background(), product)
	
	return nil
}
