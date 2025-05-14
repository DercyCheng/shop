package service

import (
	"context"
	
	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/domain/valueobject"
)

// UserService 用户服务接口
type UserService interface {
	// 根据ID获取用户
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	
	// 根据手机号获取用户
	GetUserByMobile(ctx context.Context, mobile string) (*entity.User, error)
	
	// 创建用户（注册）
	RegisterUser(ctx context.Context, mobile, password, nickname string) (*entity.User, error)
	
	// 用户列表（分页查询）
	ListUsers(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error)
	
	// 更新用户信息
	UpdateUser(ctx context.Context, user *entity.User) error
	
	// 更新用户密码
	UpdatePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	
	// 重置用户密码（忘记密码）
	ResetPassword(ctx context.Context, mobile, newPassword, verificationCode string) error
	
	// 绑定微信账号
	BindWechat(ctx context.Context, userID int64, openID, unionID string) error
	
	// 解绑微信账号
	UnbindWechat(ctx context.Context, userID int64) error
}

// AuthService 认证服务接口
type AuthService interface {
	// 用户登录
	Login(ctx context.Context, credential *valueobject.Credential) (*valueobject.LoginResponse, error)
	
	// 验证Token
	ValidateToken(ctx context.Context, token string) (int64, string, error)
	
	// 刷新Token
	RefreshToken(ctx context.Context, refreshToken string) (*valueobject.LoginResponse, error)
	
	// 退出登录
	Logout(ctx context.Context, token string) error
}
