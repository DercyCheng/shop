package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/repository"
)

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo    repository.UserRepository
	authService AuthService
	logger      *zap.Logger
	mongoClient *mongo.Client // 添加MongoDB客户端
}

// NewUserService 创建用户服务
func NewUserService(
	userRepo repository.UserRepository,
	authService AuthService,
	logger *zap.Logger,
	mongoClient *mongo.Client, // 添加MongoDB客户端参数
) UserService {
	return &UserServiceImpl{
		userRepo:    userRepo,
		authService: authService,
		logger:      logger,
		mongoClient: mongoClient, // 初始化MongoDB客户端
	}
}

// CreateUser 创建用户
func (s *UserServiceImpl) CreateUser(ctx context.Context, user *entity.User, password string) (*entity.User, error) {
	return s.authService.RegisterUser(ctx, user, password)
}

// GetUserByID 通过ID获取用户
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get user by ID", zap.Int64("user_id", id), zap.Error(err))
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return user, nil
}

// GetUserByUsername 通过用户名获取用户
func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.Error("Failed to get user by username", zap.String("username", username), zap.Error(err))
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return user, nil
}

// GetUserByPhone 通过手机号获取用户
func (s *UserServiceImpl) GetUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	user, err := s.userRepo.GetUserByPhone(ctx, phone)
	if err != nil {
		s.logger.Error("Failed to get user by phone", zap.String("phone", phone), zap.Error(err))
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return user, nil
}

