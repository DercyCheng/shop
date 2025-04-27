package service

import (
	"context"

	"shop/backend/profile/internal/domain/entity"
)

// UserFavService defines the interface for user favorite operations
type UserFavService interface {
	// AddUserFav adds a favorite for a user
	AddUserFav(ctx context.Context, userID, goodsID int64) error

	// DeleteUserFav removes a favorite for a user
	DeleteUserFav(ctx context.Context, userID, goodsID int64) error

	// CheckUserFav checks if a user has favorited a goods
	CheckUserFav(ctx context.Context, userID, goodsID int64) (bool, error)

	// GetUserFavList retrieves a list of user favorites with pagination and product details
	GetUserFavList(ctx context.Context, userID int64, page, pageSize int) ([]*entity.UserFavWithGoods, int64, error)
}
