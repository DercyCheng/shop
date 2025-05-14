package service

import (
	"context"
	"errors"
	"time"
	
	"shop/backend/product/internal/domain/entity"
)

var (
	ErrBannerNotFound = errors.New("banner not found")
	ErrInvalidBanner  = errors.New("invalid banner data")
)

// BannerServiceImpl 轮播图服务实现
type BannerServiceImpl struct {
	bannerRepo BannerRepository
}

// NewBannerService 创建轮播图服务实例
func NewBannerService(bannerRepo BannerRepository) BannerService {
	return &BannerServiceImpl{
		bannerRepo: bannerRepo,
	}
}

// GetBannerByID 根据ID获取轮播图
func (s *BannerServiceImpl) GetBannerByID(ctx context.Context, id int64) (*entity.Banner, error) {
	banner, err := s.bannerRepo.GetBannerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if banner == nil {
		return nil, ErrBannerNotFound
	}
	
	return banner, nil
}

// ListBanners 获取轮播图列表
func (s *BannerServiceImpl) ListBanners(ctx context.Context) ([]*entity.Banner, error) {
	return s.bannerRepo.ListBanners(ctx)
}

// CreateBanner 创建轮播图
func (s *BannerServiceImpl) CreateBanner(ctx context.Context, banner *entity.Banner) (*entity.Banner, error) {
	// 基本参数验证
	if banner.Image == "" {
		return nil, ErrInvalidBanner
	}
	
	// 设置初始值
	now := time.Now()
	banner.CreatedAt = now
	banner.UpdatedAt = now
	
	// 保存轮播图
	if err := s.bannerRepo.CreateBanner(ctx, banner); err != nil {
		return nil, err
	}
	
	return banner, nil
}

// UpdateBanner 更新轮播图
func (s *BannerServiceImpl) UpdateBanner(ctx context.Context, banner *entity.Banner) error {
	// 检查轮播图是否存在
	existingBanner, err := s.bannerRepo.GetBannerByID(ctx, banner.ID)
	if err != nil {
		return err
	}
	
	if existingBanner == nil {
		return ErrBannerNotFound
	}
	
	// 更新时间
	banner.UpdatedAt = time.Now()
	banner.CreatedAt = existingBanner.CreatedAt
	
	// 保存轮播图
	if err := s.bannerRepo.UpdateBanner(ctx, banner); err != nil {
		return err
	}
	
	return nil
}

// DeleteBanner 删除轮播图
func (s *BannerServiceImpl) DeleteBanner(ctx context.Context, id int64) error {
	// 检查轮播图是否存在
	banner, err := s.bannerRepo.GetBannerByID(ctx, id)
	if err != nil {
		return err
	}
	
	if banner == nil {
		return ErrBannerNotFound
	}
	
	// 删除轮播图
	if err := s.bannerRepo.DeleteBanner(ctx, id); err != nil {
		return err
	}
	
	return nil
}
