package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"shop/backend/profile/internal/domain/entity"
	"shop/backend/profile/internal/repository"
	"shop/backend/profile/pkg/client"
)

// UserFavServiceImpl implements the UserFavService interface
type UserFavServiceImpl struct {
	userFavRepo   repository.UserFavRepository
	productClient client.ProductClient
	logger        *zap.Logger
}

// NewUserFavService creates a new UserFavService implementation
func NewUserFavService(
	userFavRepo repository.UserFavRepository,
	productClient client.ProductClient,
	logger *zap.Logger,
) UserFavService {
	return &UserFavServiceImpl{
		userFavRepo:   userFavRepo,
		productClient: productClient,
		logger:        logger,
	}
}

// AddUserFav adds a favorite for a user
func (s *UserFavServiceImpl) AddUserFav(ctx context.Context, userID, goodsID int64) error {
	// Check if the product exists
	exists, err := s.productClient.CheckProductExists(ctx, goodsID)
	if err != nil {
		s.logger.Error("Failed to check product existence", zap.Error(err))
		return fmt.Errorf("failed to check product existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("product with ID %d does not exist", goodsID)
	}

	// Check if user already has this favorite
	isFav, err := s.userFavRepo.CheckUserFav(ctx, userID, goodsID)
	if err != nil {
		s.logger.Error("Failed to check user favorite", zap.Error(err))
		return fmt.Errorf("failed to check user favorite: %w", err)
	}

	if isFav {
		return fmt.Errorf("user already favorited this product")
	}

	// Add favorite
	if err := s.userFavRepo.AddUserFav(ctx, userID, goodsID); err != nil {
		s.logger.Error("Failed to add user favorite",
			zap.Int64("user_id", userID),
			zap.Int64("goods_id", goodsID),
			zap.Error(err))
		return fmt.Errorf("failed to add user favorite: %w", err)
	}

	s.logger.Info("Added user favorite",
		zap.Int64("user_id", userID),
		zap.Int64("goods_id", goodsID))
	return nil
}

// DeleteUserFav removes a favorite for a user
func (s *UserFavServiceImpl) DeleteUserFav(ctx context.Context, userID, goodsID int64) error {
	// Check if user has this favorite
	isFav, err := s.userFavRepo.CheckUserFav(ctx, userID, goodsID)
	if err != nil {
		s.logger.Error("Failed to check user favorite", zap.Error(err))
		return fmt.Errorf("failed to check user favorite: %w", err)
	}

	if !isFav {
		return fmt.Errorf("user hasn't favorited this product")
	}

	// Delete favorite
	if err := s.userFavRepo.DeleteUserFav(ctx, userID, goodsID); err != nil {
		s.logger.Error("Failed to delete user favorite",
			zap.Int64("user_id", userID),
			zap.Int64("goods_id", goodsID),
			zap.Error(err))
		return fmt.Errorf("failed to delete user favorite: %w", err)
	}

	s.logger.Info("Deleted user favorite",
		zap.Int64("user_id", userID),
		zap.Int64("goods_id", goodsID))
	return nil
}

// CheckUserFav checks if a user has favorited a goods
func (s *UserFavServiceImpl) CheckUserFav(ctx context.Context, userID, goodsID int64) (bool, error) {
	return s.userFavRepo.CheckUserFav(ctx, userID, goodsID)
}

// GetUserFavList retrieves a list of user favorites with pagination and product details
func (s *UserFavServiceImpl) GetUserFavList(ctx context.Context, userID int64, page, pageSize int) ([]*entity.UserFavWithGoods, int64, error) {
	// Get favorites from repository
	favs, total, err := s.userFavRepo.GetUserFavList(ctx, userID, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to get user favorite list",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get user favorite list: %w", err)
	}

	// No favorites found
	if len(favs) == 0 {
		return []*entity.UserFavWithGoods{}, total, nil
	}

	// Collect product IDs
	goodsIDs := make([]int64, 0, len(favs))
	for _, fav := range favs {
		goodsIDs = append(goodsIDs, fav.GoodsID)
	}

	// Get product details from product service
	productsInfo, err := s.productClient.GetProductsByIDs(ctx, goodsIDs)
	if err != nil {
		s.logger.Error("Failed to get product details",
			zap.Int64s("goods_ids", goodsIDs),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get product details: %w", err)
	}

	// Map product info to favorites
	result := make([]*entity.UserFavWithGoods, 0, len(favs))
	for _, fav := range favs {
		favWithGoods := &entity.UserFavWithGoods{
			UserFav:   *fav,
			GoodsInfo: nil,
		}

		// Find product info
		if info, ok := productsInfo[fav.GoodsID]; ok {
			favWithGoods.GoodsInfo = info
		}

		result = append(result, favWithGoods)
	}

	return result, total, nil
}
