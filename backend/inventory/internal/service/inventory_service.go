package service

import (
	"context"
	"time"
	
	"shop/backend/inventory/internal/domain/entity"
	"shop/backend/inventory/internal/domain/valueobject"
)

// InventoryService 库存服务接口
type InventoryService interface {
	// 基本库存查询
	GetInventory(ctx context.Context, productID int64) (*entity.Inventory, error)
	BatchGetInventory(ctx context.Context, productIDs []int64) ([]*entity.Inventory, error)
	GetAvailableStock(ctx context.Context, productID int64) (int, error)
	CheckStockAvailable(ctx context.Context, productID int64, quantity int) (bool, error)
	
	// 库存操作
	SetInventory(ctx context.Context, productID int64, stock int, operator string) error
	AddStock(ctx context.Context, productID int64, quantity int, remark string) error
	AdjustStock(ctx context.Context, productID int64, newStock int, operator string, remark string) error
	
	// 库存历史
	GetInventoryHistory(ctx context.Context, productID int64, page, pageSize int) ([]*entity.InventoryHistory, int64, error)
}

// InventoryLockService 库存锁定服务接口
type InventoryLockService interface {
	// 库存锁定相关
	LockInventory(ctx context.Context, lockKey string, items []LockItem, timeoutSeconds int) (*LockResult, error)
	UnlockInventory(ctx context.Context, lockKey string) error
	ConfirmReduce(ctx context.Context, lockKey string) error
	GetLockDetail(ctx context.Context, lockKey string) (*entity.StockSellDetail, error)
	GetExpiredLocks(ctx context.Context, before time.Time) ([]string, error)
}

// LockItem 锁定项目
type LockItem struct {
	ProductID int64
	Quantity  int
}

// LockFailItem 锁定失败项
type LockFailItem struct {
	ProductID int64
	Quantity  int
	Available int
	Reason    string
}

// LockResult 锁定结果
type LockResult struct {
	Success   bool
	Message   string
	FailItems []*LockFailItem
}

// WarehouseService 仓库服务接口
type WarehouseService interface {
	// 仓库管理
	GetWarehouse(ctx context.Context, id int) (*entity.Warehouse, error)
	GetWarehouseByName(ctx context.Context, name string) (*entity.Warehouse, error)
	ListWarehouses(ctx context.Context, page, pageSize int) ([]*entity.Warehouse, int64, error)
	CreateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error
	UpdateWarehouse(ctx context.Context, warehouse *entity.Warehouse) error
	DeleteWarehouse(ctx context.Context, id int) error
}
