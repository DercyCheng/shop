package entity

import (
	"time"
)

// ReservationStatus represents the status of a stock reservation
type ReservationStatus string

const (
	// ReservationPending indicates the reservation is active but not finalized
	ReservationPending ReservationStatus = "PENDING"
	// ReservationCommitted indicates the reservation has been converted to an actual sale
	ReservationCommitted ReservationStatus = "COMMITTED"
	// ReservationCancelled indicates the reservation was cancelled
	ReservationCancelled ReservationStatus = "CANCELLED"
	// ReservationExpired indicates the reservation expired without being committed
	ReservationExpired ReservationStatus = "EXPIRED"
)

// Reservation represents a temporary hold on inventory items
type Reservation struct {
	ID        int64             `gorm:"primaryKey"`
	OrderID   string            `gorm:"size:50;not null;index"`
	Status    ReservationStatus `gorm:"size:20;not null;default:'PENDING'"`
	Items     []ReservationItem `gorm:"foreignKey:ReservationID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time `gorm:"index"`
}

// IsActive checks if the reservation is still active (pending and not expired)
func (r *Reservation) IsActive() bool {
	return r.Status == ReservationPending && time.Now().Before(r.ExpiresAt)
}

// IsExpired checks if the reservation has expired
func (r *Reservation) IsExpired() bool {
	return r.Status == ReservationPending && time.Now().After(r.ExpiresAt)
}

// CanCommit checks if the reservation can be committed
func (r *Reservation) CanCommit() bool {
	return r.Status == ReservationPending && !r.IsExpired()
}

// CanCancel checks if the reservation can be cancelled
func (r *Reservation) CanCancel() bool {
	return r.Status == ReservationPending
}

// ReservationItem represents a single item within a reservation
type ReservationItem struct {
	ID            int64 `gorm:"primaryKey"`
	ReservationID int64 `gorm:"index;not null"`
	ProductID     int64 `gorm:"not null"`
	WarehouseID   int64 `gorm:"not null"`
	Quantity      int   `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
