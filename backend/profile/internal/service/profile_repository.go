package service

import (
	"context"
	
	"shop/backend/profile/internal/domain/entity"
)

// ProfileRepository 个人信息仓储接口
type ProfileRepository interface {
	// 收藏相关
	GetFavoritesByUserID(ctx context.Context, userID int64, offset, limit int) ([]*entity.UserFav, int64, error)
	GetFavoriteByID(ctx context.Context, id int64) (*entity.UserFav, error)
	CreateFavorite(ctx context.Context, fav *entity.UserFav) error
	DeleteFavorite(ctx context.Context, id int64) error
	IsFavorite(ctx context.Context, userID, goodsID int64) (bool, error)
	
	// 地址相关
	GetAddressesByUserID(ctx context.Context, userID int64) ([]*entity.Address, error)
	GetAddressByID(ctx context.Context, id int64) (*entity.Address, error)
	CreateAddress(ctx context.Context, address *entity.Address) error
	UpdateAddress(ctx context.Context, address *entity.Address) error
	DeleteAddress(ctx context.Context, id int64) error
	SetDefaultAddress(ctx context.Context, userID, addressID int64) error
	GetDefaultAddress(ctx context.Context, userID int64) (*entity.Address, error)
	
	// 用户反馈相关
	GetFeedbacksByUserID(ctx context.Context, userID int64, offset, limit int) ([]*entity.UserFeedback, int64, error)
	GetFeedbackByID(ctx context.Context, id int64) (*entity.UserFeedback, error)
	CreateFeedback(ctx context.Context, feedback *entity.UserFeedback) error
	UpdateFeedback(ctx context.Context, feedback *entity.UserFeedback) error
	DeleteFeedback(ctx context.Context, id int64) error
	
	// 浏览历史相关
	GetBrowsingHistories(ctx context.Context, userID int64, offset, limit int) ([]*entity.BrowsingHistory, int64, error)
	AddBrowsingHistory(ctx context.Context, history *entity.BrowsingHistory) error
	DeleteBrowsingHistory(ctx context.Context, userID int64, ids []int64) error
	ClearBrowsingHistory(ctx context.Context, userID int64) error
}
