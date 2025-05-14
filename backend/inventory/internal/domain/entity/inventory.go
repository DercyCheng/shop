package entity

import (
	"time"
)

// Inventory 库存实体
type Inventory struct {
	ID              int64     `gorm:"primaryKey"`
	ProductID       int64     `gorm:"column:goods;index:idx_goods_warehouse,unique;not null;comment:'商品ID'"`
	Stock           int       `gorm:"column:stocks;not null;default:0;comment:'库存数量'"`
	Version         int       `gorm:"not null;default:0;comment:'乐观锁版本号'"`
	WarehouseID     int       `gorm:"not null;default:1;index:idx_goods_warehouse,unique;comment:'仓库ID'"`
	LockStock       int       `gorm:"column:lock_stocks;not null;default:0;comment:'锁定库存数量'"`
	AlertThreshold  int       `gorm:"default:10;comment:'预警阈值'"`
	CreatedAt       time.Time `gorm:"type:datetime(3)"`
	UpdatedAt       time.Time `gorm:"type:datetime(3)"`
	DeletedAt       *time.Time `gorm:"type:datetime(3)"`
}

// TableName 指定表名
func (Inventory) TableName() string {
	return "inventory"
}

// AvailableStock 获取可用库存数量
func (i *Inventory) AvailableStock() int {
	available := i.Stock - i.LockStock
	if available < 0 {
		return 0
	}
	return available
}

// IsAvailable 判断库存是否充足
func (i *Inventory) IsAvailable(quantity int) bool {
	return i.AvailableStock() >= quantity
}

// NeedsAlert 判断是否需要库存预警
func (i *Inventory) NeedsAlert() bool {
	return i.Stock <= i.AlertThreshold
}
