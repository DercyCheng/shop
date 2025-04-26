package entity

import (
	"time"
)

// User 用户实体
type User struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Nickname string    `json:"nickname"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Password string    `json:"-"` // 密码不输出到JSON
	Avatar   string    `json:"avatar"`
	Gender   string    `json:"gender"` // "male", "female", "unknown"
	Birthday time.Time `json:"birthday"`
	Status   int       `json:"status"` // 1: 正常, 2: 禁用, 3: 锁定

	// 第三方账号信息
	WechatOpenID  string `json:"-"`
	WechatUnionID string `json:"-"`

	// 登录相关信息
	LoginFailCount int       `json:"-"`
	LastLoginAt    time.Time `json:"last_login_at"`
	SessionID      string    `json:"-"`

	// 角色信息
	Role int `json:"role"` // 1: 普通用户, 2: 管理员

	// 审计字段
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
}

// IsActive 检查用户是否处于活动状态
func (u *User) IsActive() bool {
	return u.Status == 1
}

// IsLocked 检查用户是否被锁定
func (u *User) IsLocked() bool {
	return u.Status == 3
}

// IsAdmin 检查用户是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == 2
}

// IncrementLoginFailCount 增加登录失败次数
func (u *User) IncrementLoginFailCount() {
	u.LoginFailCount++
}

// ResetLoginFailCount 重置登录失败次数
func (u *User) ResetLoginFailCount() {
	u.LoginFailCount = 0
}

// Lock 锁定用户账号
func (u *User) Lock() {
	u.Status = 3
}

// Disable 禁用用户账号
func (u *User) Disable() {
	u.Status = 2
}

// Enable 启用用户账号
func (u *User) Enable() {
	u.Status = 1
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin() {
	u.LastLoginAt = time.Now()
}
