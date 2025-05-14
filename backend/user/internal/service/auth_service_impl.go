package service

import (
	"context"
	"time"
	
	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/domain/valueobject"
	"shop/backend/user/pkg/jwt"
	
	"golang.org/x/crypto/bcrypt"
)

// AuthServiceImpl 认证服务实现
type AuthServiceImpl struct {
	repo      UserRepository
	jwtSecret string
	jwtExpire time.Duration
}

// NewAuthService 创建认证服务实例
func NewAuthService(repo UserRepository) AuthService {
	// 实际项目中，这些配置应该从配置文件或环境变量中获取
	return &AuthServiceImpl{
		repo:      repo,
		jwtSecret: "your-jwt-secret",
		jwtExpire: 24 * time.Hour,
	}
}

// Login 用户登录
func (s *AuthServiceImpl) Login(ctx context.Context, credential *valueobject.Credential) (*valueobject.LoginResponse, error) {
	var user *entity.User
	var err error
	
	// 根据凭证类型处理不同的登录方式
	switch credential.Type {
	case valueobject.CredentialTypePassword:
		// 通过手机号获取用户
		user, err = s.repo.GetByMobile(ctx, credential.Username)
		if err != nil {
			return nil, err
		}
		
		if user == nil {
			return nil, ErrUserNotFound
		}
		
		// 检查用户状态
		if user.IsDisabled() {
			return nil, ErrUserDisabled
		}
		
		if user.IsLocked() {
			return nil, ErrUserLocked
		}
		
		// 验证密码
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credential.Password))
		if err != nil {
			// 增加登录失败计数
			user.IncrementLoginFailCount()
			
			// 如果连续失败10次，则锁定账号
			if user.LoginFailCount >= 10 {
				user.Lock()
			}
			
			// 更新用户状态
			if updateErr := s.repo.Update(ctx, user); updateErr != nil {
				// 记录日志但不影响主流程
			}
			
			return nil, ErrInvalidCredentials
		}
		
		// 登录成功，重置登录失败计数并更新最后登录时间
		user.ResetLoginFailCount()
		user.UpdateLastLogin()
		if updateErr := s.repo.Update(ctx, user); updateErr != nil {
			// 记录日志但不影响主流程
		}
		
	case valueobject.CredentialTypeSMS:
		// TODO: 实现短信验证码登录
		// 验证短信验证码
		// 获取用户信息
		return nil, errors.New("SMS login not implemented")
		
	case valueobject.CredentialTypeWechat:
		// TODO: 实现微信登录
		// 验证微信授权信息
		// 获取或创建用户
		return nil, errors.New("Wechat login not implemented")
		
	default:
		return nil, errors.New("unsupported credential type")
	}
	
	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.ID, user.Mobile, user.Role, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, err
	}
	
	// 生成刷新令牌(有效期长一些)
	refreshToken, err := jwt.GenerateToken(user.ID, user.Mobile, user.Role, s.jwtSecret+"_refresh", s.jwtExpire*30)
	if err != nil {
		return nil, err
	}
	
	// 构造登录响应
	response := &valueobject.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpireTime:   time.Now().Add(s.jwtExpire).Unix(),
		UserInfo: map[string]interface{}{
			"id":       user.ID,
			"mobile":   user.Mobile,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	}
	
	return response, nil
}

// ValidateToken 验证令牌有效性
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (int64, string, error) {
	// 解析和验证JWT令牌
	claims, err := jwt.ParseToken(token, s.jwtSecret)
	if err != nil {
		return 0, "", err
	}
	
	// 从claims中获取用户ID和角色
	userID := claims.UserId
	role := claims.Role
	
	// 验证用户是否存在且状态正常（可选，但建议进行）
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return 0, "", err
	}
	
	if user == nil {
		return 0, "", ErrUserNotFound
	}
	
	if !user.IsValid() {
		if user.IsDisabled() {
			return 0, "", ErrUserDisabled
		}
		if user.IsLocked() {
			return 0, "", ErrUserLocked
		}
	}
	
	return userID, role, nil
}

// RefreshToken 刷新令牌
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*valueobject.LoginResponse, error) {
	// 解析和验证刷新令牌
	claims, err := jwt.ParseToken(refreshToken, s.jwtSecret+"_refresh")
	if err != nil {
		return nil, err
	}
	
	// 从claims中获取用户ID
	userID := claims.UserId
	
	// 获取用户信息
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, ErrUserNotFound
	}
	
	if !user.IsValid() {
		if user.IsDisabled() {
			return nil, ErrUserDisabled
		}
		if user.IsLocked() {
			return nil, ErrUserLocked
		}
	}
	
	// 生成新的JWT令牌
	newToken, err := jwt.GenerateToken(user.ID, user.Mobile, user.Role, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, err
	}
	
	// 生成新的刷新令牌
	newRefreshToken, err := jwt.GenerateToken(user.ID, user.Mobile, user.Role, s.jwtSecret+"_refresh", s.jwtExpire*30)
	if err != nil {
		return nil, err
	}
	
	// 构造响应
	response := &valueobject.LoginResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		ExpireTime:   time.Now().Add(s.jwtExpire).Unix(),
		UserInfo: map[string]interface{}{
			"id":       user.ID,
			"mobile":   user.Mobile,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	}
	
	return response, nil
}

// Logout 退出登录
func (s *AuthServiceImpl) Logout(ctx context.Context, token string) error {
	// 解析token获取用户ID
	claims, err := jwt.ParseToken(token, s.jwtSecret)
	if err != nil {
		return err
	}
	
	// 在实际应用中，可能需要将token加入黑名单，防止被再次使用
	// 这里简化实现，不做额外处理
	
	return nil
}