// GetUserByEmail 通过邮箱获取用户
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to get user by email", zap.String("email", email), zap.Error(err))
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// 准备操作日志
	log := &UserOperationLog{
		Operation: "update_user",
		UserID:    user.ID,
		Status:    "processing",
		Timestamp: time.Now(),
		UserAgent: "system",
		Details: map[string]interface{}{
			"update_fields": map[string]interface{}{
				"nickname": user.Nickname,
				"email":    user.Email,
				"phone":    user.Phone,
				"gender":   user.Gender,
				"avatar":   user.Avatar,
			},
		},
	}

	// 先验证用户是否存在
	existingUser, err := s.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		s.logger.Error("User not found for update", zap.Int64("user_id", user.ID), zap.Error(err))

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("用户不存在: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)

		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 保留一些不应由客户端更新的字段
	user.Password = existingUser.Password
	user.Status = existingUser.Status
	user.Role = existingUser.Role
	user.LoginFailCount = existingUser.LoginFailCount
	user.LastLoginAt = existingUser.LastLoginAt
	user.CreatedAt = existingUser.CreatedAt

	// 更新用户
	updatedUser, err := s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update user", zap.Int64("user_id", user.ID), zap.Error(err))

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("更新用户失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)

		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	s.logger.Info("User updated successfully", zap.Int64("user_id", updatedUser.ID))

	// 记录成功日志
	log.Status = "success"
	s.authService.LogUserOperation(ctx, log, nil)

	return updatedUser, nil
}

// ChangePassword 修改用户密码
func (s *UserServiceImpl) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	// 准备操作日志
	log := &UserOperationLog{
		Operation: "change_password",
		UserID:    userID,
		Status:    "processing",
		Timestamp: time.Now(),
		UserAgent: "system",
	}

	// 验证旧密码
	valid, err := s.authService.VerifyPassword(ctx, userID, oldPassword)
	if err != nil {
		s.logger.Error("Failed to verify old password", zap.Int64("user_id", userID), zap.Error(err))

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("验证旧密码失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)

		return fmt.Errorf("验证旧密码失败: %w", err)
	}

	if !valid {
		s.logger.Info("Invalid old password", zap.Int64("user_id", userID))

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = "旧密码不正确"
		s.authService.LogUserOperation(ctx, log, nil)

		return errors.New("旧密码不正确")
	}

	// 哈希新密码
	hashedPassword, err := s.authService.HashPassword(newPassword)
	if err != nil {
		s.logger.Error("Failed to hash new password", zap.Int64("user_id", userID), zap.Error(err))

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("密码加密失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)

		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	err = s.userRepo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		s.logger.Error("Failed to update password", zap.Int64("user_id", userID), zap.Error(err))

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("更新密码失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)

		return fmt.Errorf("更新密码失败: %w", err)
	}

	s.logger.Info("Password changed successfully", zap.Int64("user_id", userID))

	// 记录成功日志
	log.Status = "success"
	s.authService.LogUserOperation(ctx, log, nil)

	return nil
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	// 准备操作日志
	log := &UserOperationLog{
		Operation: "delete_user",
		UserID:    id,
		Status:    "processing",
		Timestamp: time.Now(),
		UserAgent: "system",
	}

	// 先获取用户信息，用于记录日志
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("用户不存在: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)

		s.logger.Error("User not found for deletion", zap.Int64("user_id", id), zap.Error(err))
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 记录用户信息用于审计
	log.Details = map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"phone":    user.Phone,
	}

	err = s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete user", zap.Int64("user_id", id), zap.Error(err))

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("删除用户失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)

		return fmt.Errorf("删除用户失败: %w", err)
	}

	s.logger.Info("User deleted successfully", zap.Int64("user_id", id))

	// 记录成功日志
	log.Status = "success"
	s.authService.LogUserOperation(ctx, log, nil)

	return nil
}

// UpdateUserStatus 更新用户状态
func (s *UserServiceImpl) UpdateUserStatus(ctx context.Context, id int64, status int) error {
	// 准备操作日志
	log := &UserOperationLog{
		Operation: "update_user_status",
		UserID:    id,
		Status:    "processing",
		Timestamp: time.Now(),
		UserAgent: "system",
		Details: map[string]interface{}{
			"new_status": status,
		},
	}

	// 验证状态值是否有效
	if status < 1 || status > 3 {
		log.Status = "failed"
		log.ErrorMsg = "无效的状态值"
		s.authService.LogUserOperation(ctx, log, nil)
		return errors.New("无效的状态值")
	}

	// 获取用户当前状态
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("用户不存在: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		s.logger.Error("User not found for status update", zap.Int64("user_id", id), zap.Error(err))
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 添加当前状态到日志
	if log.Details != nil {
		details := log.Details.(map[string]interface{})
		details["old_status"] = user.Status
		details["username"] = user.Username
	}

	err = s.userRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		s.logger.Error("Failed to update user status", zap.Int64("user_id", id), zap.Int("status", status), zap.Error(err))
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("更新用户状态失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		return fmt.Errorf("更新用户状态失败: %w", err)
	}

	s.logger.Info("User status updated successfully", zap.Int64("user_id", id), zap.Int("status", status))

	// 记录成功日志
	log.Status = "success"
	s.authService.LogUserOperation(ctx, log, nil)

	return nil
}

// LockUser 锁定用户
func (s *UserServiceImpl) LockUser(ctx context.Context, id int64) error {
	// 准备操作日志
	log := &UserOperationLog{
		Operation: "lock_user",
		UserID:    id,
		Status:    "processing",
		Timestamp: time.Now(),
		UserAgent: "system",
	}

	// 获取用户信息
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("用户不存在: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		s.logger.Error("User not found for locking", zap.Int64("user_id", id), zap.Error(err))
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 记录用户详情
	log.Details = map[string]interface{}{
		"username":        user.Username,
		"previous_status": user.Status,
	}

	// 如果用户已经被锁定，记录但不报错
	if user.IsLocked() {
		log.Status = "success"
		log.Details.(map[string]interface{})["note"] = "用户已经处于锁定状态"
		s.authService.LogUserOperation(ctx, log, nil)
		return nil
	}

	err = s.userRepo.UpdateStatus(ctx, id, 3) // 3 表示锁定状态
	if err != nil {
		s.logger.Error("Failed to lock user", zap.Int64("user_id", id), zap.Error(err))
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("锁定用户失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		return fmt.Errorf("锁定用户失败: %w", err)
	}

	s.logger.Info("User locked successfully", zap.Int64("user_id", id))

	// 记录成功日志
	log.Status = "success"
	s.authService.LogUserOperation(ctx, log, nil)

	return nil
}

// UnlockUser 解锁用户
func (s *UserServiceImpl) UnlockUser(ctx context.Context, id int64) error {
	// 准备操作日志
	log := &UserOperationLog{
		Operation: "unlock_user",
		UserID:    id,
		Status:    "processing",
		Timestamp: time.Now(),
		UserAgent: "system",
	}

	// 获取用户信息
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("用户不存在: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		s.logger.Error("User not found for unlock", zap.Int64("user_id", id), zap.Error(err))
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 记录用户详情
	log.Details = map[string]interface{}{
		"username":         user.Username,
		"previous_status":  user.Status,
		"login_fail_count": user.LoginFailCount,
	}

	// 如果用户未被锁定，记录但继续执行
	if !user.IsLocked() {
		log.Details.(map[string]interface{})["note"] = "用户未处于锁定状态"
	}

	// 重置登录失败计数
	err = s.userRepo.UpdateLoginFailCount(ctx, id, 0)
	if err != nil {
		s.logger.Error("Failed to reset login fail count", zap.Int64("user_id", id), zap.Error(err))
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("重置登录失败计数失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		return fmt.Errorf("重置登录失败计数失败: %w", err)
	}

	// 设置用户状态为正常
	err = s.userRepo.UpdateStatus(ctx, id, 1) // 1 表示正常状态
	if err != nil {
		s.logger.Error("Failed to unlock user", zap.Int64("user_id", id), zap.Error(err))
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("解锁用户失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		return fmt.Errorf("解锁用户失败: %w", err)
	}

	s.logger.Info("User unlocked successfully", zap.Int64("user_id", id))

	// 记录成功日志
	log.Status = "success"
	s.authService.LogUserOperation(ctx, log, nil)

	return nil
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

	// 获取用户列表
	users, err := s.userRepo.ListUsers(ctx, offset, pageSize)
	if err != nil {
		s.logger.Error("Failed to list users", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Error(err))
		return nil, 0, fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 获取用户总数
	total, err := s.userRepo.CountUsers(ctx)
	if err != nil {
		s.logger.Error("Failed to count users", zap.Error(err))
		return nil, 0, fmt.Errorf("获取用户总数失败: %w", err)
	}

	return users, total, nil
}

// BindWechat 绑定微信
func (s *UserServiceImpl) BindWechat(ctx context.Context, userID int64, openID, unionID string) error {
	// 准备操作日志
	log := &UserOperationLog{
		Operation: "bind_wechat",
		UserID:    userID,
		Status:    "processing",
		Timestamp: time.Now(),
		UserAgent: "system",
		Details: map[string]interface{}{
			"open_id":  openID,
			"union_id": unionID,
		},
	}

	// 检查该微信是否已被其他账号绑定
	existingUser, err := s.userRepo.GetUserByWechatOpenID(ctx, openID)
	if err == nil && existingUser.ID != userID {
		log.Status = "failed"
		log.ErrorMsg = "该微信账号已被其他用户绑定"
		log.Details.(map[string]interface{})["existing_user_id"] = existingUser.ID
		s.authService.LogUserOperation(ctx, log, nil)
		return errors.New("该微信账号已被其他用户绑定")
	}

	// 获取要绑定的用户
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("User not found for wechat binding", zap.Int64("user_id", userID), zap.Error(err))
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("用户不存在: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 如果用户已经绑定了微信，记录并更新
	if user.WechatOpenID != "" {
		log.Details.(map[string]interface{})["previous_open_id"] = user.WechatOpenID
		log.Details.(map[string]interface{})["previous_union_id"] = user.WechatUnionID
	}

	// 设置微信信息
	user.WechatOpenID = openID
	user.WechatUnionID = unionID

	// 更新用户
	_, err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		s.logger.Error("Failed to bind wechat", zap.Int64("user_id", userID), zap.String("open_id", openID), zap.Error(err))
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("绑定微信失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		return fmt.Errorf("绑定微信失败: %w", err)
	}

	s.logger.Info("Wechat bound successfully", zap.Int64("user_id", userID), zap.String("open_id", openID))

	// 记录成功日志
	log.Status = "success"
	s.authService.LogUserOperation(ctx, log, nil)

	return nil
}

// UnbindWechat 解绑微信
func (s *UserServiceImpl) UnbindWechat(ctx context.Context, userID int64) error {
	// 准备操作日志
	log := &UserOperationLog{
		Operation: "unbind_wechat",
		UserID:    userID,
		Status:    "processing",
		Timestamp: time.Now(),
		UserAgent: "system",
	}

	// 获取用户
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("User not found for wechat unbinding", zap.Int64("user_id", userID), zap.Error(err))
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("用户不存在: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 记录当前的微信信息
	log.Details = map[string]interface{}{
		"previous_open_id":  user.WechatOpenID,
		"previous_union_id": user.WechatUnionID,
	}

	// 如果用户没有绑定微信，记录并返回
	if user.WechatOpenID == "" && user.WechatUnionID == "" {
		log.Status = "success"
		log.Details.(map[string]interface{})["note"] = "用户未绑定微信，无需解绑"
		s.authService.LogUserOperation(ctx, log, nil)
		return nil
	}

	// 清除微信绑定信息
	user.WechatOpenID = ""
	user.WechatUnionID = ""

	// 更新用户
	_, err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		s.logger.Error("Failed to unbind wechat", zap.Int64("user_id", userID), zap.Error(err))
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("解绑微信失败: %v", err)
		s.authService.LogUserOperation(ctx, log, nil)
		return fmt.Errorf("解绑微信失败: %w", err)
	}

	s.logger.Info("Wechat unbound successfully", zap.Int64("user_id", userID))

	// 记录成功日志
	log.Status = "success"
	s.authService.LogUserOperation(ctx, log, nil)

	return nil
}

// GetUserPermissions 获取用户权限
func (s *UserServiceImpl) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	// 获取用户信息
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user for permissions", zap.Int64("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 根据用户角色返回对应权限
	// 实际项目中可能会从数据库或缓存中获取更复杂的权限配置
	var permissions []string

	switch user.Role {
	case 1: // 普通用户
		permissions = []string{
			"user:info:self",
			"user:update:self",
			"order:create",
			"order:list:self",
			"order:detail:self",
			"product:view",
		}
	case 2: // 管理员
		permissions = []string{
			"user:info:all",
			"user:update:all",
			"user:delete",
			"user:list",
			"order:list:all",
			"order:detail:all",
			"order:update",
			"product:manage",
			"system:manage",
		}
	default:
		permissions = []string{
			"user:info:self",
			"user:update:self",
		}
	}

	return permissions, nil
}
