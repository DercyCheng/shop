package repository

import (
	"context"
	"errors"
	"time"
	
	"shop/backend/profile/internal/domain/entity"
	"shop/backend/profile/internal/repository/cache"
	"shop/backend/profile/internal/service"
	
	"gorm.io/gorm"
)

// ProfileRepositoryImpl 个人信息仓储实现
type ProfileRepositoryImpl struct {
	db    *gorm.DB
	cache cache.ProfileCache
}

// NewProfileRepository 创建个人信息仓储实例
func NewProfileRepository(db *gorm.DB, cache cache.ProfileCache) service.ProfileRepository {
	return &ProfileRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// ======= 收藏相关 =======

// GetFavoritesByUserID 获取用户收藏列表
func (r *ProfileRepositoryImpl) GetFavoritesByUserID(ctx context.Context, userID int64, offset, limit int) ([]*entity.UserFav, int64, error) {
	var favorites []*entity.UserFav
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&entity.UserFav{}).
		Where("user = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页获取记录
	if err := r.db.WithContext(ctx).
		Where("user = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&favorites).Error; err != nil {
		return nil, 0, err
	}
	
	return favorites, total, nil
}

// GetFavoriteByID 获取收藏详情
func (r *ProfileRepositoryImpl) GetFavoriteByID(ctx context.Context, id int64) (*entity.UserFav, error) {
	var fav entity.UserFav
	result := r.db.WithContext(ctx).First(&fav, id)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return &fav, nil
}

// CreateFavorite 创建收藏
func (r *ProfileRepositoryImpl) CreateFavorite(ctx context.Context, fav *entity.UserFav) error {
	return r.db.WithContext(ctx).Save(fav).Error
}

// DeleteFavorite 删除收藏
func (r *ProfileRepositoryImpl) DeleteFavorite(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.UserFav{}, id).Error
}

// IsFavorite 检查商品是否已收藏
func (r *ProfileRepositoryImpl) IsFavorite(ctx context.Context, userID, goodsID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.UserFav{}).
		Where("user = ? AND goods = ?", userID, goodsID).
		Count(&count).Error
	
	return count > 0, err
}

// ======= 地址相关 =======

// GetAddressesByUserID 获取用户地址列表
func (r *ProfileRepositoryImpl) GetAddressesByUserID(ctx context.Context, userID int64) ([]*entity.Address, error) {
	var addresses []*entity.Address
	
	err := r.db.WithContext(ctx).
		Where("user = ?", userID).
		Order("is_default DESC, usage_count DESC"). // 默认地址排在最前，然后是使用频率高的地址
		Find(&addresses).Error
	
	return addresses, err
}

// GetAddressByID 获取地址详情
func (r *ProfileRepositoryImpl) GetAddressByID(ctx context.Context, id int64) (*entity.Address, error) {
	var address entity.Address
	result := r.db.WithContext(ctx).First(&address, id)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return &address, nil
}

// CreateAddress 创建地址
func (r *ProfileRepositoryImpl) CreateAddress(ctx context.Context, address *entity.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

// UpdateAddress 更新地址
func (r *ProfileRepositoryImpl) UpdateAddress(ctx context.Context, address *entity.Address) error {
	return r.db.WithContext(ctx).Save(address).Error
}

// DeleteAddress 删除地址
func (r *ProfileRepositoryImpl) DeleteAddress(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.Address{}, id).Error
}

// SetDefaultAddress 设置默认地址
func (r *ProfileRepositoryImpl) SetDefaultAddress(ctx context.Context, userID, addressID int64) error {
	// 使用事务确保数据一致性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 先将该用户的所有地址设为非默认
		if err := tx.Model(&entity.Address{}).
			Where("user = ?", userID).
			Update("is_default", false).Error; err != nil {
			return err
		}
		
		// 2. 再将指定地址设为默认
		if err := tx.Model(&entity.Address{}).
			Where("id = ? AND user = ?", addressID, userID).
			Update("is_default", true).Error; err != nil {
			return err
		}
		
		return nil
	})
}

// GetDefaultAddress 获取默认地址
func (r *ProfileRepositoryImpl) GetDefaultAddress(ctx context.Context, userID int64) (*entity.Address, error) {
	var address entity.Address
	result := r.db.WithContext(ctx).
		Where("user = ? AND is_default = ?", userID, true).
		First(&address)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return &address, nil
}

// ======= 用户反馈相关 =======

// GetFeedbacksByUserID 获取用户反馈列表
func (r *ProfileRepositoryImpl) GetFeedbacksByUserID(ctx context.Context, userID int64, offset, limit int) ([]*entity.UserFeedback, int64, error) {
	var feedbacks []*entity.UserFeedback
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&entity.UserFeedback{}).
		Where("user = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页获取记录
	if err := r.db.WithContext(ctx).
		Where("user = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}
	
	return feedbacks, total, nil
}

// GetFeedbackByID 获取反馈详情
func (r *ProfileRepositoryImpl) GetFeedbackByID(ctx context.Context, id int64) (*entity.UserFeedback, error) {
	var feedback entity.UserFeedback
	result := r.db.WithContext(ctx).First(&feedback, id)
	
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return &feedback, nil
}

// CreateFeedback 创建反馈
func (r *ProfileRepositoryImpl) CreateFeedback(ctx context.Context, feedback *entity.UserFeedback) error {
	return r.db.WithContext(ctx).Create(feedback).Error
}

// UpdateFeedback 更新反馈
func (r *ProfileRepositoryImpl) UpdateFeedback(ctx context.Context, feedback *entity.UserFeedback) error {
	return r.db.WithContext(ctx).Save(feedback).Error
}

// DeleteFeedback 删除反馈
func (r *ProfileRepositoryImpl) DeleteFeedback(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.UserFeedback{}, id).Error
}

// ======= 浏览历史相关 =======

// GetBrowsingHistories 获取浏览历史
func (r *ProfileRepositoryImpl) GetBrowsingHistories(ctx context.Context, userID int64, offset, limit int) ([]*entity.BrowsingHistory, int64, error) {
	var histories []*entity.BrowsingHistory
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&entity.BrowsingHistory{}).
		Where("user = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页获取记录
	if err := r.db.WithContext(ctx).
		Where("user = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}
	
	return histories, total, nil
}

// AddBrowsingHistory 添加浏览历史
func (r *ProfileRepositoryImpl) AddBrowsingHistory(ctx context.Context, history *entity.BrowsingHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// DeleteBrowsingHistory 删除浏览历史
func (r *ProfileRepositoryImpl) DeleteBrowsingHistory(ctx context.Context, userID int64, ids []int64) error {
	return r.db.WithContext(ctx).
		Where("user = ? AND id IN ?", userID, ids).
		Delete(&entity.BrowsingHistory{}).Error
}

// ClearBrowsingHistory 清空浏览历史
func (r *ProfileRepositoryImpl) ClearBrowsingHistory(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Where("user = ?", userID).
		Delete(&entity.BrowsingHistory{}).Error
}
