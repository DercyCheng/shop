package service

import (
	"context"
	
	"shop/backend/user/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 根据ID获取用户
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	
	// 根据手机号获取用户
	GetByMobile(ctx context.Context, mobile string) (*entity.User, error)
	
	// 创建用户
	Create(ctx context.Context, user *entity.User) error
	
	// 更新用户信息
	Update(ctx context.Context, user *entity.User) error
	
	// 删除用户（软删除）
	Delete(ctx context.Context, id int64) error
	
	// 根据微信OpenID获取用户
	GetByWechatOpenID(ctx context.Context, openID string) (*entity.User, error)
	
	// 检查手机号是否已存在
	MobileExists(ctx context.Context, mobile string) (bool, error)
	
	// 获取用户列表（分页）
	List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)
}
