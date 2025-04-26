package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/domain/valueobject"
	"shop/backend/user/internal/repository"
)

// 常量定义
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// AuthServiceImpl 认证服务实现
type AuthServiceImpl struct {
	userRepo        repository.UserRepository
	redisClient     *redis.Client
	jwtSecret       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	verificationTTL time.Duration
}

// NewAuthService 创建认证服务
func NewAuthService(
	userRepo repository.UserRepository,
	redisClient *redis.Client,
	jwtSecret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	verificationTTL time.Duration,
) AuthService {
	return &AuthServiceImpl{
		userRepo:        userRepo,
		redisClient:     redisClient,
		jwtSecret:       jwtSecret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		verificationTTL: verificationTTL,
	}
}

// Login 用户登录
func (s *AuthServiceImpl) Login(ctx context.Context, username, password string) (*entity.User, *valueobject.Credential, error) {
	// 首先尝试通过用户名查找用户
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		// 如果用户名不存在，尝试通过手机号查找
		user, err = s.userRepo.GetUserByPhone(ctx, username)
		if err != nil {
			// 如果手机号也不存在，尝试通过邮箱查找
			user, err = s.userRepo.GetUserByEmail(ctx, username)
			if err != nil {
				return nil, nil, errors.New("用户不存在")
			}
		}
	}

	// 检查用户状态
	if user.Status == 2 {
		return nil, nil, errors.New("账号已被禁用")
	}
	if user.Status == 3 {
		return nil, nil, errors.New("账号已被锁定")
	}

	// 验证密码
	valid, err := s.userRepo.VerifyPassword(ctx, user.ID, password)
	if err != nil {
		return nil, nil, fmt.Errorf("验证密码失败: %w", err)
	}

	// 如果密码不正确
	if !valid {
		// 增加登录失败次数
		user.IncrementLoginFailCount()
		err = s.userRepo.UpdateLoginFailCount(ctx, user.ID, user.LoginFailCount)
		if err != nil {
			return nil, nil, fmt.Errorf("更新登录失败次数失败: %w", err)
		}

		// 如果失败次数超过阈值，锁定账号
		if user.LoginFailCount >= 5 {
			user.Lock()
			err = s.userRepo.UpdateStatus(ctx, user.ID, user.Status)
			if err != nil {
				return nil, nil, fmt.Errorf("锁定账号失败: %w", err)
			}
			return nil, nil, errors.New("密码错误次数过多，账号已被锁定")
		}

		return nil, nil, errors.New("密码不正确")
	}

	// 如果密码正确，重置登录失败次数并更新最后登录时间
	user.ResetLoginFailCount()
	user.UpdateLastLogin()
	err = s.userRepo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("更新登录信息失败: %w", err)
	}

	// 生成令牌
	credential, err := s.GenerateToken(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return user, credential, nil
}

// Logout 用户登出
func (s *AuthServiceImpl) Logout(ctx context.Context, token string) error {
	// 验证令牌
	claims, err := s.ValidateToken(ctx, token)
	if err != nil {
		return errors.New("无效的令牌")
	}

	// 使令牌失效
	if err := s.userRepo.InvalidateToken(ctx, token); err != nil {
		return fmt.Errorf("使令牌失效失败: %w", err)
	}

	// 记录登出信息
	// TODO: 记录日志

	return nil
}

// RefreshToken 刷新访问令牌
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*valueobject.Credential, error) {
	// 验证刷新令牌
	userID, err := s.userRepo.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 获取用户信息
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账号状态异常")
	}

	// 使旧的刷新令牌失效
	if err := s.userRepo.InvalidateRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("使刷新令牌失效失败: %w", err)
	}

	// 生成新的令牌
	credential, err := s.GenerateToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return credential, nil
}

// GenerateToken 生成用户令牌
func (s *AuthServiceImpl) GenerateToken(ctx context.Context, user *entity.User) (*valueobject.Credential, error) {
	now := time.Now()
	accessClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"type":     TokenTypeAccess,
		"iat":      now.Unix(),
		"exp":      now.Add(s.accessTokenTTL).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    TokenTypeRefresh,
		"iat":     now.Unix(),
		"exp":     now.Add(s.refreshTokenTTL).Unix(),
	}

	// 生成访问令牌
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	// 生成刷新令牌
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	// 存储令牌
	if err := s.userRepo.StoreToken(ctx, accessTokenString, user.ID, s.accessTokenTTL); err != nil {
		return nil, fmt.Errorf("存储访问令牌失败: %w", err)
	}
	if err := s.userRepo.StoreRefreshToken(ctx, refreshTokenString, user.ID, s.refreshTokenTTL); err != nil {
		return nil, fmt.Errorf("存储刷新令牌失败: %w", err)
	}

	return &valueobject.Credential{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    s.accessTokenTTL,
	}, nil
}

