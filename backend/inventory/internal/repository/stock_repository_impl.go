package repository

import (
	"context"

	"shop/inventory/internal/domain/entity"
	"shop/inventory/internal/repository/dao"
)

// StockRepositoryImpl implements the StockRepository interface
type StockRepositoryImpl struct {
	stockDAO *dao.StockDAO
}

// NewStockRepository creates a new StockRepositoryImpl
func NewStockRepository(stockDAO *dao.StockDAO) StockRepository {
	return &StockRepositoryImpl{stockDAO: stockDAO}
}

// Create creates a new stock record
func (r *StockRepositoryImpl) Create(ctx context.Context, stock *entity.Stock) error {
	return r.stockDAO.Create(ctx, stock)
}

// GetByID retrieves a stock record by its ID
func (r *StockRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Stock, error) {
	return r.stockDAO.GetByID(ctx, id)
}

// GetByProductAndWarehouse retrieves a stock record by product ID and warehouse ID
func (r *StockRepositoryImpl) GetByProductAndWarehouse(ctx context.Context, productID, warehouseID int64) (*entity.Stock, error) {
	return r.stockDAO.GetByProductAndWarehouse(ctx, productID, warehouseID)
}

// List retrieves stocks with pagination
func (r *StockRepositoryImpl) List(ctx context.Context, page, pageSize int, warehouseID int64, lowStockOnly, outOfStockOnly bool) ([]*entity.Stock, int64, error) {
	return r.stockDAO.List(ctx, page, pageSize, warehouseID, lowStockOnly, outOfStockOnly)
}

// Update updates a stock record
func (r *StockRepositoryImpl) Update(ctx context.Context, stock *entity.Stock) error {
	return r.stockDAO.Update(ctx, stock)
}

// Delete deletes a stock record
func (r *StockRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.stockDAO.Delete(ctx, id)
}

// GetByProductIDs retrieves stock records for multiple products in a specific warehouse
func (r *StockRepositoryImpl) GetByProductIDs(ctx context.Context, productIDs []int64, warehouseID int64) ([]*entity.Stock, error) {
	return r.stockDAO.GetByProductIDs(ctx, productIDs, warehouseID)
}

// IncrementStock increases the stock quantity
func (r *StockRepositoryImpl) IncrementStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	return r.stockDAO.IncrementStock(ctx, productID, warehouseID, quantity)
}

// DecrementStock decreases the stock quantity
func (r *StockRepositoryImpl) DecrementStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	return r.stockDAO.DecrementStock(ctx, productID, warehouseID, quantity)
}

// ReserveStock places a reservation on stock
func (r *StockRepositoryImpl) ReserveStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	return r.stockDAO.ReserveStock(ctx, productID, warehouseID, quantity)
}

// CancelReservation removes a reservation from stock
func (r *StockRepositoryImpl) CancelReservation(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	return r.stockDAO.CancelReservation(ctx, productID, warehouseID, quantity)
}

// CommitReservation converts a reservation to an actual stock reduction
func (r *StockRepositoryImpl) CommitReservation(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	return r.stockDAO.CommitReservation(ctx, productID, warehouseID, quantity)
}
