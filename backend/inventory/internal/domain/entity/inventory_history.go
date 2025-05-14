package entity

import (
	"time"
)

// OperationType 库存操作类型
type OperationType string

const (
	OperationLock     OperationType = "lock"     // 锁定库存
	OperationUnlock   OperationType = "unlock"   // 解锁库存
	OperationDecrease OperationType = "decrease" // 减少库存
	OperationIncrease OperationType = "increase" // 增加库存
	OperationAdjust   OperationType = "adjust"   // 调整库存（盘点）
)

// InventoryHistory 库存变更历史
type InventoryHistory struct {
	ID          int64        `gorm:"primaryKey"`
	ProductID   int64        `gorm:"column:goods;index;not null;comment:'商品ID'"`
	WarehouseID int          `gorm:"not null;comment:'仓库ID'"`
	Quantity    int          `gorm:"not null;comment:'变更数量（正数增加，负数减少）'"`
	Operation   OperationType `gorm:"column:operation_type;type:varchar(20);not null;comment:'操作类型：lock, unlock, decrease, increase, adjust'"`
	Operator    string       `gorm:"type:varchar(50);comment:'操作人'"`
	OrderSN     string       `gorm:"column:order_sn;type:varchar(50);index;comment:'相关订单号'"`
	Remark      string       `gorm:"type:varchar(255);comment:'备注'"`
	CreatedAt   time.Time    `gorm:"type:datetime(3)"`
}

// TableName 指定表名
func (InventoryHistory) TableName() string {
	return "inventory_history"
}

// InventoryChangeRecord MongoDB版本的库存变更记录
type InventoryChangeRecord struct {
	ProductID    int64        `bson:"product_id"`
	OrderSn      string       `bson:"order_sn"`
	Operation    string       `bson:"operation"`
	Quantity     int32        `bson:"quantity"`
	BeforeStock  int32        `bson:"before_stock"`
	AfterStock   int32        `bson:"after_stock"`
	Operator     string       `bson:"operator"`
	OperatorID   int64        `bson:"operator_id"`
	Reason       string       `bson:"reason"`
	Timestamp    time.Time    `bson:"timestamp"`
	Success      bool         `bson:"success"`
	ErrorMessage string       `bson:"error_message,omitempty"`
}

// InventoryActivitySummary 库存操作统计摘要
type InventoryActivitySummary struct {
	ProductID     int64  `bson:"product_id"`
	Operation     string `bson:"operation"`
	Count         int    `bson:"count"`
	TotalQuantity int    `bson:"total_quantity"`
	MaxQuantity   int    `bson:"max_quantity"`
	MinQuantity   int    `bson:"min_quantity"`
	AvgQuantity   float64 `bson:"avg_quantity"`
}
