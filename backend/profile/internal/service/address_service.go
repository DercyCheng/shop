package service

import (
	"context"

	"shop/backend/profile/internal/domain/entity"
)

// AddressService defines the interface for address operations
type AddressService interface {
	// CreateAddress creates a new address for a user
	CreateAddress(ctx context.Context, address *entity.Address) (*entity.Address, error)

	// UpdateAddress updates an existing address
	UpdateAddress(ctx context.Context, address *entity.Address) error

	// DeleteAddress deletes an address by ID and user ID
	DeleteAddress(ctx context.Context, id, userID int64) error

	// GetAddressByID retrieves an address by ID
	GetAddressByID(ctx context.Context, id, userID int64) (*entity.Address, error)

	// GetAddressList retrieves addresses for a user with pagination
	GetAddressList(ctx context.Context, userID int64, page, pageSize int) ([]*entity.Address, int64, error)

	// GetDefaultAddress gets the default address for a user
	GetDefaultAddress(ctx context.Context, userID int64) (*entity.Address, error)
}
