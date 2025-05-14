package entity

import (
	"time"
)

// Warehouse 仓库实体
type Warehouse struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null;comment:'仓库名称'"`
	Address   string    `gorm:"type:varchar(255);not null;comment:'仓库地址'"`
	Contact   string    `gorm:"type:varchar(50);comment:'联系人'"`
	Phone     string    `gorm:"type:varchar(20);comment:'联系电话'"`
	Status    int8      `gorm:"type:tinyint(1);default:1;index;comment:'状态：1-正常，0-禁用'"`
	CreatedAt time.Time `gorm:"type:datetime(3)"`
	UpdatedAt time.Time `gorm:"type:datetime(3)"`
	DeletedAt *time.Time `gorm:"type:datetime(3)"`
}

// TableName 指定表名
func (Warehouse) TableName() string {
	return "warehouse"
}

// IsActive 判断仓库是否处于可用状态
func (w *Warehouse) IsActive() bool {
	return w.Status == 1
}
