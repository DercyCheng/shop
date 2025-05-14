package repository

import (
	"context"
	"errors"
	
	"shop/backend/product/internal/domain/entity"
	"shop/backend/product/internal/repository/cache"
	"shop/backend/product/internal/service"
	
	"gorm.io/gorm"
)

// BannerRepositoryImpl 轮播图仓储实现
type BannerRepositoryImpl struct {
	db    *gorm.DB
	cache cache.ProductCache
}

// NewBannerRepository 创建轮播图仓储实例
func NewBannerRepository(db *gorm.DB, cache cache.ProductCache) service.BannerRepository {
	return &BannerRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// GetBannerByID 根据ID获取轮播图
func (r *BannerRepositoryImpl) GetBannerByID(ctx context.Context, id int64) (*entity.Banner, error) {
	var banner entity.Banner
	result := r.db.WithContext(ctx).First(&banner, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return &banner, nil
}

// ListBanners 获取轮播图列表
func (r *BannerRepositoryImpl) ListBanners(ctx context.Context) ([]*entity.Banner, error) {
	// 尝试从缓存获取
	banners, err := r.cache.GetBanners(ctx)
	if err == nil && len(banners) > 0 {
		return banners, nil
	}
	
	// 从数据库获取
	var dbBanners []*entity.Banner
	if err := r.db.WithContext(ctx).Order("`index`").Find(&dbBanners).Error; err != nil {
		return nil, err
	}
	
	// 将轮播图列表放入缓存
	if err := r.cache.SetBanners(ctx, dbBanners); err != nil {
		// 缓存失败只记录日志，不影响主流程
		// log.Printf("Cache banners failed: %v", err)
	}
	
	return dbBanners, nil
}

// CreateBanner 创建轮播图
func (r *BannerRepositoryImpl) CreateBanner(ctx context.Context, banner *entity.Banner) error {
	if err := r.db.WithContext(ctx).Create(banner).Error; err != nil {
		return err
	}
	
	// 清除轮播图列表缓存，强制重新加载
	if err := r.cache.DeleteBanners(ctx); err != nil {
		// 缓存删除失败只记录日志，不影响主流程
		// log.Printf("Delete banners cache failed: %v", err)
	}
	
	return nil
}

// UpdateBanner 更新轮播图
func (r *BannerRepositoryImpl) UpdateBanner(ctx context.Context, banner *entity.Banner) error {
	if err := r.db.WithContext(ctx).Save(banner).Error; err != nil {
		return err
	}
	
	// 清除轮播图列表缓存，强制重新加载
	if err := r.cache.DeleteBanners(ctx); err != nil {
		// 缓存删除失败只记录日志，不影响主流程
		// log.Printf("Delete banners cache failed: %v", err)
	}
	
	return nil
}

// DeleteBanner 删除轮播图
func (r *BannerRepositoryImpl) DeleteBanner(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Banner{}, id).Error; err != nil {
		return err
	}
	
	// 清除轮播图列表缓存，强制重新加载
	if err := r.cache.DeleteBanners(ctx); err != nil {
		// 缓存删除失败只记录日志，不影响主流程
		// log.Printf("Delete banners cache failed: %v", err)
	}
	
	return nil
}
