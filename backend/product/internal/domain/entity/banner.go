package entity

import (
	"time"
)

// Banner represents a banner image in the system
type Banner struct {
	ID        int64      `gorm:"primaryKey;column:id"`
	Image     string     `gorm:"column:image;type:varchar(255);not null"`
	URL       string     `gorm:"column:url;type:varchar(255)"`
	Index     int        `gorm:"column:index;default:0"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

// TableName returns the table name for the Banner entity
func (Banner) TableName() string {
	return "banner"
}
