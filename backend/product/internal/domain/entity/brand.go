package entity

import (
	"time"
)

// Brand represents a product brand in the system
type Brand struct {
	ID        int64      `gorm:"primaryKey;column:id"`
	Name      string     `gorm:"column:name;type:varchar(50);not null"`
	Logo      string     `gorm:"column:logo;type:varchar(255)"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

// TableName returns the table name for the Brand entity
func (Brand) TableName() string {
	return "brands"
}
