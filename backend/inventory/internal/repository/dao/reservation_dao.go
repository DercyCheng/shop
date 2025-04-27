package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"shop/inventory/internal/domain/entity"
)

// ReservationDAO provides data access methods for reservation entity
type ReservationDAO struct {
	db *gorm.DB
}

// NewReservationDAO creates a new ReservationDAO
func NewReservationDAO(db *gorm.DB) *ReservationDAO {
	return &ReservationDAO{db: db}
}

// Create inserts a new reservation
func (d *ReservationDAO) Create(ctx context.Context, reservation *entity.Reservation) error {
	return d.db.WithContext(ctx).Create(reservation).Error
}

// GetByID retrieves a reservation by ID with its items
func (d *ReservationDAO) GetByID(ctx context.Context, id int64) (*entity.Reservation, error) {
	var reservation entity.Reservation
	if err := d.db.WithContext(ctx).Preload("Items").First(&reservation, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}
	return &reservation, nil
}

// GetByOrderID retrieves a reservation by order ID with its items
func (d *ReservationDAO) GetByOrderID(ctx context.Context, orderID string) (*entity.Reservation, error) {
	var reservation entity.Reservation
	if err := d.db.WithContext(ctx).Preload("Items").Where("order_id = ?", orderID).First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}
	return &reservation, nil
}

// Update updates a reservation
func (d *ReservationDAO) Update(ctx context.Context, reservation *entity.Reservation) error {
	return d.db.WithContext(ctx).Save(reservation).Error
}

// Delete deletes a reservation
func (d *ReservationDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Delete(&entity.Reservation{}, id).Error
}

// GetExpiredReservations retrieves all pending reservations that have expired
func (d *ReservationDAO) GetExpiredReservations(ctx context.Context) ([]*entity.Reservation, error) {
	var reservations []*entity.Reservation

	if err := d.db.WithContext(ctx).
		Preload("Items").
		Where("status = ? AND expires_at < ?", entity.ReservationPending, time.Now()).
		Find(&reservations).Error; err != nil {
		return nil, err
	}

	return reservations, nil
}

// CommitReservation marks a reservation as committed with concurrency control
func (d *ReservationDAO) CommitReservation(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var reservation entity.Reservation

		// Lock the record for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Items").
			First(&reservation, id).Error; err != nil {
			return err
		}

		// Check if the reservation can be committed
		if !reservation.CanCommit() {
			return errors.New("reservation cannot be committed: it is either expired or not in pending state")
		}

		// Update the status
		reservation.Status = entity.ReservationCommitted

		// Save the updated reservation
		return tx.Save(&reservation).Error
	})
}

// CancelReservation marks a reservation as cancelled with concurrency control
func (d *ReservationDAO) CancelReservation(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var reservation entity.Reservation

		// Lock the record for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Items").
			First(&reservation, id).Error; err != nil {
			return err
		}

		// Check if the reservation can be cancelled
		if !reservation.CanCancel() {
			return errors.New("reservation cannot be cancelled: it is not in pending state")
		}

		// Update the status
		reservation.Status = entity.ReservationCancelled

		// Save the updated reservation
		return tx.Save(&reservation).Error
	})
}
