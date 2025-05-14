package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// UserFav 用户收藏
type UserFav struct {
	ID            int64      `json:"id"`
	UserID        int64      `json:"user_id" gorm:"column:user"`
	GoodsID       int64      `json:"goods_id" gorm:"column:goods"`
	CategoryID    int64      `json:"category_id" gorm:"column:category_id"`
	Remark        string     `json:"remark"`
	PriceWhenFav  float64    `json:"price_when_fav"`
	Notification  bool       `json:"notification"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

// TableName 指定表名
func (UserFav) TableName() string {
	return "user_fav"
}

// Address 用户地址
type Address struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id" gorm:"column:user"`
	Province     string     `json:"province"`
	City         string     `json:"city"`
	District     string     `json:"district"`
	Address      string     `json:"address"`
	SignerName   string     `json:"signer_name"`
	SignerMobile string     `json:"signer_mobile"`
	IsDefault    bool       `json:"is_default"`
	Label        string     `json:"label"`
	Postcode     string     `json:"postcode"`
	UsageCount   int        `json:"usage_count"`
	LastUsedAt   *time.Time `json:"last_used_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

// TableName 指定表名
func (Address) TableName() string {
	return "address"
}

// JsonArray 自定义JSON数组类型
type JsonArray []string

// Value 实现driver.Valuer接口，用于存储到数据库
func (j JsonArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口，用于从数据库加载
func (j *JsonArray) Scan(value interface{}) error {
	if value == nil {
		*j = JsonArray{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, j)
}

// UserFeedback 用户反馈
type UserFeedback struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id" gorm:"column:user"`
	FeedbackType int        `json:"feedback_type"`
	Subject      string     `json:"subject"`
	Content      string     `json:"content"`
	FileURLs     JsonArray  `json:"file_urls" gorm:"column:file_urls"`
	Status       int        `json:"status"`
	OrderSn      string     `json:"order_sn"`
	AdminReply   string     `json:"admin_reply"`
	ReplyAt      *time.Time `json:"reply_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

// TableName 指定表名
func (UserFeedback) TableName() string {
	return "user_feedback"
}

// BrowsingHistory 浏览历史
type BrowsingHistory struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id" gorm:"column:user"`
	GoodsID   int64     `json:"goods_id" gorm:"column:goods"`
	Source    string    `json:"source"`
	StayTime  int       `json:"stay_time"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (BrowsingHistory) TableName() string {
	return "browsing_history"
}

// UserSetting 用户偏好设置
type UserSetting struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id" gorm:"column:user"`
	NotifyNewOrder bool      `json:"notify_new_order"`
	NotifyPromotion bool     `json:"notify_promotion"`
	NotifySystem   bool      `json:"notify_system"`
	PrivacyShowFav bool      `json:"privacy_show_fav"`
	ThemeColor     string    `json:"theme_color"`
	Language       string    `json:"language"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 指定表名
func (UserSetting) TableName() string {
	return "user_setting"
}
