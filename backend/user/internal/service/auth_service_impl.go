package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/domain/valueobject"
	"shop/backend/user/internal/repository"
	"shop/backend/user/pkg/jwt"
)

// AuthServiceImpl 认证服务实现
type AuthServiceImpl struct {
	userRepo     repository.UserRepository
	jwtUtil      *jwt.JWTUtil
	logger       *zap.Logger
	accessTTL    time.Duration
	refreshTTL   time.Duration
	maxLoginFail int
	mongoClient  *mongo.Client // 添加MongoDB客户端
}

// NewAuthService 创建认证服务
func NewAuthService(
	userRepo repository.UserRepository,
	jwtUtil *jwt.JWTUtil,
	logger *zap.Logger,
	accessTTL, refreshTTL time.Duration,
	maxLoginFail int,
	mongoClient *mongo.Client, // 添加MongoDB客户端参数
) AuthService {
	return &AuthServiceImpl{
		userRepo:     userRepo,
		jwtUtil:      jwtUtil,
		logger:       logger,
		accessTTL:    accessTTL,
		refreshTTL:   refreshTTL,
		maxLoginFail: maxLoginFail,
		mongoClient:  mongoClient, // 初始化MongoDB客户端
	}
}

// Login 用户登录
func (s *AuthServiceImpl) Login(ctx context.Context, username, password string) (*entity.User, *valueobject.Credential, error) {
	var user *entity.User
	var err error

	// 支持使用用户名、手机号或邮箱登录
	if len(username) == 11 && username[0] == '1' { // 简单判断是否为手机号
		user, err = s.userRepo.GetUserByPhone(ctx, username)
	} else if isEmail(username) {
		user, err = s.userRepo.GetUserByEmail(ctx, username)
	} else {
		user, err = s.userRepo.GetUserByUsername(ctx, username)
	}

	// 准备操作日志
	log := &UserOperationLog{
		Operation: "login",
		UserAgent: "system", // 通常从请求中获取，这里是系统调用默认值
		Timestamp: time.Now(),
	}

	if err != nil {
		s.logger.Error("Failed to find user", zap.String("username", username), zap.Error(err))

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("用户不存在: %v", err)
		log.Details = map[string]interface{}{
			"login_type": "username",
			"username":   username,
		}

		// 这里使用0作为用户ID，表示未找到用户
		log.UserID = 0
		s.logUserOperationInternal(ctx, log)

		return nil, nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 更新日志中的用户ID
	log.UserID = user.ID

	// 验证用户状态
	if !user.IsActive() {
		s.logger.Warn("Login attempt on inactive account", zap.Int64("user_id", user.ID))

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = "账号已被禁用或锁定"
		log.Details = map[string]interface{}{
			"user_status": user.Status,
		}
		s.logUserOperationInternal(ctx, log)

		if user.IsLocked() {
			return nil, nil, errors.New("账号已被锁定，请联系客服")
		}
		return nil, nil, errors.New("账号已被禁用")
	}

	// 验证密码
	valid, err := s.userRepo.VerifyPassword(ctx, user.ID, password)
	if err != nil || !valid {
		// 增加登录失败次数
		user.IncrementLoginFailCount()
		err := s.userRepo.UpdateLoginFailCount(ctx, user.ID, user.LoginFailCount)
		if err != nil {
			s.logger.Error("Failed to update login fail count", zap.Int64("user_id", user.ID), zap.Error(err))
		}

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = "密码错误"
		log.Details = map[string]interface{}{
			"login_fail_count": user.LoginFailCount,
		}
		s.logUserOperationInternal(ctx, log)

		// 检查是否超过最大失败次数
		if user.LoginFailCount >= s.maxLoginFail {
			user.Lock()
			err := s.userRepo.UpdateStatus(ctx, user.ID, user.Status)
			if err != nil {
				s.logger.Error("Failed to lock user", zap.Int64("user_id", user.ID), zap.Error(err))
			}
			s.logger.Warn("User account locked due to too many login failures", zap.Int64("user_id", user.ID))
			return nil, nil, errors.New("登录失败次数过多，账号已被锁定")
		}

		s.logger.Info("Login failed: invalid password", zap.Int64("user_id", user.ID), zap.Int("fail_count", user.LoginFailCount))
		return nil, nil, errors.New("用户名或密码不正确")
	}

	// 重置登录失败次数
	if user.LoginFailCount > 0 {
		user.ResetLoginFailCount()
		err = s.userRepo.UpdateLoginFailCount(ctx, user.ID, user.LoginFailCount)
		if err != nil {
			s.logger.Error("Failed to reset login fail count", zap.Int64("user_id", user.ID), zap.Error(err))
			// 不阻止登录流程继续
		}
	}

	// 登录成功，更新最后登录时间
	user.UpdateLastLogin()
	err = s.userRepo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		s.logger.Error("Failed to update last login time", zap.Int64("user_id", user.ID), zap.Error(err))
		// 不阻止登录流程继续
	}

	// 生成令牌
	credential, err := s.GenerateToken(ctx, user)
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Int64("user_id", user.ID), zap.Error(err))

		// 记录失败日志
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("生成令牌失败: %v", err)
		s.logUserOperationInternal(ctx, log)

		return nil, nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	// 记录成功日志
	log.Status = "success"
	log.Details = map[string]interface{}{
		"last_login": user.LastLoginAt,
	}
	s.logUserOperationInternal(ctx, log)

	return user, credential, nil
}

