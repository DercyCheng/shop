package valueobject

import "time"

// Credential 用户凭证值对象
type Credential struct {
	AccessToken  string        // JWT访问令牌
	RefreshToken string        // 刷新令牌
	ExpiresIn    time.Duration // 访问令牌过期时间
}

// TokenClaims JWT令牌声明
type TokenClaims struct {
	UserID    int64    // 用户ID
	Username  string   // 用户名
	Role      int      // 角色
	IssuedAt  int64    // 签发时间
	ExpiresAt int64    // 过期时间
	TokenType string   // 令牌类型: access 或 refresh
	Scopes    []string // 权限范围
}
