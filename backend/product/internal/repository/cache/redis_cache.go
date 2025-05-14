package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"shop/backend/product/internal/domain/entity"
	
	"github.com/go-redis/redis/v8"
)

// Constants for cache keys
const (
	ProductKeyPrefix    = "product:id:"
	CategoryKeyPrefix   = "category:id:"
	BrandKeyPrefix      = "brand:id:"
	BannerKeyPrefix     = "banner:all"
	CategoryTreePrefix  = "category:tree"
	HotProductsPrefix   = "product:hot"
	NewProductsPrefix   = "product:new"
	ProductListPrefix   = "product:list:"
	CategoryBrandsPrefix = "category:brands:"
	
	// Default expiration times
	DefaultExpiration = 24 * time.Hour
	ShortExpiration   = 1 * time.Hour
	LongExpiration    = 7 * 24 * time.Hour
)

// ProductCache 商品缓存接口
type ProductCache interface {
	// 商品相关
	GetProduct(ctx context.Context, id int64) (*entity.Product, error)
	SetProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	
	// 分类相关
	GetCategory(ctx context.Context, id int64) (*entity.Category, error)
	SetCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int64) error
	GetCategoryTree(ctx context.Context) ([]*entity.Category, error)
	SetCategoryTree(ctx context.Context, categories []*entity.Category) error
	DeleteCategoryTree(ctx context.Context) error
	
	// 品牌相关
	GetBrand(ctx context.Context, id int64) (*entity.Brand, error)
	SetBrand(ctx context.Context, brand *entity.Brand) error
	DeleteBrand(ctx context.Context, id int64) error
	
	// 轮播图相关
	GetBanners(ctx context.Context) ([]*entity.Banner, error)
	SetBanners(ctx context.Context, banners []*entity.Banner) error
	DeleteBanners(ctx context.Context) error
	
	// 热门/新品
	GetHotProducts(ctx context.Context) ([]*entity.Product, error)
	SetHotProducts(ctx context.Context, products []*entity.Product) error
	GetNewProducts(ctx context.Context) ([]*entity.Product, error)
	SetNewProducts(ctx context.Context, products []*entity.Product) error
	
	// 分类品牌关联
	GetCategoryBrands(ctx context.Context, categoryID int64) ([]*entity.Brand, error)
	SetCategoryBrands(ctx context.Context, categoryID int64, brands []*entity.Brand) error
	DeleteCategoryBrands(ctx context.Context, categoryID int64) error
}

// RedisProductCache Redis实现的商品缓存
type RedisProductCache struct {
	client *redis.Client
}

// NewRedisProductCache 创建Redis商品缓存实例
func NewRedisProductCache(client *redis.Client) ProductCache {
	return &RedisProductCache{
		client: client,
	}
}

// GetProduct 获取商品缓存
func (c *RedisProductCache) GetProduct(ctx context.Context, id int64) (*entity.Product, error) {
	key := fmt.Sprintf("%s%d", ProductKeyPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var product entity.Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, err
	}
	
	return &product, nil
}

// SetProduct 设置商品缓存
func (c *RedisProductCache) SetProduct(ctx context.Context, product *entity.Product) error {
	key := fmt.Sprintf("%s%d", ProductKeyPrefix, product.ID)
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, key, data, DefaultExpiration).Err()
}

// DeleteProduct 删除商品缓存
func (c *RedisProductCache) DeleteProduct(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s%d", ProductKeyPrefix, id)
	return c.client.Del(ctx, key).Err()
}

// GetCategory 获取分类缓存
func (c *RedisProductCache) GetCategory(ctx context.Context, id int64) (*entity.Category, error) {
	key := fmt.Sprintf("%s%d", CategoryKeyPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var category entity.Category
	if err := json.Unmarshal(data, &category); err != nil {
		return nil, err
	}
	
	return &category, nil
}

// SetCategory 设置分类缓存
func (c *RedisProductCache) SetCategory(ctx context.Context, category *entity.Category) error {
	key := fmt.Sprintf("%s%d", CategoryKeyPrefix, category.ID)
	data, err := json.Marshal(category)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, key, data, DefaultExpiration).Err()
}

// DeleteCategory 删除分类缓存
func (c *RedisProductCache) DeleteCategory(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s%d", CategoryKeyPrefix, id)
	// 同时删除分类树缓存
	if err := c.DeleteCategoryTree(ctx); err != nil {
		return err
	}
	return c.client.Del(ctx, key).Err()
}

