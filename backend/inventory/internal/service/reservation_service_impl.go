package service

import (
	"context"
	"errors"
	"time"

	"shop/inventory/internal/domain/entity"
	"shop/inventory/internal/repository"
)

// ReservationServiceImpl implements the ReservationService interface
type ReservationServiceImpl struct {
	reservationRepo repository.ReservationRepository
	stockRepo       repository.StockRepository
}

// NewReservationService creates a new ReservationServiceImpl
func NewReservationService(
	reservationRepo repository.ReservationRepository,
	stockRepo repository.StockRepository,
) ReservationService {
	return &ReservationServiceImpl{
		reservationRepo: reservationRepo,
		stockRepo:       stockRepo,
	}
}

// CreateReservation creates a new stock reservation
func (s *ReservationServiceImpl) CreateReservation(
	ctx context.Context,
	orderID string,
	items []*entity.ReservationItem,
	expirationMinutes int,
) (*entity.Reservation, error) {
	// Validate inputs
	if orderID == "" {
		return nil, errors.New("order ID cannot be empty")
	}
	if len(items) == 0 {
		return nil, errors.New("reservation must have at least one item")
	}
	if expirationMinutes <= 0 {
		expirationMinutes = 30 // Default expiration time: 30 minutes
	}

	// Check if reservation already exists for this order
	existingReservation, err := s.reservationRepo.GetByOrderID(ctx, orderID)
	if err == nil && existingReservation != nil {
		return nil, errors.New("reservation already exists for this order")
	}

	// Reserve stock for each item
	for _, item := range items {
		if item.ProductID <= 0 || item.WarehouseID <= 0 || item.Quantity <= 0 {
			return nil, errors.New("invalid item data")
		}

		// Try to reserve the stock
		_, err := s.stockRepo.ReserveStock(ctx, item.ProductID, item.WarehouseID, item.Quantity)
		if err != nil {
			// If any reservation fails, we need to cancel all previous reservations
			s.rollbackReservations(ctx, items, items[0:]) // This will cancel all reservations processed so far
			return nil, err
		}
	}

	// Create the reservation
	reservation := &entity.Reservation{
		OrderID:   orderID,
		Status:    entity.ReservationPending,
		Items:     items,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(expirationMinutes) * time.Minute),
	}

	// Save the reservation
	err = s.reservationRepo.Create(ctx, reservation)
	if err != nil {
		// If saving fails, we need to cancel all reservations
		s.rollbackReservations(ctx, items, nil)
		return nil, err
	}

	return reservation, nil
}

// GetReservationByID retrieves a reservation by ID
func (s *ReservationServiceImpl) GetReservationByID(ctx context.Context, id int64) (*entity.Reservation, error) {
	if id <= 0 {
		return nil, errors.New("invalid reservation ID")
	}
	return s.reservationRepo.GetByID(ctx, id)
}

// GetReservationByOrderID retrieves a reservation by order ID
func (s *ReservationServiceImpl) GetReservationByOrderID(ctx context.Context, orderID string) (*entity.Reservation, error) {
	if orderID == "" {
		return nil, errors.New("order ID cannot be empty")
	}
	return s.reservationRepo.GetByOrderID(ctx, orderID)
}

// CommitReservation confirms a reservation (e.g., after payment)
func (s *ReservationServiceImpl) CommitReservation(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid reservation ID")
	}

	// Get the reservation
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if the reservation can be committed
	if !reservation.CanCommit() {
		return errors.New("reservation cannot be committed: it is either expired or not in pending state")
	}

	// Commit each reservation item
	for _, item := range reservation.Items {
		_, err := s.stockRepo.CommitReservation(ctx, item.ProductID, item.WarehouseID, item.Quantity)
		if err != nil {
			return err
		}
	}

	// Update the reservation status
	return s.reservationRepo.CommitReservation(ctx, id)
}

// CancelReservation cancels a reservation and returns stock
func (s *ReservationServiceImpl) CancelReservation(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid reservation ID")
	}

	// Get the reservation
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if the reservation can be cancelled
	if !reservation.CanCancel() {
		return errors.New("reservation cannot be cancelled: it is not in pending state")
	}

	// Cancel each reservation item
	for _, item := range reservation.Items {
		_, err := s.stockRepo.CancelReservation(ctx, item.ProductID, item.WarehouseID, item.Quantity)
		if err != nil {
			return err
		}
	}

	// Update the reservation status
	return s.reservationRepo.CancelReservation(ctx, id)
}

// ProcessExpiredReservations finds and cancels all expired reservations
func (s *ReservationServiceImpl) ProcessExpiredReservations(ctx context.Context) (int, error) {
	// Get all expired reservations
	expiredReservations, err := s.reservationRepo.GetExpiredReservations(ctx)
	if err != nil {
		return 0, err
	}

	processedCount := 0
	for _, reservation := range expiredReservations {
		// Skip if the reservation is not actually expired
		if !reservation.IsExpired() {
			continue
		}

		// Cancel each reservation item
		for _, item := range reservation.Items {
			_, err := s.stockRepo.CancelReservation(ctx, item.ProductID, item.WarehouseID, item.Quantity)
			if err != nil {
				continue // Continue with other items even if one fails
			}
		}

		// Update the reservation status to expired
		reservation.Status = entity.ReservationExpired
		err = s.reservationRepo.Update(ctx, reservation)
		if err == nil {
			processedCount++
		}
	}

	return processedCount, nil
}

// ExtendReservation extends the expiration time of a reservation
func (s *ReservationServiceImpl) ExtendReservation(ctx context.Context, id int64, additionalMinutes int) error {
	if id <= 0 {
		return errors.New("invalid reservation ID")
	}
	if additionalMinutes <= 0 {
		return errors.New("additional minutes must be greater than zero")
	}

	// Get the reservation
	reservation, err := s.reservationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if the reservation is in pending state and not already expired
	if reservation.Status != entity.ReservationPending {
		return errors.New("only pending reservations can be extended")
	}
	if reservation.IsExpired() {
		return errors.New("expired reservations cannot be extended")
	}

	// Extend the expiration time
	reservation.ExpiresAt = reservation.ExpiresAt.Add(time.Duration(additionalMinutes) * time.Minute)
	reservation.UpdatedAt = time.Now()

	// Update the reservation
	return s.reservationRepo.Update(ctx, reservation)
}

// rollbackReservations helps in rolling back stock reservations when a multi-item reservation fails
// completedItems: all items in the reservation
// processedItems: items that have been processed so far (subset of completedItems)
func (s *ReservationServiceImpl) rollbackReservations(
	ctx context.Context,
	completedItems []*entity.ReservationItem,
	processedItems []*entity.ReservationItem,
) {
	for _, item := range processedItems {
		s.stockRepo.CancelReservation(ctx, item.ProductID, item.WarehouseID, item.Quantity)
		// We don't check for errors here as we want to try to cancel as many reservations as possible
	}
}
