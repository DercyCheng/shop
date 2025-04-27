package entity

import (
	"time"
)

// Warehouse represents a physical location where inventory is stored
type Warehouse struct {
	ID        int64  `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Address   string `gorm:"size:255"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
