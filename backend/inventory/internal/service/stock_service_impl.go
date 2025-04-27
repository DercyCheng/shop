package service

import (
	"context"
	"errors"

	"shop/inventory/internal/domain/entity"
	"shop/inventory/internal/repository"
)

// StockServiceImpl implements the StockService interface
type StockServiceImpl struct {
	stockRepo repository.StockRepository
}

// NewStockService creates a new StockServiceImpl
func NewStockService(stockRepo repository.StockRepository) StockService {
	return &StockServiceImpl{stockRepo: stockRepo}
}

// GetStockByID retrieves a stock by ID
func (s *StockServiceImpl) GetStockByID(ctx context.Context, id int64) (*entity.Stock, error) {
	return s.stockRepo.GetByID(ctx, id)
}

// GetStockByProductAndWarehouse retrieves a stock by product ID and warehouse ID
func (s *StockServiceImpl) GetStockByProductAndWarehouse(ctx context.Context, productID, warehouseID int64) (*entity.Stock, error) {
	return s.stockRepo.GetByProductAndWarehouse(ctx, productID, warehouseID)
}

// ListStocks retrieves stocks with pagination
func (s *StockServiceImpl) ListStocks(ctx context.Context, page, pageSize int, warehouseID int64, lowStockOnly, outOfStockOnly bool) ([]*entity.Stock, int64, error) {
	// Ensure minimum page and pageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	return s.stockRepo.List(ctx, page, pageSize, warehouseID, lowStockOnly, outOfStockOnly)
}

// CreateStock creates a new stock record
func (s *StockServiceImpl) CreateStock(ctx context.Context, stock *entity.Stock) error {
	// Validate the stock data
	if stock.ProductID <= 0 {
		return errors.New("invalid product ID")
	}
	if stock.WarehouseID <= 0 {
		return errors.New("invalid warehouse ID")
	}
	if stock.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	// Check if stock already exists for this product and warehouse
	existingStock, err := s.stockRepo.GetByProductAndWarehouse(ctx, stock.ProductID, stock.WarehouseID)
	if err == nil && existingStock != nil {
		return errors.New("stock already exists for this product and warehouse")
	}

	return s.stockRepo.Create(ctx, stock)
}

// UpdateStock updates a stock record
func (s *StockServiceImpl) UpdateStock(ctx context.Context, stock *entity.Stock) error {
	// Validate the stock data
	if stock.ID <= 0 {
		return errors.New("invalid stock ID")
	}
	if stock.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	// Check if stock exists
	existingStock, err := s.stockRepo.GetByID(ctx, stock.ID)
	if err != nil {
		return err
	}
	if existingStock == nil {
		return errors.New("stock not found")
	}

	return s.stockRepo.Update(ctx, stock)
}

// DeleteStock deletes a stock record
func (s *StockServiceImpl) DeleteStock(ctx context.Context, id int64) error {
	// Check if stock exists
	existingStock, err := s.stockRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingStock == nil {
		return errors.New("stock not found")
	}

	return s.stockRepo.Delete(ctx, id)
}

// GetStocksByProductIDs retrieves stock records for multiple products in a warehouse
func (s *StockServiceImpl) GetStocksByProductIDs(ctx context.Context, productIDs []int64, warehouseID int64) ([]*entity.Stock, error) {
	if len(productIDs) == 0 {
		return []*entity.Stock{}, nil
	}
	if warehouseID <= 0 {
		return nil, errors.New("invalid warehouse ID")
	}

	return s.stockRepo.GetByProductIDs(ctx, productIDs, warehouseID)
}

// AddStock increases the available stock quantity
func (s *StockServiceImpl) AddStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}
	if warehouseID <= 0 {
		return nil, errors.New("invalid warehouse ID")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	return s.stockRepo.IncrementStock(ctx, productID, warehouseID, quantity)
}

// RemoveStock decreases the available stock quantity
func (s *StockServiceImpl) RemoveStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}
	if warehouseID <= 0 {
		return nil, errors.New("invalid warehouse ID")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	return s.stockRepo.DecrementStock(ctx, productID, warehouseID, quantity)
}

// ReserveStock places a reservation on stock
func (s *StockServiceImpl) ReserveStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}
	if warehouseID <= 0 {
		return nil, errors.New("invalid warehouse ID")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	return s.stockRepo.ReserveStock(ctx, productID, warehouseID, quantity)
}

// CancelStockReservation removes a reservation from stock
func (s *StockServiceImpl) CancelStockReservation(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}
	if warehouseID <= 0 {
		return nil, errors.New("invalid warehouse ID")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	return s.stockRepo.CancelReservation(ctx, productID, warehouseID, quantity)
}

// CommitStockReservation converts a reservation to an actual stock reduction
func (s *StockServiceImpl) CommitStockReservation(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}
	if warehouseID <= 0 {
		return nil, errors.New("invalid warehouse ID")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	return s.stockRepo.CommitReservation(ctx, productID, warehouseID, quantity)
}

// CheckStockAvailability checks if a product is available in the requested quantity
func (s *StockServiceImpl) CheckStockAvailability(ctx context.Context, productID, warehouseID int64, quantity int) (bool, error) {
	if productID <= 0 {
		return false, errors.New("invalid product ID")
	}
	if warehouseID <= 0 {
		return false, errors.New("invalid warehouse ID")
	}
	if quantity <= 0 {
		return false, errors.New("quantity must be greater than zero")
	}

	stock, err := s.stockRepo.GetByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return false, err
	}

	return stock.Available() >= quantity, nil
}