// RegisterUser 注册用户
func (s *AuthServiceImpl) RegisterUser(ctx context.Context, user *entity.User, password string) (*entity.User, error) {
	// 准备操作日志
	log := &UserOperationLog{
		Operation: "register",
		UserAgent: "system", // 通常从请求中获取
		UserID:    0,        // 注册前没有用户ID
		Timestamp: time.Now(),
	}

	// 验证用户数据
	if user.Username == "" {
		log.Status = "failed"
		log.ErrorMsg = "用户名不能为空"
		s.logUserOperationInternal(ctx, log)
		return nil, errors.New("用户名不能为空")
	}

	if password == "" {
		log.Status = "failed"
		log.ErrorMsg = "密码不能为空"
		s.logUserOperationInternal(ctx, log)
		return nil, errors.New("密码不能为空")
	}

	// 检查用户名是否已存在
	_, err := s.userRepo.GetUserByUsername(ctx, user.Username)
	if err == nil {
		log.Status = "failed"
		log.ErrorMsg = "用户名已被使用"
		log.Details = map[string]interface{}{
			"username": user.Username,
		}
		s.logUserOperationInternal(ctx, log)
		return nil, errors.New("用户名已被使用")
	}

	// 如果提供了手机号，检查是否已注册
	if user.Phone != "" {
		_, err = s.userRepo.GetUserByPhone(ctx, user.Phone)
		if err == nil {
			log.Status = "failed"
			log.ErrorMsg = "手机号已被注册"
			log.Details = map[string]interface{}{
				"phone": user.Phone,
			}
			s.logUserOperationInternal(ctx, log)
			return nil, errors.New("手机号已被注册")
		}
	}

	// 如果提供了邮箱，检查是否已注册
	if user.Email != "" {
		_, err = s.userRepo.GetUserByEmail(ctx, user.Email)
		if err == nil {
			log.Status = "failed"
			log.ErrorMsg = "邮箱已被注册"
			log.Details = map[string]interface{}{
				"email": user.Email,
			}
			s.logUserOperationInternal(ctx, log)
			return nil, errors.New("邮箱已被注册")
		}
	}

	// 设置默认值
	if user.Nickname == "" {
		user.Nickname = user.Username
	}
	if user.Gender == "" {
		user.Gender = "unknown"
	}
	if user.Status == 0 {
		user.Status = 1 // 正常状态
	}
	if user.Role == 0 {
		user.Role = 1 // 普通用户
	}

	// 哈希密码
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("密码加密失败: %v", err)
		s.logUserOperationInternal(ctx, log)
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	createdUser, err := s.userRepo.CreateUser(ctx, user, hashedPassword)
	if err != nil {
		s.logger.Error("Failed to create user", zap.String("username", user.Username), zap.Error(err))
		log.Status = "failed"
		log.ErrorMsg = fmt.Sprintf("创建用户失败: %v", err)
		s.logUserOperationInternal(ctx, log)
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	s.logger.Info("User registered successfully", zap.Int64("user_id", createdUser.ID), zap.String("username", createdUser.Username))

	// 记录成功日志
	log.Status = "success"
	log.UserID = createdUser.ID
	log.Details = map[string]interface{}{
		"user_id":  createdUser.ID,
		"username": createdUser.Username,
	}
	s.logUserOperationInternal(ctx, log)

	return createdUser, nil
}

// LogUserOperation 记录用户操作日志
func (s *AuthServiceImpl) LogUserOperation(ctx context.Context, log *UserOperationLog, req *http.Request) error {
	// 从请求中获取更多信息
	if req != nil {
		log.IP = getClientIP(req)
		log.UserAgent = req.UserAgent()
	}

	// 确保时间戳存在
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	return s.logUserOperationInternal(ctx, log)
}

// logUserOperationInternal 内部记录用户操作日志方法
func (s *AuthServiceImpl) logUserOperationInternal(ctx context.Context, log *UserOperationLog) error {
	if s.mongoClient == nil {
		s.logger.Warn("MongoDB client not available, skipping log")
		return errors.New("MongoDB client not available")
	}

	// 创建MongoDB对象
	collection := s.mongoClient.Database("user_logs").Collection("user_operations")

	// 确保时间戳存在
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 脱敏处理请求数据（如果有）
	if len(log.RequestData) > 0 {
		var requestData map[string]interface{}
		if err := json.Unmarshal([]byte(log.RequestData), &requestData); err == nil {
			// 脱敏密码字段
			if password, exists := requestData["password"]; exists {
				requestData["password"] = "******"
			}
			// 脱敏其他敏感字段...

			// 重新序列化
			if sanitizedData, err := json.Marshal(requestData); err == nil {
				log.RequestData = string(sanitizedData)
			}
		}
	}

	// 插入日志
	_, err := collection.InsertOne(ctx, log)
	if err != nil {
		s.logger.Error("Failed to insert user operation log", zap.Error(err))
		return err
	}

	return nil
}

// RefreshToken 刷新访问令牌
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*valueobject.Credential, error) {
	// 验证刷新令牌
	userID, err := s.userRepo.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Info("Invalid refresh token", zap.Error(err))
		return nil, errors.New("刷新令牌无效或已过期")
	}

	// 获取用户信息
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user for refresh token", zap.Int64("user_id", userID), zap.Error(err))
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 验证用户状态
	if !user.IsActive() {
		s.logger.Warn("Refresh token attempt on inactive account", zap.Int64("user_id", user.ID))
		return nil, errors.New("账号已被禁用或锁定")
	}

	// 生成新的令牌
	credential, err := s.GenerateToken(ctx, user)
	if err != nil {
		s.logger.Error("Failed to generate new token", zap.Int64("user_id", user.ID), zap.Error(err))
		return nil, fmt.Errorf("生成新令牌失败: %w", err)
	}

	// 使旧的刷新令牌失效
	err = s.userRepo.InvalidateRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Error("Failed to invalidate old refresh token", zap.Int64("user_id", user.ID), zap.Error(err))
		// 不阻止刷新流程继续
	}

	return credential, nil
}

