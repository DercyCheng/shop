package valueobject

// Credential 用户凭证值对象
type Credential struct {
	Username string
	Password string
	Type     CredentialType
	Code     string  // 验证码（可选，用于短信登录等）
}

// CredentialType 凭证类型
type CredentialType int

// 凭证类型枚举
const (
	CredentialTypePassword CredentialType = iota // 用户名密码
	CredentialTypeSMS                          // 短信验证码
	CredentialTypeWechat                       // 微信授权
)

// LoginResponse 登录响应值对象
type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	ExpireTime   int64       `json:"expire_time"`
	UserInfo     interface{} `json:"user_info"`
}
