package entity

import (
	"encoding/json"
	"time"
)

// StockStatus 库存锁定状态
type StockStatus int

const (
	StockLocked    StockStatus = 1 // 已锁定
	StockReduced   StockStatus = 2 // 已扣减
	StockReturned  StockStatus = 3 // 已归还
)

// StockSellDetail 库存扣减明细
type StockSellDetail struct {
	ID          int64      `gorm:"primaryKey"`
	OrderSN     string     `gorm:"column:order_sn;type:varchar(50);uniqueIndex;not null;comment:'订单号'"`
	Status      StockStatus `gorm:"type:int;default:1;index;not null;comment:'状态：1:锁定，2:已扣减，3:已归还'"`
	Detail      string     `gorm:"type:json;comment:'库存扣减明细，结构为[{goods_id:1, num:2, warehouse_id:1}]'"`
	LockTime    *time.Time `gorm:"type:datetime(3);comment:'锁定时间'"`
	ConfirmTime *time.Time `gorm:"type:datetime(3);comment:'确认时间'"`
	CreatedAt   time.Time  `gorm:"type:datetime(3)"`
	UpdatedAt   time.Time  `gorm:"type:datetime(3)"`
	DeletedAt   *time.Time `gorm:"type:datetime(3)"`
	
	// 非数据库字段，用于Detail的JSON转换
	DetailItems []*StockDetail `gorm:"-"`
}

// StockDetail 库存操作详情项
type StockDetail struct {
	ProductID   int64 `json:"goods_id"`
	Quantity    int   `json:"num"`
	WarehouseID int   `json:"warehouse_id"`
}

// TableName 指定表名
func (StockSellDetail) TableName() string {
	return "stock_sell_detail"
}

// BeforeSave 保存前的钩子函数，将DetailItems转换为JSON字符串
func (s *StockSellDetail) BeforeSave() error {
	if len(s.DetailItems) > 0 {
		data, err := json.Marshal(s.DetailItems)
		if err != nil {
			return err
		}
		s.Detail = string(data)
	}
	return nil
}

// AfterFind 查询后的钩子函数，将JSON字符串解析为DetailItems
func (s *StockSellDetail) AfterFind() error {
	if s.Detail != "" {
		return json.Unmarshal([]byte(s.Detail), &s.DetailItems)
	}
	return nil
}
