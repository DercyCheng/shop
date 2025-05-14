package entity

import (
	"time"
)

// Banner 轮播图实体
type Banner struct {
	ID        int64     `json:"id"`
	Image     string    `json:"image"`
	URL       string    `json:"url"`
	Index     int       `json:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
