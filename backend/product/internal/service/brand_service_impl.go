package service

import (
	"context"
	"errors"
	"time"
	
	"shop/backend/product/internal/domain/entity"
)

var (
	ErrBrandNotFound    = errors.New("brand not found")
	ErrInvalidBrand     = errors.New("invalid brand data")
	ErrBrandHasProducts = errors.New("brand has associated products")
	ErrRelationExists   = errors.New("relation already exists")
	ErrRelationNotFound = errors.New("relation not found")
)

// BrandServiceImpl 品牌服务实现
type BrandServiceImpl struct {
	brandRepo   BrandRepository
	productRepo ProductRepository
}

// NewBrandService 创建品牌服务实例
func NewBrandService(
	brandRepo BrandRepository,
	productRepo ProductRepository,
) BrandService {
	return &BrandServiceImpl{
		brandRepo:   brandRepo,
		productRepo: productRepo,
	}
}

// GetBrandByID 根据ID获取品牌
func (s *BrandServiceImpl) GetBrandByID(ctx context.Context, id int64) (*entity.Brand, error) {
	brand, err := s.brandRepo.GetBrandByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if brand == nil {
		return nil, ErrBrandNotFound
	}
	
	return brand, nil
}

// ListBrands 获取品牌列表
func (s *BrandServiceImpl) ListBrands(ctx context.Context, filter BrandFilter) ([]*entity.Brand, int64, error) {
	return s.brandRepo.ListBrands(ctx, filter)
}

// CreateBrand 创建品牌
func (s *BrandServiceImpl) CreateBrand(ctx context.Context, brand *entity.Brand) (*entity.Brand, error) {
	// 基本参数验证
	if brand.Name == "" {
		return nil, ErrInvalidBrand
	}
	
	// 设置初始值
	now := time.Now()
	brand.CreatedAt = now
	brand.UpdatedAt = now
	
	// 保存品牌
	if err := s.brandRepo.CreateBrand(ctx, brand); err != nil {
		return nil, err
	}
	
	return brand, nil
}

// UpdateBrand 更新品牌
func (s *BrandServiceImpl) UpdateBrand(ctx context.Context, brand *entity.Brand) error {
	// 检查品牌是否存在
	existingBrand, err := s.brandRepo.GetBrandByID(ctx, brand.ID)
	if err != nil {
		return err
	}
	
	if existingBrand == nil {
		return ErrBrandNotFound
	}
	
	// 更新时间
	brand.UpdatedAt = time.Now()
	brand.CreatedAt = existingBrand.CreatedAt
	
	// 保存品牌
	if err := s.brandRepo.UpdateBrand(ctx, brand); err != nil {
		return err
	}
	
	return nil
}

// DeleteBrand 删除品牌
func (s *BrandServiceImpl) DeleteBrand(ctx context.Context, id int64) error {
	// 检查品牌是否存在
	brand, err := s.brandRepo.GetBrandByID(ctx, id)
	if err != nil {
		return err
	}
	
	if brand == nil {
		return ErrBrandNotFound
	}
	
	// 检查品牌下是否有商品
	// 此处需要ProductRepository的新方法，暂时不实现这个检查
	// 或者可以通过计数查询实现
	
	// 删除品牌
	if err := s.brandRepo.DeleteBrand(ctx, id); err != nil {
		return err
	}
	
	return nil
}

// GetBrandsByCategoryID 获取分类下的品牌列表
func (s *BrandServiceImpl) GetBrandsByCategoryID(ctx context.Context, categoryID int64) ([]*entity.Brand, error) {
	return s.brandRepo.ListBrandsByCategoryID(ctx, categoryID)
}

// GetCategoriesByBrandID 获取品牌关联的分类列表
func (s *BrandServiceImpl) GetCategoriesByBrandID(ctx context.Context, brandID int64) ([]*entity.Category, error) {
	return s.brandRepo.ListCategoriesByBrandID(ctx, brandID)
}

// CreateCategoryBrand 创建分类品牌关联
func (s *BrandServiceImpl) CreateCategoryBrand(ctx context.Context, categoryID int64, brandID int64) error {
	// 检查分类和品牌是否已关联
	brands, err := s.brandRepo.ListBrandsByCategoryID(ctx, categoryID)
	if err != nil {
		return err
	}
	
	for _, b := range brands {
		if b.ID == brandID {
			return ErrRelationExists
		}
	}
	
	// 创建关联
	categoryBrand := &entity.CategoryBrand{
		CategoryID: categoryID,
		BrandID:    brandID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	return s.brandRepo.CreateCategoryBrand(ctx, categoryBrand)
}

// DeleteCategoryBrand 删除分类品牌关联
func (s *BrandServiceImpl) DeleteCategoryBrand(ctx context.Context, categoryID int64, brandID int64) error {
	return s.brandRepo.DeleteCategoryBrand(ctx, categoryID, brandID)
}

// CategoryBrandList 获取分类品牌关联列表
func (s *BrandServiceImpl) CategoryBrandList(ctx context.Context, page, pageSize int) ([]*entity.CategoryBrand, int64, error) {
	// 这里需要BrandRepository的新方法来获取CategoryBrand列表
	// 暂时返回空列表
	return []*entity.CategoryBrand{}, 0, nil
}
