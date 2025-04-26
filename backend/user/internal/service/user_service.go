package service

import (
	"context"

	"shop/backend/user/internal/domain/entity"
)

// UserService 用户服务接口
type UserService interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, user *entity.User, password string) (*entity.User, error)

	// GetUserByID 通过ID获取用户
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)

	// GetUserByUsername 通过用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)

	// GetUserByPhone 通过手机号获取用户
	GetUserByPhone(ctx context.Context, phone string) (*entity.User, error)

	// GetUserByEmail 通过邮箱获取用户
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)

	// ChangePassword 修改用户密码
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error

	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, id int64) error

	// LockUser 锁定用户
	LockUser(ctx context.Context, id int64) error

	// UnlockUser 解锁用户
	UnlockUser(ctx context.Context, id int64) error

	// ListUsers 获取用户列表
	ListUsers(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error)

	// UpdateUserStatus 更新用户状态
	UpdateUserStatus(ctx context.Context, id int64, status int) error

	// BindWechat 绑定微信
	BindWechat(ctx context.Context, userID int64, openID, unionID string) error

	// UnbindWechat 解绑微信
	UnbindWechat(ctx context.Context, userID int64) error

	// GetUserPermissions 获取用户权限
	GetUserPermissions(ctx context.Context, userID int64) ([]string, error)
}