// GetCategoryTree 获取分类树缓存
func (c *RedisProductCache) GetCategoryTree(ctx context.Context) ([]*entity.Category, error) {
	key := CategoryTreePrefix
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var categories []*entity.Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return nil, err
	}
	
	return categories, nil
}

// SetCategoryTree 设置分类树缓存
func (c *RedisProductCache) SetCategoryTree(ctx context.Context, categories []*entity.Category) error {
	key := CategoryTreePrefix
	data, err := json.Marshal(categories)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, key, data, DefaultExpiration).Err()
}

// DeleteCategoryTree 删除分类树缓存
func (c *RedisProductCache) DeleteCategoryTree(ctx context.Context) error {
	key := CategoryTreePrefix
	return c.client.Del(ctx, key).Err()
}

// GetBrand 获取品牌缓存
func (c *RedisProductCache) GetBrand(ctx context.Context, id int64) (*entity.Brand, error) {
	key := fmt.Sprintf("%s%d", BrandKeyPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var brand entity.Brand
	if err := json.Unmarshal(data, &brand); err != nil {
		return nil, err
	}
	
	return &brand, nil
}

// SetBrand 设置品牌缓存
func (c *RedisProductCache) SetBrand(ctx context.Context, brand *entity.Brand) error {
	key := fmt.Sprintf("%s%d", BrandKeyPrefix, brand.ID)
	data, err := json.Marshal(brand)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, key, data, DefaultExpiration).Err()
}

// DeleteBrand 删除品牌缓存
func (c *RedisProductCache) DeleteBrand(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s%d", BrandKeyPrefix, id)
	return c.client.Del(ctx, key).Err()
}

// GetBanners 获取轮播图缓存
func (c *RedisProductCache) GetBanners(ctx context.Context) ([]*entity.Banner, error) {
	data, err := c.client.Get(ctx, BannerKeyPrefix).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var banners []*entity.Banner
	if err := json.Unmarshal(data, &banners); err != nil {
		return nil, err
	}
	
	return banners, nil
}

// SetBanners 设置轮播图缓存
func (c *RedisProductCache) SetBanners(ctx context.Context, banners []*entity.Banner) error {
	data, err := json.Marshal(banners)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, BannerKeyPrefix, data, DefaultExpiration).Err()
}

// DeleteBanners 删除轮播图缓存
func (c *RedisProductCache) DeleteBanners(ctx context.Context) error {
	return c.client.Del(ctx, BannerKeyPrefix).Err()
}

// GetHotProducts 获取热门商品缓存
func (c *RedisProductCache) GetHotProducts(ctx context.Context) ([]*entity.Product, error) {
	data, err := c.client.Get(ctx, HotProductsPrefix).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var products []*entity.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}
	
	return products, nil
}

// SetHotProducts 设置热门商品缓存
func (c *RedisProductCache) SetHotProducts(ctx context.Context, products []*entity.Product) error {
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, HotProductsPrefix, data, ShortExpiration).Err()
}

// GetNewProducts 获取新品缓存
func (c *RedisProductCache) GetNewProducts(ctx context.Context) ([]*entity.Product, error) {
	data, err := c.client.Get(ctx, NewProductsPrefix).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var products []*entity.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}
	
	return products, nil
}

// SetNewProducts 设置新品缓存
func (c *RedisProductCache) SetNewProducts(ctx context.Context, products []*entity.Product) error {
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, NewProductsPrefix, data, ShortExpiration).Err()
}

// GetCategoryBrands 获取分类品牌关联缓存
func (c *RedisProductCache) GetCategoryBrands(ctx context.Context, categoryID int64) ([]*entity.Brand, error) {
	key := fmt.Sprintf("%s%d", CategoryBrandsPrefix, categoryID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var brands []*entity.Brand
	if err := json.Unmarshal(data, &brands); err != nil {
		return nil, err
	}
	
	return brands, nil
}

// SetCategoryBrands 设置分类品牌关联缓存
func (c *RedisProductCache) SetCategoryBrands(ctx context.Context, categoryID int64, brands []*entity.Brand) error {
	key := fmt.Sprintf("%s%d", CategoryBrandsPrefix, categoryID)
	data, err := json.Marshal(brands)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, key, data, DefaultExpiration).Err()
}

// DeleteCategoryBrands 删除分类品牌关联缓存
func (c *RedisProductCache) DeleteCategoryBrands(ctx context.Context, categoryID int64) error {
	key := fmt.Sprintf("%s%d", CategoryBrandsPrefix, categoryID)
	return c.client.Del(ctx, key).Err()
}