// GenerateToken 生成用户令牌
func (s *AuthServiceImpl) GenerateToken(ctx context.Context, user *entity.User) (*valueobject.Credential, error) {
	// 生成访问令牌
	accessToken, err := s.jwtUtil.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	// 生成刷新令牌
	refreshToken, err := s.jwtUtil.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	// 存储令牌到缓存
	err = s.userRepo.StoreToken(ctx, accessToken, user.ID, s.accessTTL)
	if err != nil {
		s.logger.Error("Failed to store access token", zap.Int64("user_id", user.ID), zap.Error(err))
		return nil, fmt.Errorf("存储访问令牌失败: %w", err)
	}

	err = s.userRepo.StoreRefreshToken(ctx, refreshToken, user.ID, s.refreshTTL)
	if err != nil {
		s.logger.Error("Failed to store refresh token", zap.Int64("user_id", user.ID), zap.Error(err))
		return nil, fmt.Errorf("存储刷新令牌失败: %w", err)
	}

	return &valueobject.Credential{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.accessTTL,
	}, nil
}

// ValidateToken 验证访问令牌
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (*valueobject.TokenClaims, error) {
	// 解析令牌
	claims, err := s.jwtUtil.ParseToken(token)
	if err != nil {
		s.logger.Info("Invalid token format", zap.Error(err))
		return nil, errors.New("令牌格式无效")
	}

	// 验证令牌类型
	if claims.Type != jwt.AccessToken {
		s.logger.Info("Invalid token type", zap.String("type", string(claims.Type)))
		return nil, errors.New("令牌类型无效")
	}

	// 验证令牌是否在缓存中
	userID, err := s.userRepo.VerifyToken(ctx, token)
	if err != nil {
		s.logger.Info("Token not found in cache", zap.Error(err))
		return nil, errors.New("令牌无效或已过期")
	}

	// 验证令牌中的用户ID与缓存中的是否一致
	if userID != claims.UserID {
		s.logger.Warn("Token user ID mismatch", zap.Int64("claims_user_id", claims.UserID), zap.Int64("cache_user_id", userID))
		return nil, errors.New("令牌用户ID不匹配")
	}

	// 转换JWT claims为应用claims
	return &valueobject.TokenClaims{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Role:      claims.Role,
		IssuedAt:  claims.IssuedAt.Unix(),
		ExpiresAt: claims.ExpiresAt.Unix(),
		TokenType: string(claims.Type),
		Scopes:    []string{}, // JWT claims中没有存储scope信息
	}, nil
}
