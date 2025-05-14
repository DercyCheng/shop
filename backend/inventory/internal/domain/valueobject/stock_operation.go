package valueobject

// StockOperation 库存操作值对象
type StockOperation struct {
	ProductID   int64
	WarehouseID int
	Quantity    int
	OrderSN     string
	Operator    string
	Remark      string
}

// LockResult 库存锁定结果
type LockResult struct {
	Success   bool
	Message   string
	FailItems []*LockFailItem
}

// LockFailItem 锁定失败项
type LockFailItem struct {
	ProductID int64
	Quantity  int
	Available int
	Reason    string
}
