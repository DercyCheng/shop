package entity

import (
	"time"
)

// Brand 品牌实体
type Brand struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Logo      string    `json:"logo"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	
	// 关联实体
	Categories []*Category `json:"categories,omitempty"`
}
