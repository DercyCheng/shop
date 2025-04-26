package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// User 数据库用户模型
type User struct {
	ID             int64      `gorm:"primaryKey;autoIncrement"`
	Username       string     `gorm:"type:varchar(50);uniqueIndex"`
	Password       string     `gorm:"type:varchar(100);not null"`
	Nickname       string     `gorm:"type:varchar(50)"`
	Email          string     `gorm:"type:varchar(100);index"`
	Phone          string     `gorm:"type:varchar(20);uniqueIndex"`
	Avatar         string     `gorm:"type:varchar(255)"`
	Gender         string     `gorm:"type:varchar(10);default:'unknown'"`
	Birthday       *time.Time `gorm:"type:datetime"`
	Status         int        `gorm:"type:tinyint;default:1;index"` // 1: 正常, 2: 禁用, 3: 锁定
	Role           int        `gorm:"type:tinyint;default:1"`       // 1: 普通用户, 2: 管理员
	LoginFailCount int        `gorm:"type:int;default:0"`
	LastLoginAt    *time.Time `gorm:"type:datetime"`
	WechatOpenID   string     `gorm:"type:varchar(50);index"`
	WechatUnionID  string     `gorm:"type:varchar(50)"`
	SessionID      string     `gorm:"type:varchar(100)"`
	CreatedAt      time.Time  `gorm:"type:datetime;index"`
	UpdatedAt      time.Time  `gorm:"type:datetime"`
	DeletedAt      *time.Time `gorm:"type:datetime;index"`
}

// UserDAO 用户数据访问对象
type UserDAO struct {
	db *gorm.DB
}

// NewUserDAO 创建用户DAO
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

// Create 创建用户
func (dao *UserDAO) Create(ctx context.Context, user *User) error {
	return dao.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (dao *UserDAO) Update(ctx context.Context, user *User) error {
	return dao.db.WithContext(ctx).Save(user).Error
}

// FindByID 通过ID查找用户
func (dao *UserDAO) FindByID(ctx context.Context, id int64) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 通过用户名查找用户
func (dao *UserDAO) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone 通过手机号查找用户
func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 通过邮箱查找用户
func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByWechatOpenID 通过微信OpenID查找用户
func (dao *UserDAO) FindByWechatOpenID(ctx context.Context, openID string) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Delete 删除用户（软删除）
func (dao *UserDAO) Delete(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Where("id = ?", id).Delete(&User{}).Error
}

// UpdatePassword 更新用户密码
func (dao *UserDAO) UpdatePassword(ctx context.Context, id int64, password string) error {
	return dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		Update("password", password).Error
}

// UpdateStatus 更新用户状态
func (dao *UserDAO) UpdateStatus(ctx context.Context, id int64, status int) error {
	return dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		Update("status", status).Error
}

// UpdateLoginFailCount 更新登录失败次数
func (dao *UserDAO) UpdateLoginFailCount(ctx context.Context, id int64, count int) error {
	return dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		Update("login_fail_count", count).Error
}

// UpdateLastLogin 更新最后登录时间
func (dao *UserDAO) UpdateLastLogin(ctx context.Context, id int64) error {
	now := time.Now()
	return dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at":    now,
			"login_fail_count": 0,
		}).Error
}

// ListUsers 获取用户列表
func (dao *UserDAO) ListUsers(ctx context.Context, offset, limit int) ([]*User, error) {
	var users []*User
	err := dao.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

// CountUsers 获取用户总数
func (dao *UserDAO) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&User{}).Count(&count).Error
	return count, err
}
