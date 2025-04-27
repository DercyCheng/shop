package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shop/product/internal/domain/entity"
	"shop/product/internal/repository"
)

// BannerServiceImpl implements BannerService interface
type BannerServiceImpl struct {
	bannerRepo repository.BannerRepository
}

// NewBannerService creates a new banner service
func NewBannerService(bannerRepo repository.BannerRepository) BannerService {
	return &BannerServiceImpl{
		bannerRepo: bannerRepo,
	}
}

// GetBannerByID retrieves a banner by ID
func (s *BannerServiceImpl) GetBannerByID(ctx context.Context, id int64) (*entity.Banner, error) {
	return s.bannerRepo.GetByID(ctx, id)
}

// ListBanners retrieves all banners ordered by index
func (s *BannerServiceImpl) ListBanners(ctx context.Context, limit int) ([]*entity.Banner, error) {
	return s.bannerRepo.List(ctx, limit)
}

// CreateBanner adds a new banner
func (s *BannerServiceImpl) CreateBanner(ctx context.Context, banner *entity.Banner) (*entity.Banner, error) {
	// Validate required fields
	if banner.Image == "" {
		return nil, errors.New("banner image is required")
	}

	// Set timestamps
	now := time.Now()
	banner.CreatedAt = now
	banner.UpdatedAt = now

	// Create banner
	if err := s.bannerRepo.Create(ctx, banner); err != nil {
		return nil, err
	}

	return banner, nil
}

// UpdateBanner modifies an existing banner
func (s *BannerServiceImpl) UpdateBanner(ctx context.Context, banner *entity.Banner) error {
	// Verify banner exists
	existingBanner, err := s.bannerRepo.GetByID(ctx, banner.ID)
	if err != nil {
		return fmt.Errorf("banner not found: %v", err)
	}

	// Keep existing values if not provided
	if banner.Image == "" {
		banner.Image = existingBanner.Image
	}

	if banner.URL == "" {
		banner.URL = existingBanner.URL
	}

	// If index not specified, keep existing
	if banner.Index == 0 {
		banner.Index = existingBanner.Index
	}

	// Update timestamp
	banner.UpdatedAt = time.Now()

	return s.bannerRepo.Update(ctx, banner)
}

// DeleteBanner removes a banner by ID
func (s *BannerServiceImpl) DeleteBanner(ctx context.Context, id int64) error {
	return s.bannerRepo.Delete(ctx, id)
}
