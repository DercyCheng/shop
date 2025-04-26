package repository

import (
	"context"
	"time"

	"shop/backend/user/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 创建用户
	CreateUser(ctx context.Context, user *entity.User, password string) (*entity.User, error)

	// 获取用户
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByWechatOpenID(ctx context.Context, openID string) (*entity.User, error)

	// 更新用户
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdatePassword(ctx context.Context, id int64, password string) error
	UpdateStatus(ctx context.Context, id int64, status int) error
	UpdateLoginFailCount(ctx context.Context, id int64, count int) error
	UpdateLastLogin(ctx context.Context, id int64) error

	// 删除用户
	DeleteUser(ctx context.Context, id int64) error

	// 用户认证
	VerifyPassword(ctx context.Context, userID int64, password string) (bool, error)

	// 令牌管理
	StoreToken(ctx context.Context, token string, userID int64, expiry time.Duration) error
	VerifyToken(ctx context.Context, token string) (int64, error)
	InvalidateToken(ctx context.Context, token string) error

	// 刷新令牌管理
	StoreRefreshToken(ctx context.Context, refreshToken string, userID int64, expiry time.Duration) error
	VerifyRefreshToken(ctx context.Context, refreshToken string) (int64, error)
	InvalidateRefreshToken(ctx context.Context, refreshToken string) error

	// 批量操作
	ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error)
	CountUsers(ctx context.Context) (int64, error)
}
