package service

import (
	"context"
	"errors"
	"time"
	
	"shop/backend/profile/internal/domain/entity"
)

// ErrFavoriteExists 商品已收藏错误
var ErrFavoriteExists = errors.New("favorite already exists")

// ErrFavoriteNotFound 收藏不存在错误
var ErrFavoriteNotFound = errors.New("favorite not found")

// ErrUserNotMatch 用户不匹配错误
var ErrUserNotMatch = errors.New("user not match")

// FavoriteServiceImpl 收藏服务实现
type FavoriteServiceImpl struct {
	repo ProfileRepository
}

// NewFavoriteService 创建收藏服务实例
func NewFavoriteService(repo ProfileRepository) FavoriteService {
	return &FavoriteServiceImpl{
		repo: repo,
	}
}

// ListFavorites 获取用户收藏列表
func (s *FavoriteServiceImpl) ListFavorites(ctx context.Context, userID int64, page, pageSize int) ([]*entity.UserFav, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.GetFavoritesByUserID(ctx, userID, offset, pageSize)
}

// AddFavorite 添加收藏
func (s *FavoriteServiceImpl) AddFavorite(ctx context.Context, userID, goodsID, categoryID int64, remark string) error {
	// 检查是否已收藏
	exists, err := s.repo.IsFavorite(ctx, userID, goodsID)
	if err != nil {
		return err
	}
	
	if exists {
		return ErrFavoriteExists
	}
	
	// TODO: 获取商品当前价格，需要调用商品服务
	// 这里暂时设置为0，实际项目中应该从商品服务获取
	priceWhenFav := 0.00
	
	// 创建收藏记录
	now := time.Now()
	fav := &entity.UserFav{
		UserID:       userID,
		GoodsID:      goodsID,
		CategoryID:   categoryID,
		Remark:       remark,
		PriceWhenFav: priceWhenFav,
		Notification: false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	
	return s.repo.CreateFavorite(ctx, fav)
}

// RemoveFavorite 删除收藏
func (s *FavoriteServiceImpl) RemoveFavorite(ctx context.Context, id int64) error {
	// 检查收藏是否存在
	fav, err := s.repo.GetFavoriteByID(ctx, id)
	if err != nil {
		return err
	}
	
	if fav == nil {
		return ErrFavoriteNotFound
	}
	
	return s.repo.DeleteFavorite(ctx, id)
}

// IsFavorite 检查商品是否已收藏
func (s *FavoriteServiceImpl) IsFavorite(ctx context.Context, userID, goodsID int64) (bool, error) {
	return s.repo.IsFavorite(ctx, userID, goodsID)
}

// SetPriceNotification 设置收藏价格变动通知
func (s *FavoriteServiceImpl) SetPriceNotification(ctx context.Context, id int64, notify bool) error {
	// 检查收藏是否存在
	fav, err := s.repo.GetFavoriteByID(ctx, id)
	if err != nil {
		return err
	}
	
	if fav == nil {
		return ErrFavoriteNotFound
	}
	
	// 更新通知设置
	fav.Notification = notify
	fav.UpdatedAt = time.Now()
	
	// 实际项目中这里应该有更新收藏的方法，这里简化处理
	// 重新保存整个实体
	return s.repo.CreateFavorite(ctx, fav)
}
