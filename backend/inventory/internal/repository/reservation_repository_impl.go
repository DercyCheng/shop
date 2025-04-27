package repository

import (
	"context"

	"shop/inventory/internal/domain/entity"
	"shop/inventory/internal/repository/dao"
)

// ReservationRepositoryImpl implements the ReservationRepository interface
type ReservationRepositoryImpl struct {
	reservationDAO *dao.ReservationDAO
}

// NewReservationRepository creates a new ReservationRepositoryImpl
func NewReservationRepository(reservationDAO *dao.ReservationDAO) ReservationRepository {
	return &ReservationRepositoryImpl{reservationDAO: reservationDAO}
}

// Create creates a new reservation
func (r *ReservationRepositoryImpl) Create(ctx context.Context, reservation *entity.Reservation) error {
	return r.reservationDAO.Create(ctx, reservation)
}

// GetByID retrieves a reservation by its ID
func (r *ReservationRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Reservation, error) {
	return r.reservationDAO.GetByID(ctx, id)
}

// GetByOrderID retrieves a reservation by order ID
func (r *ReservationRepositoryImpl) GetByOrderID(ctx context.Context, orderID string) (*entity.Reservation, error) {
	return r.reservationDAO.GetByOrderID(ctx, orderID)
}

// Update updates a reservation
func (r *ReservationRepositoryImpl) Update(ctx context.Context, reservation *entity.Reservation) error {
	return r.reservationDAO.Update(ctx, reservation)
}

// Delete deletes a reservation
func (r *ReservationRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.reservationDAO.Delete(ctx, id)
}

// GetExpiredReservations retrieves all expired reservations
func (r *ReservationRepositoryImpl) GetExpiredReservations(ctx context.Context) ([]*entity.Reservation, error) {
	return r.reservationDAO.GetExpiredReservations(ctx)
}

// CommitReservation marks a reservation as committed
func (r *ReservationRepositoryImpl) CommitReservation(ctx context.Context, id int64) error {
	return r.reservationDAO.CommitReservation(ctx, id)
}

// CancelReservation marks a reservation as cancelled
func (r *ReservationRepositoryImpl) CancelReservation(ctx context.Context, id int64) error {
	return r.reservationDAO.CancelReservation(ctx, id)
}
