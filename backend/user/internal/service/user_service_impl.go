package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/repository"
	"shop/backend/user/internal/repository/cache"
)

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo    repository.UserRepository
	userCache   *cache.RedisUserCache
	authService AuthService
}

// NewUserService 创建用户服务
func NewUserService(
	userRepo repository.UserRepository,
	userCache *cache.RedisUserCache,
	authService AuthService,
) UserService {
	return &UserServiceImpl{
		userRepo:    userRepo,
		userCache:   userCache,
		authService: authService,
	}
}

// CreateUser 创建用户
func (s *UserServiceImpl) CreateUser(ctx context.Context, user *entity.User, password string) (*entity.User, error) {
	// 验证用户名是否已存在
	existingUser, err := s.userRepo.GetUserByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 如果提供了手机号，验证手机号是否已存在
	if user.Phone != "" {
		existingUser, err = s.userRepo.GetUserByPhone(ctx, user.Phone)
		if err == nil && existingUser != nil {
			return nil, errors.New("手机号已注册")
		}
	}

	// 如果提供了邮箱，验证邮箱是否已存在
	if user.Email != "" {
		existingUser, err = s.userRepo.GetUserByEmail(ctx, user.Email)
		if err == nil && existingUser != nil {
			return nil, errors.New("邮箱已注册")
		}
	}

	// 设置用户默认信息
	if user.Nickname == "" {
		user.Nickname = user.Username
	}
	user.Status = 1 // 正常状态
	user.Role = 1   // 普通用户
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// 创建用户
	createdUser, err := s.userRepo.CreateUser(ctx, user, password)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return createdUser, nil
}

// GetUserByID 通过ID获取用户
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	// 先尝试从缓存获取
	cachedUser, err := s.userCache.GetUser(ctx, id)
	if err == nil {
		return cachedUser, nil
	}

	// 缓存未命中，从数据库获取
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 存入缓存
	if err := s.userCache.SetUser(ctx, user, 0); err != nil {
		// 缓存错误只记录不返回
		fmt.Printf("缓存用户信息失败: %v", err)
	}

	return user, nil
}

// GetUserByUsername 通过用户名获取用户
func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 存入缓存
	if err := s.userCache.SetUser(ctx, user, 0); err != nil {
		// 缓存错误只记录不返回
		fmt.Printf("缓存用户信息失败: %v", err)
	}

	return user, nil
}

// GetUserByPhone 通过手机号获取用户
func (s *UserServiceImpl) GetUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	user, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 存入缓存
	if err := s.userCache.SetUser(ctx, user, 0); err != nil {
		// 缓存错误只记录不返回
		fmt.Printf("缓存用户信息失败: %v", err)
	}

	return user, nil
}

// GetUserByEmail 通过邮箱获取用户
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 存入缓存
	if err := s.userCache.SetUser(ctx, user, 0); err != nil {
		// 缓存错误只记录不返回
		fmt.Printf("缓存用户信息失败: %v", err)
	}

	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// 验证用户是否存在
	existingUser, err := s.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 只允许更新部分字段
	existingUser.Nickname = user.Nickname
	existingUser.Avatar = user.Avatar
	existingUser.Gender = user.Gender
	existingUser.Birthday = user.Birthday
	existingUser.Email = user.Email
	existingUser.Phone = user.Phone
	existingUser.UpdatedAt = time.Now()

	// 更新用户
	updatedUser, err := s.userRepo.UpdateUser(ctx, existingUser)
	if err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	// 更新缓存
	if err := s.userCache.DeleteUser(ctx, user.ID); err != nil {
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return updatedUser, nil
}

// ChangePassword 修改用户密码
func (s *UserServiceImpl) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	// 验证旧密码
	valid, err := s.authService.VerifyPassword(ctx, userID, oldPassword)
	if err != nil {
		return fmt.Errorf("验证密码失败: %w", err)
	}
	if !valid {
		return errors.New("旧密码不正确")
	}

	// 哈希新密码
	hashedPassword, err := s.authService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码哈希失败: %w", err)
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	// 删除用户
	if err := s.userRepo.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	// 删除缓存
	if err := s.userCache.DeleteUser(ctx, id); err != nil {
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// LockUser 锁定用户
func (s *UserServiceImpl) LockUser(ctx context.Context, id int64) error {
	// 更新用户状态为锁定
	if err := s.userRepo.UpdateStatus(ctx, id, 3); err != nil {
		return fmt.Errorf("锁定用户失败: %w", err)
	}

	// 删除缓存
	if err := s.userCache.DeleteUser(ctx, id); err != nil {
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// UnlockUser 解锁用户
func (s *UserServiceImpl) UnlockUser(ctx context.Context, id int64) error {
	// 更新用户状态为正常
	if err := s.userRepo.UpdateStatus(ctx, id, 1); err != nil {
		return fmt.Errorf("解锁用户失败: %w", err)
	}

	// 删除缓存
	if err := s.userCache.DeleteUser(ctx, id); err != nil {
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// ListUsers 获取用户列表
func (s *UserServiceImpl) ListUsers(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	// 计算分页
	offset := (page - 1) * pageSize
	limit := pageSize

	// 获取用户列表
	users, err := s.userRepo.ListUsers(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 获取用户总数
	total, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户总数失败: %w", err)
	}

	return users, total, nil
}

// UpdateUserStatus 更新用户状态
func (s *UserServiceImpl) UpdateUserStatus(ctx context.Context, id int64, status int) error {
	// 验证状态值
	if status < 1 || status > 3 {
		return errors.New("无效的状态值")
	}

	// 更新用户状态
	if err := s.userRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("更新用户状态失败: %w", err)
	}

	// 删除缓存
	if err := s.userCache.DeleteUser(ctx, id); err != nil {
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// BindWechat 绑定微信
func (s *UserServiceImpl) BindWechat(ctx context.Context, userID int64, openID, unionID string) error {
	// 验证用户是否存在
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 更新微信信息
	user.WechatOpenID = openID
	user.WechatUnionID = unionID
	user.UpdatedAt = time.Now()

	// 更新用户
	_, err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("绑定微信失败: %w", err)
	}

	// 删除缓存
	if err := s.userCache.DeleteUser(ctx, userID); err != nil {
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// UnbindWechat 解绑微信
func (s *UserServiceImpl) UnbindWechat(ctx context.Context, userID int64) error {
	// 验证用户是否存在
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 清除微信信息
	user.WechatOpenID = ""
	user.WechatUnionID = ""
	user.UpdatedAt = time.Now()

	// 更新用户
	_, err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("解绑微信失败: %w", err)
	}

	// 删除缓存
	if err := s.userCache.DeleteUser(ctx, userID); err != nil {
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// GetUserPermissions 获取用户权限
func (s *UserServiceImpl) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	// 获取用户
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 根据用户角色返回不同权限
	var permissions []string
	if user.Role == 1 { // 普通用户
		permissions = []string{
			"user:info",
			"user:update",
			"user:password",
		}
	} else if user.Role == 2 { // 管理员
		permissions = []string{
			"user:info",
			"user:update",
			"user:password",
			"user:list",
			"user:create",
			"user:delete",
			"user:status",
			"admin:access",
		}
	}

	return permissions, nil
}
