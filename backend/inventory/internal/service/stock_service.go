package service

import (
	"context"

	"shop/inventory/internal/domain/entity"
)

// StockService defines methods for stock-related operations
type StockService interface {
	// GetStockByID retrieves a stock by ID
	GetStockByID(ctx context.Context, id int64) (*entity.Stock, error)

	// GetStockByProductAndWarehouse retrieves a stock by product ID and warehouse ID
	GetStockByProductAndWarehouse(ctx context.Context, productID, warehouseID int64) (*entity.Stock, error)

	// ListStocks retrieves stocks with pagination
	ListStocks(ctx context.Context, page, pageSize int, warehouseID int64, lowStockOnly, outOfStockOnly bool) ([]*entity.Stock, int64, error)

	// CreateStock creates a new stock record
	CreateStock(ctx context.Context, stock *entity.Stock) error

	// UpdateStock updates a stock record
	UpdateStock(ctx context.Context, stock *entity.Stock) error

	// DeleteStock deletes a stock record
	DeleteStock(ctx context.Context, id int64) error

	// GetStocksByProductIDs retrieves stock records for multiple products in a warehouse
	GetStocksByProductIDs(ctx context.Context, productIDs []int64, warehouseID int64) ([]*entity.Stock, error)

	// AddStock increases the available stock quantity
	AddStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error)

	// RemoveStock decreases the available stock quantity
	RemoveStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error)

	// ReserveStock places a reservation on stock
	ReserveStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error)

	// CancelStockReservation removes a reservation from stock
	CancelStockReservation(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error)

	// CommitStockReservation converts a reservation to an actual stock reduction
	CommitStockReservation(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error)

	// CheckStockAvailability checks if a product is available in the requested quantity
	CheckStockAvailability(ctx context.Context, productID, warehouseID int64, quantity int) (bool, error)
}
