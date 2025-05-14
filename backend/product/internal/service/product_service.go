package service

import (
	"context"
	
	"shop/backend/product/internal/domain/entity"
)

// ProductService 商品服务接口
type ProductService interface {
	// 商品管理相关接口
	GetProductByID(ctx context.Context, id int64) (*entity.Product, error)
	GetProductBySN(ctx context.Context, goodsSN string) (*entity.Product, error)
	ListProducts(ctx context.Context, filter ProductFilter) ([]*entity.Product, int64, error)
	BatchGetProducts(ctx context.Context, ids []int64) ([]*entity.Product, error)
	CreateProduct(ctx context.Context, product *entity.Product, skus []*entity.ProductSKU, attrs []*entity.ProductAttribute, specs []*entity.ProductSpec, images []string) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product, skus []*entity.ProductSKU, attrs []*entity.ProductAttribute, specs []*entity.ProductSpec, images []string) error
	DeleteProduct(ctx context.Context, id int64) error
	
	// 商品SKU相关接口
	GetSKUByID(ctx context.Context, id int64) (*entity.ProductSKU, error)
	GetSKUsByProductID(ctx context.Context, productID int64) ([]*entity.ProductSKU, error)
	UpdateSKU(ctx context.Context, sku *entity.ProductSKU) error
	
	// 商品状态变更接口
	SetOnSale(ctx context.Context, id int64, onSale bool) error
	RecordClick(ctx context.Context, id int64) error
	UpdateSoldCount(ctx context.Context, id int64, count int) error
}

// CategoryService 分类服务接口
type CategoryService interface {
	// 分类管理相关接口
	GetCategoryByID(ctx context.Context, id int64) (*entity.Category, error)
	GetAllCategories(ctx context.Context) ([]*entity.Category, error)
	GetCategoryTree(ctx context.Context) ([]*entity.Category, error)
	GetSubCategories(ctx context.Context, parentID int64) ([]*entity.Category, error)
	CreateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

// BrandService 品牌服务接口
type BrandService interface {
	// 品牌管理相关接口
	GetBrandByID(ctx context.Context, id int64) (*entity.Brand, error)
	ListBrands(ctx context.Context, filter BrandFilter) ([]*entity.Brand, int64, error)
	CreateBrand(ctx context.Context, brand *entity.Brand) (*entity.Brand, error)
	UpdateBrand(ctx context.Context, brand *entity.Brand) error
	DeleteBrand(ctx context.Context, id int64) error
	
	// 分类品牌关联接口
	GetBrandsByCategoryID(ctx context.Context, categoryID int64) ([]*entity.Brand, error)
	GetCategoriesByBrandID(ctx context.Context, brandID int64) ([]*entity.Category, error)
	CreateCategoryBrand(ctx context.Context, categoryID int64, brandID int64) error
	DeleteCategoryBrand(ctx context.Context, categoryID int64, brandID int64) error
	CategoryBrandList(ctx context.Context, page, pageSize int) ([]*entity.CategoryBrand, int64, error)
}

// BannerService 轮播图服务接口
type BannerService interface {
	// 轮播图管理相关接口
	GetBannerByID(ctx context.Context, id int64) (*entity.Banner, error) 
	ListBanners(ctx context.Context) ([]*entity.Banner, error)
	CreateBanner(ctx context.Context, banner *entity.Banner) (*entity.Banner, error)
	UpdateBanner(ctx context.Context, banner *entity.Banner) error
	DeleteBanner(ctx context.Context, id int64) error
}

// SearchService 搜索服务接口
type SearchService interface {
	// 搜索相关接口
	SearchProducts(ctx context.Context, params *SearchParams) (*SearchResult, error)
	IndexProduct(ctx context.Context, product *entity.Product) error
	BatchIndexProducts(ctx context.Context, products []*entity.Product) error
	DeleteProductIndex(ctx context.Context, id int64) error
	SyncProductIndex(ctx context.Context) error
}