// ValidateToken 验证访问令牌
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*valueobject.TokenClaims, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌有效性
	if !token.Valid {
		return nil, errors.New("无效的令牌")
	}

	// 获取声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("无效的令牌声明")
	}

	// 验证令牌类型
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != TokenTypeAccess {
		return nil, errors.New("无效的令牌类型")
	}

	// 验证令牌是否在Redis中有效
	userID, err := s.userRepo.VerifyToken(ctx, tokenString)
	if err != nil {
		return nil, errors.New("令牌已失效")
	}

	// 构建Token声明
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("无效的用户ID")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("无效的用户名")
	}

	roleFloat, ok := claims["role"].(float64)
	if !ok {
		return nil, errors.New("无效的角色信息")
	}

	issuedAtFloat, ok := claims["iat"].(float64)
	if !ok {
		return nil, errors.New("无效的签发时间")
	}

	expiresAtFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("无效的过期时间")
	}

	return &valueobject.TokenClaims{
		UserID:    int64(userIDFloat),
		Username:  username,
		Role:      int(roleFloat),
		IssuedAt:  int64(issuedAtFloat),
		ExpiresAt: int64(expiresAtFloat),
		TokenType: tokenType,
	}, nil
}

// RegisterUser 注册用户
func (s *AuthServiceImpl) RegisterUser(ctx context.Context, user *entity.User, password string) (*entity.User, error) {
	// 哈希密码
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("密码哈希失败: %w", err)
	}

	// 设置用户默认信息
	user.Status = 1 // 正常状态
	user.Role = 1   // 普通用户

	// 创建用户
	createdUser, err := s.userRepo.CreateUser(ctx, user, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return createdUser, nil
}

// VerifyPassword 验证用户密码
func (s *AuthServiceImpl) VerifyPassword(ctx context.Context, userID int64, password string) (bool, error) {
	return s.userRepo.VerifyPassword(ctx, userID, password)
}

// HashPassword 哈希用户密码
func (s *AuthServiceImpl) HashPassword(password string) (string, error) {
	// 生成随机盐值
	salt := generateSalt(16)

	// 哈希密码（使用MD5+盐值）
	hasher := md5.New()
	_, err := hasher.Write([]byte(password + salt))
	if err != nil {
		return "", err
	}

	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	// 返回格式: hashedPassword:salt
	return fmt.Sprintf("%s:%s", hashedPassword, salt), nil
}

// GenerateVerificationCode 生成验证码
func (s *AuthServiceImpl) GenerateVerificationCode(ctx context.Context, phone string, codeType string) (string, error) {
	// 生成随机验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 验证码键
	key := fmt.Sprintf("verification:%s:%s", codeType, phone)

	// 存储验证码到Redis
	err := s.redisClient.Set(ctx, key, code, s.verificationTTL).Err()
	if err != nil {
		return "", fmt.Errorf("存储验证码失败: %w", err)
	}

	// TODO: 发送验证码到用户手机（实际应用中这里会调用短信服务）

	return code, nil
}

// VerifyVerificationCode 验证验证码
func (s *AuthServiceImpl) VerifyVerificationCode(ctx context.Context, phone, code, codeType string) (bool, error) {
	// 验证码键
	key := fmt.Sprintf("verification:%s:%s", codeType, phone)

	// 从Redis获取验证码
	storedCode, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, errors.New("验证码不存在或已过期")
		}
		return false, fmt.Errorf("获取验证码失败: %w", err)
	}

	// 验证码校验
	if storedCode != code {
		return false, nil
	}

	// 验证成功后使验证码失效
	err = s.redisClient.Del(ctx, key).Err()
	if err != nil {
		// 这里只记录错误，不影响验证结果
		fmt.Printf("删除验证码失败: %v", err)
	}

	return true, nil
}

// 生成随机盐值
func generateSalt(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	salt := make([]byte, length)
	for i := range salt {
		salt[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(salt)
}
