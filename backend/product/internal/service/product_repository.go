package service

import (
	"context"
	
	"shop/backend/product/internal/domain/entity"
)

// ProductRepository 商品仓储接口
type ProductRepository interface {
	// 商品相关
	GetProductByID(ctx context.Context, id int64) (*entity.Product, error)
	GetProductBySN(ctx context.Context, goodsSN string) (*entity.Product, error)
	ListProducts(ctx context.Context, filter ProductFilter) ([]*entity.Product, int64, error)
	BatchGetProducts(ctx context.Context, ids []int64) ([]*entity.Product, error)
	CreateProduct(ctx context.Context, product *entity.Product) error
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	
	// SKU相关
	GetSKUByID(ctx context.Context, id int64) (*entity.ProductSKU, error)
	GetSKUsByProductID(ctx context.Context, productID int64) ([]*entity.ProductSKU, error)
	CreateSKU(ctx context.Context, sku *entity.ProductSKU) error
	UpdateSKU(ctx context.Context, sku *entity.ProductSKU) error
	DeleteSKU(ctx context.Context, id int64) error
	
	// 属性相关
	GetAttributesByProductID(ctx context.Context, productID int64) ([]*entity.ProductAttribute, error)
	SaveAttributes(ctx context.Context, attributes []*entity.ProductAttribute) error
	DeleteAttributeByProductID(ctx context.Context, productID int64) error
	
	// 图片相关
	GetImagesByProductID(ctx context.Context, productID int64) ([]*entity.ProductImage, error)
	SaveImages(ctx context.Context, images []*entity.ProductImage) error
	DeleteImageByProductID(ctx context.Context, productID int64) error
	
	// 规格相关
	GetSpecsByProductID(ctx context.Context, productID int64) ([]*entity.ProductSpec, error)
	SaveSpecs(ctx context.Context, specs []*entity.ProductSpec) error
	DeleteSpecByProductID(ctx context.Context, productID int64) error
}

// ProductFilter 商品过滤条件
type ProductFilter struct {
	Name       string
	CategoryID int64
	BrandID    int64
	OnSale     *bool
	IsNew      *bool
	IsHot      *bool
	ShipFree   *bool
	PriceMin   *float64
	PriceMax   *float64
	OrderBy    string
	Page       int
	PageSize   int
}

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	GetCategoryByID(ctx context.Context, id int64) (*entity.Category, error)
	ListAllCategories(ctx context.Context) ([]*entity.Category, error)
	ListCategoriesByParentID(ctx context.Context, parentID int64) ([]*entity.Category, error)
	GetCategoryTree(ctx context.Context, rootID int64) (*entity.Category, error)
	CreateCategory(ctx context.Context, category *entity.Category) error
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int64) error
}

// BrandRepository 品牌仓储接口
type BrandRepository interface {
	GetBrandByID(ctx context.Context, id int64) (*entity.Brand, error)
	ListBrands(ctx context.Context, filter BrandFilter) ([]*entity.Brand, int64, error)
	CreateBrand(ctx context.Context, brand *entity.Brand) error
	UpdateBrand(ctx context.Context, brand *entity.Brand) error
	DeleteBrand(ctx context.Context, id int64) error
	
	// 分类品牌关联
	ListBrandsByCategoryID(ctx context.Context, categoryID int64) ([]*entity.Brand, error)
	ListCategoriesByBrandID(ctx context.Context, brandID int64) ([]*entity.Category, error)
	CreateCategoryBrand(ctx context.Context, categoryBrand *entity.CategoryBrand) error
	DeleteCategoryBrand(ctx context.Context, categoryID, brandID int64) error
}

// BrandFilter 品牌过滤条件
type BrandFilter struct {
	Name     string
	Page     int
	PageSize int
}

// BannerRepository 轮播图仓储接口
type BannerRepository interface {
	GetBannerByID(ctx context.Context, id int64) (*entity.Banner, error)
	ListBanners(ctx context.Context) ([]*entity.Banner, error)
	CreateBanner(ctx context.Context, banner *entity.Banner) error
	UpdateBanner(ctx context.Context, banner *entity.Banner) error
	DeleteBanner(ctx context.Context, id int64) error
}

// SearchRepository 商品搜索仓储接口
type SearchRepository interface {
	SearchProducts(ctx context.Context, params SearchParams) (*SearchResult, error)
	IndexProduct(ctx context.Context, product *entity.Product) error
	BatchIndexProducts(ctx context.Context, products []*entity.Product) error
	DeleteProductIndex(ctx context.Context, id int64) error
	SyncProductIndex(ctx context.Context) error
}

// SearchParams 搜索参数
type SearchParams struct {
	Keyword    string
	CategoryID int64
	BrandID    int64
	PriceMin   float64
	PriceMax   float64
	OnSale     bool
	IsNew      bool
	IsHot      bool
	ShipFree   bool
	Page       int
	PageSize   int
	OrderBy    string
}

// SearchResult 搜索结果
type SearchResult struct {
	Total  int64
	Page   int
	Size   int
	Pages  int
	Goods  []*entity.Product
}
