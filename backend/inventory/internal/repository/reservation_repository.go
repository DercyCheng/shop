package repository

import (
	"context"

	"shop/inventory/internal/domain/entity"
)

// ReservationRepository defines the interface for reservation data access
type ReservationRepository interface {
	// Create creates a new reservation
	Create(ctx context.Context, reservation *entity.Reservation) error

	// GetByID retrieves a reservation by its ID
	GetByID(ctx context.Context, id int64) (*entity.Reservation, error)

	// GetByOrderID retrieves a reservation by order ID
	GetByOrderID(ctx context.Context, orderID string) (*entity.Reservation, error)

	// Update updates a reservation
	Update(ctx context.Context, reservation *entity.Reservation) error

	// Delete deletes a reservation
	Delete(ctx context.Context, id int64) error

	// GetExpiredReservations retrieves all expired reservations
	GetExpiredReservations(ctx context.Context) ([]*entity.Reservation, error)

	// CommitReservation marks a reservation as committed
	CommitReservation(ctx context.Context, id int64) error

	// CancelReservation marks a reservation as cancelled
	CancelReservation(ctx context.Context, id int64) error
}
