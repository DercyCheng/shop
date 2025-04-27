package service

import (
	"context"

	"shop/inventory/internal/domain/entity"
)

// ReservationService defines methods for reservation-related operations
type ReservationService interface {
	// CreateReservation creates a new stock reservation
	CreateReservation(ctx context.Context, orderID string, items []*entity.ReservationItem, expirationMinutes int) (*entity.Reservation, error)

	// GetReservationByID retrieves a reservation by ID
	GetReservationByID(ctx context.Context, id int64) (*entity.Reservation, error)

	// GetReservationByOrderID retrieves a reservation by order ID
	GetReservationByOrderID(ctx context.Context, orderID string) (*entity.Reservation, error)

	// CommitReservation confirms a reservation (e.g., after payment)
	CommitReservation(ctx context.Context, id int64) error

	// CancelReservation cancels a reservation and returns stock
	CancelReservation(ctx context.Context, id int64) error

	// ProcessExpiredReservations finds and cancels all expired reservations
	ProcessExpiredReservations(ctx context.Context) (int, error)

	// ExtendReservation extends the expiration time of a reservation
	ExtendReservation(ctx context.Context, id int64, additionalMinutes int) error
}
