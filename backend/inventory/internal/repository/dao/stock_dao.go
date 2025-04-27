package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"shop/inventory/internal/domain/entity"
)

// StockDAO provides data access methods for stock entity
type StockDAO struct {
	db *gorm.DB
}

// NewStockDAO creates a new StockDAO
func NewStockDAO(db *gorm.DB) *StockDAO {
	return &StockDAO{db: db}
}

// Create inserts a new stock record
func (d *StockDAO) Create(ctx context.Context, stock *entity.Stock) error {
	return d.db.WithContext(ctx).Create(stock).Error
}

// GetByID retrieves a stock by ID
func (d *StockDAO) GetByID(ctx context.Context, id int64) (*entity.Stock, error) {
	var stock entity.Stock
	if err := d.db.WithContext(ctx).First(&stock, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stock not found")
		}
		return nil, err
	}
	return &stock, nil
}

// GetByProductAndWarehouse retrieves a stock by product ID and warehouse ID
func (d *StockDAO) GetByProductAndWarehouse(ctx context.Context, productID, warehouseID int64) (*entity.Stock, error) {
	var stock entity.Stock
	if err := d.db.WithContext(ctx).Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stock not found")
		}
		return nil, err
	}
	return &stock, nil
}

// List retrieves stocks with pagination
func (d *StockDAO) List(ctx context.Context, page, pageSize int, warehouseID int64, lowStockOnly, outOfStockOnly bool) ([]*entity.Stock, int64, error) {
	var stocks []*entity.Stock
	var total int64

	query := d.db.WithContext(ctx)

	if warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	// Count total before applying pagination
	countQuery := query
	if err := countQuery.Model(&entity.Stock{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply low stock filter if requested
	if lowStockOnly {
		query = query.Where("quantity - reserved <= low_stock_threshold AND quantity - reserved > 0")
	}

	// Apply out of stock filter if requested
	if outOfStockOnly {
		query = query.Where("quantity - reserved <= 0")
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&stocks).Error; err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

// Update updates a stock record
func (d *StockDAO) Update(ctx context.Context, stock *entity.Stock) error {
	return d.db.WithContext(ctx).Save(stock).Error
}

// Delete deletes a stock record
func (d *StockDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Delete(&entity.Stock{}, id).Error
}

// GetByProductIDs retrieves stock records for multiple products in a specific warehouse
func (d *StockDAO) GetByProductIDs(ctx context.Context, productIDs []int64, warehouseID int64) ([]*entity.Stock, error) {
	var stocks []*entity.Stock

	query := d.db.WithContext(ctx).Where("product_id IN ?", productIDs)

	if warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	if err := query.Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}

// IncrementStock increases the stock quantity with concurrency control
func (d *StockDAO) IncrementStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	var stock entity.Stock

	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the record for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
			First(&stock).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create a new stock record if it doesn't exist
				stock = entity.Stock{
					ProductID:   productID,
					WarehouseID: warehouseID,
					Quantity:    quantity,
				}
				return tx.Create(&stock).Error
			}
			return err
		}

		// Increment the quantity
		stock.Quantity += quantity

		// Save the updated stock
		return tx.Save(&stock).Error
	})

	if err != nil {
		return nil, err
	}

	return &stock, nil
}

// DecrementStock decreases the stock quantity with concurrency control
func (d *StockDAO) DecrementStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	var stock entity.Stock

	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the record for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
			First(&stock).Error; err != nil {
			return err
		}

		// Check if there's enough stock
		if stock.Available() < quantity {
			return errors.New("insufficient stock")
		}

		// Decrement the quantity
		if !stock.Decrement(quantity) {
			return errors.New("failed to decrement stock")
		}

		// Save the updated stock
		return tx.Save(&stock).Error
	})

	if err != nil {
		return nil, err
	}

	return &stock, nil
}

// ReserveStock places a reservation on stock
func (d *StockDAO) ReserveStock(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	var stock entity.Stock

	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the record for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
			First(&stock).Error; err != nil {
			return err
		}

		// Check if reservation is possible
		if !stock.CanReserve(quantity) {
			return errors.New("insufficient stock for reservation")
		}

		// Reserve the quantity
		if !stock.Reserve(quantity) {
			return errors.New("failed to reserve stock")
		}

		// Save the updated stock
		return tx.Save(&stock).Error
	})

	if err != nil {
		return nil, err
	}

	return &stock, nil
}

// CancelReservation removes a reservation from stock
func (d *StockDAO) CancelReservation(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	var stock entity.Stock

	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the record for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
			First(&stock).Error; err != nil {
			return err
		}

		// Cancel the reservation
		if !stock.CancelReservation(quantity) {
			return errors.New("failed to cancel reservation: reserved quantity is less than requested")
		}

		// Save the updated stock
		return tx.Save(&stock).Error
	})

	if err != nil {
		return nil, err
	}

	return &stock, nil
}

// CommitReservation converts a reservation to an actual stock reduction
func (d *StockDAO) CommitReservation(ctx context.Context, productID, warehouseID int64, quantity int) (*entity.Stock, error) {
	var stock entity.Stock

	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the record for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
			First(&stock).Error; err != nil {
			return err
		}

		// Commit the reservation
		if !stock.CommitReservation(quantity) {
			return errors.New("failed to commit reservation: reserved quantity is less than requested")
		}

		// Save the updated stock
		return tx.Save(&stock).Error
	})

	if err != nil {
		return nil, err
	}

	return &stock, nil
}
