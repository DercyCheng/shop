package entity

import (
	"time"
)

// User 用户实体
type User struct {
	ID            int64      `json:"id"`
	Mobile        string     `json:"mobile"`
	Password      string     `json:"-"` // 不对外暴露密码
	Nickname      string     `json:"nickname"`
	Avatar        string     `json:"avatar"`
	Birthday      *time.Time `json:"birthday"`
	Gender        string     `json:"gender"`
	Role          int        `json:"role"`
	Status        int        `json:"status"`
	LoginFailCount int        `json:"login_fail_count"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	WechatOpenID  string     `json:"wechat_open_id"`
	WechatUnionID string     `json:"wechat_union_id"`
	SessionID     string     `json:"session_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

// UserStatusEnum 用户状态枚举
const (
	UserStatusNormal  = 1 // 正常
	UserStatusDisabled = 2 // 禁用
	UserStatusLocked  = 3 // 锁定
)

// UserRoleEnum 用户角色枚举
const (
	UserRoleNormal = 1 // 普通用户
	UserRoleAdmin  = 2 // 管理员
)

// IsValid 检查用户状态是否有效
func (u *User) IsValid() bool {
	return u.Status == UserStatusNormal
}

// IsLocked 检查用户是否被锁定
func (u *User) IsLocked() bool {
	return u.Status == UserStatusLocked
}

// IsDisabled 检查用户是否被禁用
func (u *User) IsDisabled() bool {
	return u.Status == UserStatusDisabled
}

// IncrementLoginFailCount 增加登录失败计数
func (u *User) IncrementLoginFailCount() {
	u.LoginFailCount++
}

// ResetLoginFailCount 重置登录失败计数
func (u *User) ResetLoginFailCount() {
	u.LoginFailCount = 0
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// Lock 锁定用户账号
func (u *User) Lock() {
	u.Status = UserStatusLocked
}

// Unlock 解锁用户账号
func (u *User) Unlock() {
	u.Status = UserStatusNormal
}
