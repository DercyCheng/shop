package entity

import (
	"time"
)

// Category represents a product category in the system
type Category struct {
	ID               int64      `gorm:"primaryKey;column:id"`
	Name             string     `gorm:"column:name;type:varchar(50);not null"`
	ParentCategoryID int64      `gorm:"column:parent_category_id;default:0;index:idx_parent_id"`
	Level            int        `gorm:"column:level;default:1"`
	IsTab            bool       `gorm:"column:is_tab;default:0"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at;index"`

	// Relations
	SubCategories  []*Category `gorm:"-"`
	ParentCategory *Category   `gorm:"-"`
}

// TableName returns the table name for the Category entity
func (Category) TableName() string {
	return "category"
}
