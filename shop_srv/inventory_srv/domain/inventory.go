package model

import (
	"database/sql/driver"
	"encoding/json"
)

//type Stock struct{
//	BaseModel
//	Name string
//	Address string
//}
type GoodsDetail struct {
	Goods int32
	Num   int32
}
type GormList []GoodsDetail

// 实现 driver.Valuer 接口，Value 返回 json value
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index;comment:商品id"`
	Stocks  int32 `gorm:"type:int;comment:仓库"`
	Version int32 `gorm:"type:int;comment:分布式锁-乐观锁"`
}

//type Delivery struct {
//	Goods   int32 `gorm:"type:int;index"`
//	Nums    int32 `gorm:"type:int"`
//	OrderSn string `gorm:"type:varchar(200)"`
//	Status  string `gorm:"type:varchar(200)"` //1.已扣减 2.已归还
//}
type StockSellDetail struct {
	OrderSn string   `gorm:"type:varchar(200);index:idx_order_sn,unique;comment:订单编号"`
	Status  int32    `gorm:"type:varchar(200);comment:1.已扣减,2.已归还"`
	Detail  GormList `gorm:"type:varchar(200);comment:详细商品"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}

//type InventoryHistory struct {
//	user   int32 `gorm:"type:int;comment:用户"`
//	goods  int32 `gorm:"type:int;comment:商品"`
//	nums   int32 `gorm:"type:int;comment:数量"`
//	order  int32 `gorm:"type:int;comment:订单编号"`
//	status int32 `gorm:"type:int;comment:状态：1.预扣减2.已经支付,幂等性"`
//}
