package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32          `gorm:"primary_key;comment:ID"`
	CreatedAt time.Time      `gorm:"column:add_time;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:update_time;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"comment:删除时间"`
	IsDeleted bool           `gorm:"comment:是否删除"`
}
