package service

import (
	"context"

	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/domain/valueobject"
)

// AuthService 认证服务接口
type AuthService interface {
	// Login 用户登录
	Login(ctx context.Context, username, password string) (*entity.User, *valueobject.Credential, error)

	// Logout 用户登出
	Logout(ctx context.Context, token string) error

	// RefreshToken 刷新访问令牌
	RefreshToken(ctx context.Context, refreshToken string) (*valueobject.Credential, error)

	// GenerateToken 生成用户令牌
	GenerateToken(ctx context.Context, user *entity.User) (*valueobject.Credential, error)

	// ValidateToken 验证访问令牌
	ValidateToken(ctx context.Context, token string) (*valueobject.TokenClaims, error)

	// RegisterUser 注册用户
	RegisterUser(ctx context.Context, user *entity.User, password string) (*entity.User, error)

	// VerifyPassword 验证用户密码
	VerifyPassword(ctx context.Context, userID int64, password string) (bool, error)

	// HashPassword 哈希用户密码
	HashPassword(password string) (string, error)

	// GenerateVerificationCode 生成验证码
	GenerateVerificationCode(ctx context.Context, phone string, codeType string) (string, error)

	// VerifyVerificationCode 验证验证码
	VerifyVerificationCode(ctx context.Context, phone, code, codeType string) (bool, error)
}
