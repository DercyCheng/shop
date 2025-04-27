package entity

import (
	"time"
)

// Stock represents inventory level for a product in a warehouse
type Stock struct {
	ID                int64 `gorm:"primaryKey"`
	ProductID         int64 `gorm:"index:idx_product_warehouse,unique;not null"`
	WarehouseID       int64 `gorm:"index:idx_product_warehouse,unique;not null"`
	Quantity          int   `gorm:"not null;default:0"`
	Reserved          int   `gorm:"not null;default:0"`
	LowStockThreshold int   `gorm:"default:5"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Available returns the quantity available for purchase (total minus reserved)
func (s *Stock) Available() int {
	available := s.Quantity - s.Reserved
	if available < 0 {
		return 0
	}
	return available
}

// IsInStock checks if the product is in stock
func (s *Stock) IsInStock() bool {
	return s.Available() > 0
}

// IsLowStock checks if the stock is below the low stock threshold
func (s *Stock) IsLowStock() bool {
	return s.Available() <= s.LowStockThreshold && s.Available() > 0
}

// CanReserve checks if the requested quantity can be reserved
func (s *Stock) CanReserve(quantity int) bool {
	return s.Available() >= quantity
}

// Reserve reduces the available quantity by reserving the specified amount
func (s *Stock) Reserve(quantity int) bool {
	if !s.CanReserve(quantity) {
		return false
	}
	s.Reserved += quantity
	return true
}

// CancelReservation returns the reserved quantity back to available stock
func (s *Stock) CancelReservation(quantity int) bool {
	if s.Reserved < quantity {
		return false
	}
	s.Reserved -= quantity
	return true
}

// CommitReservation converts a reservation into an actual reduction of stock
func (s *Stock) CommitReservation(quantity int) bool {
	if s.Reserved < quantity {
		return false
	}
	s.Quantity -= quantity
	s.Reserved -= quantity
	return true
}

// Increment increases the quantity
func (s *Stock) Increment(quantity int) {
	s.Quantity += quantity
}

// Decrement decreases the quantity if possible
func (s *Stock) Decrement(quantity int) bool {
	if s.Available() < quantity {
		return false
	}
	s.Quantity -= quantity
	return true
}
