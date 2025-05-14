package repository

import (
	"context"
	
	"shop/backend/inventory/internal/domain/entity"
	"shop/backend/inventory/internal/domain/valueobject"
)

// InventoryRepository 库存仓储接口
type InventoryRepository interface {
	// 库存基本操作
	GetInventory(ctx context.Context, productID int64, warehouseID int) (*entity.Inventory, error)
	BatchGetInventory(ctx context.Context, productIDs []int64, warehouseID int) ([]*entity.Inventory, error)
	SetInventory(ctx context.Context, inventory *entity.Inventory) error
	
	// 库存锁定和扣减
	LockStock(ctx context.Context, orderSN string, items []*valueobject.StockOperation) (*valueobject.LockResult, error)
	UnlockStock(ctx context.Context, orderSN string) error
	ReduceStock(ctx context.Context, orderSN string) error
	
	// 库存调整
	IncreaseStock(ctx context.Context, productID int64, warehouseID int, quantity int, remark string) error
	DecreaseStock(ctx context.Context, productID int64, warehouseID int, quantity int, remark string) error
	AdjustStock(ctx context.Context, productID int64, warehouseID int, newStock int, operator string, remark string) error
	
	// 库存锁定记录操作
	GetStockSellDetail(ctx context.Context, orderSN string) (*entity.StockSellDetail, error)
	UpdateStockSellDetailStatus(ctx context.Context, orderSN string, status entity.StockStatus) error
	CreateStockSellDetail(ctx context.Context, detail *entity.StockSellDetail) error
	
	// 历史记录
	RecordInventoryHistory(ctx context.Context, history *entity.InventoryHistory) error
	GetInventoryHistoryByProductID(ctx context.Context, productID int64, page, pageSize int) ([]*entity.InventoryHistory, int64, error)
}

// WarehouseRepository 仓库仓储接口
type WarehouseRepository interface {
	// 仓库基本操作
	GetWarehouse(ctx context.Context, id int) (*entity.Warehouse, error)
	GetWarehouseByName(ctx context.Context, name string) (*entity.Warehouse, error)
	ListWarehouses(ctx context.Context, page, pageSize int) ([]*entity.Warehouse, int64, error)
	CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error
	DeleteWarehouse(ctx context.Context, id int) error
	
	// 更多仓库相关操作...
}
