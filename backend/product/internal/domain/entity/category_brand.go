package entity

import (
	"time"
)

// CategoryBrand represents a relationship between a category and a brand
type CategoryBrand struct {
	ID         int64      `gorm:"primaryKey;column:id"`
	CategoryID int64      `gorm:"column:category_id;not null;uniqueIndex:idx_category_brand,priority:1"`
	BrandID    int64      `gorm:"column:brands_id;not null;uniqueIndex:idx_category_brand,priority:2"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;index"`

	// Relations
	Category *Category `gorm:"-"`
	Brand    *Brand    `gorm:"-"`
}

// TableName returns the table name for the CategoryBrand entity
func (CategoryBrand) TableName() string {
	return "goods_category_brand"
}
