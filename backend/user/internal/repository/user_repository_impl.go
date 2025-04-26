package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/repository/cache"
	"shop/backend/user/internal/repository/dao"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	userDAO   *dao.UserDAO
	userCache *cache.RedisUserCache
}

// NewUserRepository 创建用户仓储
func NewUserRepository(userDAO *dao.UserDAO, userCache *cache.RedisUserCache) UserRepository {
	return &UserRepositoryImpl{
		userDAO:   userDAO,
		userCache: userCache,
	}
}

// CreateUser 创建用户
func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *entity.User, password string) (*entity.User, error) {
	// 转换为DAO模型
	userModel := &dao.User{
		Username:  user.Username,
		Password:  password,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if !user.Birthday.IsZero() {
		userModel.Birthday = &user.Birthday
	}

	// 创建用户
	err := r.userDAO.Create(ctx, userModel)
	if err != nil {
		return nil, err
	}

	// 转换为领域实体并返回
	return r.daoToEntity(userModel), nil
}

// GetUserByID 通过ID获取用户
func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	// 先尝试从缓存获取
	cachedUser, err := r.userCache.GetUser(ctx, id)
	if err == nil {
		return cachedUser, nil
	}

	// 缓存未命中，从数据库获取
	userModel, err := r.userDAO.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在: %w", err)
		}
		return nil, err
	}

	// 转换为领域实体
	user := r.daoToEntity(userModel)

	// 缓存用户信息
	if err := r.userCache.SetUser(ctx, user, 0); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("缓存用户信息失败: %v", err)
	}

	return user, nil
}

// GetUserByUsername 通过用户名获取用户
func (r *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	userModel, err := r.userDAO.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在: %w", err)
		}
		return nil, err
	}

	// 转换为领域实体
	user := r.daoToEntity(userModel)

	// 缓存用户信息
	if err := r.userCache.SetUser(ctx, user, 0); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("缓存用户信息失败: %v", err)
	}

	return user, nil
}

// GetUserByPhone 通过手机号获取用户
func (r *UserRepositoryImpl) GetUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	userModel, err := r.userDAO.FindByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在: %w", err)
		}
		return nil, err
	}

	// 转换为领域实体
	user := r.daoToEntity(userModel)

	// 缓存用户信息
	if err := r.userCache.SetUser(ctx, user, 0); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("缓存用户信息失败: %v", err)
	}

	return user, nil
}

// GetUserByEmail 通过邮箱获取用户
func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	userModel, err := r.userDAO.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在: %w", err)
		}
		return nil, err
	}

	// 转换为领域实体
	user := r.daoToEntity(userModel)

	// 缓存用户信息
	if err := r.userCache.SetUser(ctx, user, 0); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("缓存用户信息失败: %v", err)
	}

	return user, nil
}

// GetUserByWechatOpenID 通过微信OpenID获取用户
func (r *UserRepositoryImpl) GetUserByWechatOpenID(ctx context.Context, openID string) (*entity.User, error) {
	userModel, err := r.userDAO.FindByWechatOpenID(ctx, openID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在: %w", err)
		}
		return nil, err
	}

	// 转换为领域实体
	user := r.daoToEntity(userModel)

	// 缓存用户信息
	if err := r.userCache.SetUser(ctx, user, 0); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("缓存用户信息失败: %v", err)
	}

	return user, nil
}

// UpdateUser 更新用户
func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// 先获取当前用户信息
	existingUser, err := r.userDAO.FindByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 更新用户信息
	existingUser.Nickname = user.Nickname
	existingUser.Avatar = user.Avatar
	existingUser.Email = user.Email
	existingUser.Phone = user.Phone
	existingUser.Gender = user.Gender
	existingUser.WechatOpenID = user.WechatOpenID
	existingUser.WechatUnionID = user.WechatUnionID
	existingUser.UpdatedAt = time.Now()

	if !user.Birthday.IsZero() {
		existingUser.Birthday = &user.Birthday
	}

	// 保存更新
	err = r.userDAO.Update(ctx, existingUser)
	if err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	// 删除缓存
	if err := r.userCache.DeleteUser(ctx, user.ID); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	// 转换为领域实体并返回
	return r.daoToEntity(existingUser), nil
}

