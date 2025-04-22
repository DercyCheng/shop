package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32          `gorm:"primary_key;comment:ID" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time;comment:创建时间" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time;comment:更新时间" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"comment:删除时间" json:"-"`
	IsDeleted bool           `gorm:"comment:是否删除" json:"-"`
}
