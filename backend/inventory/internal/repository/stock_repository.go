package repository

import (
	"context"

	"shop/inventory/internal/domain/entity"
)

// StockRepository defines the interface for stock data access
type StockRepository interface {
	// Create creates a new stock record
	Create(ctx context.Context, stock *entity.Stock) error

	// GetByID retrieves a stock record by its ID
	GetByID(ctx context.Context, id int64) (*entity.Stock, error)

	// GetByProductAndWarehouse retrieves a stock record by product ID and warehouse ID
	GetByProductAndWarehouse(ctx context.Context, productID, warehouseID int64) (*entity.Stock, error)

	// List retrieves stocks with pagination
	List(ctx context.Context, page, pageSize int, warehouseID int64, lowStockOnly, outOfStockOnly bool) ([]*entity.Stock, int64, error)

	// Update updates a stock record
	Update(ctx context.Context, stock *entity.Stock) error

	// Delete deletes a stock record
	Delete(ctx context.Context, id int64) error

	// GetByProductIDs retrieves stock records for multiple products in a specific warehouse
	GetByProductIDs(ctx context.Context, productIDs []int64, warehouseID int64) ([]*entity.Stock, error)

	// IncrementStock increases the stock quantity
	IncrementStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error)

	// DecrementStock decreases the stock quantity
	DecrementStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error)

	// ReserveStock places a reservation on stock
	ReserveStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error)

	// CancelReservation removes a reservation from stock
	CancelReservation(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error)

	// CommitReservation converts a reservation to an actual stock reduction
	CommitReservation(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error)
}
