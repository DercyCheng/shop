package service

import (
	"context"
	"errors"
	"time"
	
	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/domain/valueobject"
	"shop/backend/user/pkg/jwt"
	
	"golang.org/x/crypto/bcrypt"
)

// ErrUserNotFound 用户不存在错误
var ErrUserNotFound = errors.New("user not found")

// ErrInvalidCredentials 无效的凭证错误
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrUserDisabled 用户被禁用错误
var ErrUserDisabled = errors.New("user is disabled")

// ErrUserLocked 用户被锁定错误
var ErrUserLocked = errors.New("user is locked")

// ErrMobileExists 手机号已存在错误
var ErrMobileExists = errors.New("mobile already exists")

// ErrInvalidCode 验证码错误
var ErrInvalidCode = errors.New("invalid verification code")

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	repo UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(repo UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

// GetUserByID 根据ID获取用户
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, ErrUserNotFound
	}
	
	return user, nil
}

// GetUserByMobile 根据手机号获取用户
func (s *UserServiceImpl) GetUserByMobile(ctx context.Context, mobile string) (*entity.User, error) {
	user, err := s.repo.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, ErrUserNotFound
	}
	
	return user, nil
}

// RegisterUser 注册新用户
func (s *UserServiceImpl) RegisterUser(ctx context.Context, mobile, password, nickname string) (*entity.User, error) {
	// 检查手机号是否已存在
	exists, err := s.repo.MobileExists(ctx, mobile)
	if err != nil {
		return nil, err
	}
	
	if exists {
		return nil, ErrMobileExists
	}
	
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	// 创建用户实体
	now := time.Now()
	user := &entity.User{
		Mobile:        mobile,
		Password:      string(hashedPassword),
		Nickname:      nickname,
		Status:        entity.UserStatusNormal,
		Role:          entity.UserRoleNormal,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	
	// 保存用户
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	
	return user, nil
}

// ListUsers 获取用户列表
func (s *UserServiceImpl) ListUsers(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	if page < 1 {
		page = 1
	}
	
	if pageSize < 1 {
		pageSize = 10
	}
	
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, offset, pageSize)
}

// UpdateUser 更新用户信息
func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *entity.User) error {
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	
	if existingUser == nil {
		return ErrUserNotFound
	}
	
	// 保留不允许修改的字段
	user.Password = existingUser.Password
	user.Status = existingUser.Status
	user.Role = existingUser.Role
	user.CreatedAt = existingUser.CreatedAt
	
	// 更新时间
	user.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, user)
}

// UpdatePassword 更新用户密码
func (s *UserServiceImpl) UpdatePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	
	if user == nil {
		return ErrUserNotFound
	}
	
	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}
	
	// 生成新密码的哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	// 更新密码
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, user)
}

// ResetPassword 重置用户密码（忘记密码）
func (s *UserServiceImpl) ResetPassword(ctx context.Context, mobile, newPassword, verificationCode string) error {
	// TODO: 验证验证码
	// 这里应该调用短信验证码服务验证码是否正确
	// 简化实现，假设验证通过
	
	user, err := s.repo.GetByMobile(ctx, mobile)
	if err != nil {
		return err
	}
	
	if user == nil {
		return ErrUserNotFound
	}
	
	// 生成新密码的哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	// 更新密码
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, user)
}

// BindWechat 绑定微信账号
func (s *UserServiceImpl) BindWechat(ctx context.Context, userID int64, openID, unionID string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	
	if user == nil {
		return ErrUserNotFound
	}
	
	// 检查OpenID是否已被其他用户绑定
	existingUser, err := s.repo.GetByWechatOpenID(ctx, openID)
	if err != nil {
		return err
	}
	
	if existingUser != nil && existingUser.ID != userID {
		return errors.New("wechat account already bound to another user")
	}
	
	// 更新微信绑定信息
	user.WechatOpenID = openID
	user.WechatUnionID = unionID
	user.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, user)
}

// UnbindWechat 解绑微信账号
func (s *UserServiceImpl) UnbindWechat(ctx context.Context, userID int64) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	
	if user == nil {
		return ErrUserNotFound
	}
	
	// 清除微信绑定信息
	user.WechatOpenID = ""
	user.WechatUnionID = ""
	user.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, user)
}
