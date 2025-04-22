package model

import (
	"database/sql/driver"
	"encoding/json"
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

type GormList []string

// 实现 driver.Valuer 接口，Value 返回 json value
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}
