package service

import (
	"context"
	
	"shop/backend/profile/internal/domain/entity"
)

// FavoriteService 收藏服务接口
type FavoriteService interface {
	// 获取用户收藏列表
	ListFavorites(ctx context.Context, userID int64, page, pageSize int) ([]*entity.UserFav, int64, error)
	
	// 添加收藏
	AddFavorite(ctx context.Context, userID, goodsID, categoryID int64, remark string) error
	
	// 删除收藏
	RemoveFavorite(ctx context.Context, id int64) error
	
	// 检查商品是否已收藏
	IsFavorite(ctx context.Context, userID, goodsID int64) (bool, error)
	
	// 设置收藏价格变动通知
	SetPriceNotification(ctx context.Context, id int64, notify bool) error
}

// AddressService 地址服务接口
type AddressService interface {
	// 获取用户地址列表
	ListAddresses(ctx context.Context, userID int64) ([]*entity.Address, error)
	
	// 获取地址详情
	GetAddress(ctx context.Context, id int64) (*entity.Address, error)
	
	// 添加地址
	AddAddress(ctx context.Context, address *entity.Address) error
	
	// 更新地址
	UpdateAddress(ctx context.Context, address *entity.Address) error
	
	// 删除地址
	DeleteAddress(ctx context.Context, id int64) error
	
	// 设置默认地址
	SetDefault(ctx context.Context, userID, addressID int64) error
	
	// 获取默认地址
	GetDefaultAddress(ctx context.Context, userID int64) (*entity.Address, error)
	
	// 使用地址（更新使用次数和最后使用时间）
	UseAddress(ctx context.Context, id int64) error
}

// FeedbackService 用户反馈服务接口
type FeedbackService interface {
	// 获取用户反馈列表
	ListFeedbacks(ctx context.Context, userID int64, page, pageSize int) ([]*entity.UserFeedback, int64, error)
	
	// 获取反馈详情
	GetFeedback(ctx context.Context, id int64) (*entity.UserFeedback, error)
	
	// 提交反馈
	SubmitFeedback(ctx context.Context, feedback *entity.UserFeedback) error
	
	// 更新反馈
	UpdateFeedback(ctx context.Context, feedback *entity.UserFeedback) error
	
	// 删除反馈
	DeleteFeedback(ctx context.Context, id int64) error
}

// BrowsingHistoryService 浏览历史服务接口
type BrowsingHistoryService interface {
	// 获取用户浏览历史
	GetHistories(ctx context.Context, userID int64, page, pageSize int) ([]*entity.BrowsingHistory, int64, error)
	
	// 添加浏览记录
	AddHistory(ctx context.Context, userID, goodsID int64, source string, stayTime int) error
	
	// 删除浏览记录
	RemoveHistories(ctx context.Context, userID int64, ids []int64) error
	
	// 清空浏览记录
	ClearHistories(ctx context.Context, userID int64) error
}