// UpdatePassword 更新用户密码
func (r *UserRepositoryImpl) UpdatePassword(ctx context.Context, id int64, password string) error {
	err := r.userDAO.UpdatePassword(ctx, id, password)
	if err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	// 删除缓存
	if err := r.userCache.DeleteUser(ctx, id); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// UpdateStatus 更新用户状态
func (r *UserRepositoryImpl) UpdateStatus(ctx context.Context, id int64, status int) error {
	err := r.userDAO.UpdateStatus(ctx, id, status)
	if err != nil {
		return fmt.Errorf("更新状态失败: %w", err)
	}

	// 删除缓存
	if err := r.userCache.DeleteUser(ctx, id); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// UpdateLoginFailCount 更新登录失败次数
func (r *UserRepositoryImpl) UpdateLoginFailCount(ctx context.Context, id int64, count int) error {
	err := r.userDAO.UpdateLoginFailCount(ctx, id, count)
	if err != nil {
		return fmt.Errorf("更新登录失败次数失败: %w", err)
	}

	// 删除缓存
	if err := r.userCache.DeleteUser(ctx, id); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// UpdateLastLogin 更新最后登录时间
func (r *UserRepositoryImpl) UpdateLastLogin(ctx context.Context, id int64) error {
	err := r.userDAO.UpdateLastLogin(ctx, id)
	if err != nil {
		return fmt.Errorf("更新最后登录时间失败: %w", err)
	}

	// 删除缓存
	if err := r.userCache.DeleteUser(ctx, id); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// DeleteUser 删除用户
func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, id int64) error {
	err := r.userDAO.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	// 删除缓存
	if err := r.userCache.DeleteUser(ctx, id); err != nil {
		// 仅记录错误，不影响返回
		fmt.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// VerifyPassword 验证用户密码
func (r *UserRepositoryImpl) VerifyPassword(ctx context.Context, userID int64, password string) (bool, error) {
	// 获取用户
	user, err := r.userDAO.FindByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("用户不存在: %w", err)
	}

	// 解析密码和盐值
	parts := strings.Split(user.Password, ":")
	if len(parts) != 2 {
		return false, errors.New("密码格式无效")
	}

	storedHash, salt := parts[0], parts[1]

	// 使用相同盐值哈希输入的密码
	hash := r.hashPassword(password, salt)

	// 比较哈希值
	return hash == storedHash, nil
}

// StoreToken 存储令牌
func (r *UserRepositoryImpl) StoreToken(ctx context.Context, token string, userID int64, expiry time.Duration) error {
	return r.userCache.SaveToken(ctx, token, userID, expiry)
}

// VerifyToken 验证令牌
func (r *UserRepositoryImpl) VerifyToken(ctx context.Context, token string) (int64, error) {
	return r.userCache.GetUserIDByToken(ctx, token)
}

// InvalidateToken 使令牌失效
func (r *UserRepositoryImpl) InvalidateToken(ctx context.Context, token string) error {
	return r.userCache.InvalidateToken(ctx, token)
}

// StoreRefreshToken 存储刷新令牌
func (r *UserRepositoryImpl) StoreRefreshToken(ctx context.Context, refreshToken string, userID int64, expiry time.Duration) error {
	return r.userCache.SaveRefreshToken(ctx, refreshToken, userID, expiry)
}

// VerifyRefreshToken 验证刷新令牌
func (r *UserRepositoryImpl) VerifyRefreshToken(ctx context.Context, refreshToken string) (int64, error) {
	return r.userCache.GetUserIDByRefreshToken(ctx, refreshToken)
}

// InvalidateRefreshToken 使刷新令牌失效
func (r *UserRepositoryImpl) InvalidateRefreshToken(ctx context.Context, refreshToken string) error {
	return r.userCache.InvalidateRefreshToken(ctx, refreshToken)
}

// ListUsers 获取用户列表
func (r *UserRepositoryImpl) ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	usersModel, err := r.userDAO.ListUsers(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 转换为领域实体列表
	var users []*entity.User
	for _, userModel := range usersModel {
		users = append(users, r.daoToEntity(userModel))
	}

	return users, nil
}

// CountUsers 获取用户总数
func (r *UserRepositoryImpl) CountUsers(ctx context.Context) (int64, error) {
	return r.userDAO.CountUsers(ctx)
}

// daoToEntity 将DAO模型转换为领域实体
func (r *UserRepositoryImpl) daoToEntity(userModel *dao.User) *entity.User {
	user := &entity.User{
		ID:             userModel.ID,
		Username:       userModel.Username,
		Nickname:       userModel.Nickname,
		Email:          userModel.Email,
		Phone:          userModel.Phone,
		Password:       userModel.Password,
		Avatar:         userModel.Avatar,
		Gender:         userModel.Gender,
		Status:         userModel.Status,
		Role:           userModel.Role,
		LoginFailCount: userModel.LoginFailCount,
		WechatOpenID:   userModel.WechatOpenID,
		WechatUnionID:  userModel.WechatUnionID,
		SessionID:      userModel.SessionID,
		CreatedAt:      userModel.CreatedAt,
		UpdatedAt:      userModel.UpdatedAt,
	}

	if userModel.Birthday != nil {
		user.Birthday = *userModel.Birthday
	}

	if userModel.LastLoginAt != nil {
		user.LastLoginAt = *userModel.LastLoginAt
	}

	if userModel.DeletedAt != nil {
		user.DeletedAt = *userModel.DeletedAt
	}

	return user
}

// hashPassword 哈希密码（仅用于验证）
func (r *UserRepositoryImpl) hashPassword(password, salt string) string {
	// 这里应该实现与 authService.HashPassword 相同的哈希算法
	// 在实际实现中，可能会将此逻辑抽取到共享的工具包
	// 这里只是一个简单的示例
	hash := fmt.Sprintf("%x", password+salt)
	return hash
}
