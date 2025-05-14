package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	
	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/repository/cache"
	"shop/backend/user/internal/service"
	
	"gorm.io/gorm"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db    *gorm.DB
	cache cache.UserCache
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB, cache cache.UserCache) service.UserRepository {
	return &UserRepositoryImpl{
		db:    db,
		cache: cache,
	}
}

// GetByID 根据ID获取用户信息
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	// 1. 尝试从缓存获取
	user, err := r.cache.GetUser(ctx, id)
	if err == nil && user != nil {
		return user, nil
	}
	
	// 2. 缓存未命中，从数据库查询
	user = &entity.User{}
	result := r.db.WithContext(ctx).First(user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	// 3. 将用户信息存入缓存
	if err := r.cache.SetUser(ctx, user); err != nil {
		// 缓存设置失败仅记录日志，不影响主流程
		// log.Printf("Failed to set user cache: %v", err)
	}
	
	return user, nil
}

// GetByMobile 根据手机号获取用户
func (r *UserRepositoryImpl) GetByMobile(ctx context.Context, mobile string) (*entity.User, error) {
	// 1. 尝试从缓存获取
	user, err := r.cache.GetUserByMobile(ctx, mobile)
	if err == nil && user != nil {
		return user, nil
	}
	
	// 2. 缓存未命中，从数据库查询
	user = &entity.User{}
	result := r.db.WithContext(ctx).Where("mobile = ?", mobile).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	// 3. 将用户信息存入缓存
	if err := r.cache.SetUser(ctx, user); err != nil {
		// 缓存设置失败仅记录日志，不影响主流程
	}
	
	return user, nil
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	// 设置创建和更新时间
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	
	// 创建用户记录
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return result.Error
	}
	
	// 将新用户信息存入缓存
	return r.cache.SetUser(ctx, user)
}

// Update 更新用户信息
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	// 设置更新时间
	user.UpdatedAt = time.Now()
	
	// 更新用户记录
	result := r.db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}
	
	// 更新缓存
	return r.cache.SetUser(ctx, user)
}

// Delete 删除用户（软删除）
func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	// 软删除用户记录
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"updated_at": now,
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	// 从缓存中删除用户
	return r.cache.DeleteUser(ctx, id)
}

// GetByWechatOpenID 根据微信OpenID获取用户
func (r *UserRepositoryImpl) GetByWechatOpenID(ctx context.Context, openID string) (*entity.User, error) {
	user := &entity.User{}
	result := r.db.WithContext(ctx).Where("wechat_open_id = ?", openID).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	
	return user, nil
}

// MobileExists 检查手机号是否已存在
func (r *UserRepositoryImpl) MobileExists(ctx context.Context, mobile string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("mobile = ?", mobile).
		Count(&count)
	
	if result.Error != nil {
		return false, result.Error
	}
	
	return count > 0, nil
}

// List 获取用户列表（分页）
func (r *UserRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64
	
	// 获取总记录数
	result := r.db.WithContext(ctx).Model(&entity.User{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	
	// 分页查询
	result = r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&users)
	
	if result.Error != nil {
		return nil, 0, result.Error
	}
	
	return users, total, nil
}
