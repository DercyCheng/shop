package service

import (
	"context"
	"fmt"
	"regexp"

	"go.uber.org/zap"

	"shop/backend/profile/internal/domain/entity"
	"shop/backend/profile/internal/repository"
)

// Regular expression for validating Chinese mobile numbers
var mobileRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// AddressServiceImpl implements the AddressService interface
type AddressServiceImpl struct {
	addressRepo repository.AddressRepository
	logger      *zap.Logger
}

// NewAddressService creates a new AddressService implementation
func NewAddressService(
	addressRepo repository.AddressRepository,
	logger *zap.Logger,
) AddressService {
	return &AddressServiceImpl{
		addressRepo: addressRepo,
		logger:      logger,
	}
}

// validateAddress validates address data
func (s *AddressServiceImpl) validateAddress(address *entity.Address) error {
	if address.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if address.Province == "" {
		return fmt.Errorf("province is required")
	}

	if address.City == "" {
		return fmt.Errorf("city is required")
	}

	if address.District == "" {
		return fmt.Errorf("district is required")
	}

	if address.Address == "" {
		return fmt.Errorf("address detail is required")
	}

	if address.SignerName == "" {
		return fmt.Errorf("signer name is required")
	}

	if !mobileRegex.MatchString(address.SignerMobile) {
		return fmt.Errorf("invalid mobile number format")
	}

	return nil
}

// CreateAddress creates a new address for a user
func (s *AddressServiceImpl) CreateAddress(ctx context.Context, address *entity.Address) (*entity.Address, error) {
	// Validate address
	if err := s.validateAddress(address); err != nil {
		return nil, err
	}

	// If this is the default address, clear other default addresses
	if address.IsDefault {
		// Get existing default address
		defaultAddr, err := s.addressRepo.GetDefaultAddress(ctx, address.UserID)
		if err != nil && err.Error() != "record not found" {
			s.logger.Error("Failed to get default address",
				zap.Int64("user_id", address.UserID),
				zap.Error(err))
			return nil, fmt.Errorf("failed to get default address: %w", err)
		}

		// If default address exists, unset it
		if defaultAddr != nil {
			defaultAddr.IsDefault = false
			if err := s.addressRepo.UpdateAddress(ctx, defaultAddr); err != nil {
				s.logger.Error("Failed to update previous default address",
					zap.Int64("address_id", defaultAddr.ID),
					zap.Error(err))
				return nil, fmt.Errorf("failed to update previous default address: %w", err)
			}
		}
	}

	// Create address
	createdAddress, err := s.addressRepo.CreateAddress(ctx, address)
	if err != nil {
		s.logger.Error("Failed to create address",
			zap.Int64("user_id", address.UserID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	s.logger.Info("Created address",
		zap.Int64("address_id", createdAddress.ID),
		zap.Int64("user_id", createdAddress.UserID))
	return createdAddress, nil
}

// UpdateAddress updates an existing address
func (s *AddressServiceImpl) UpdateAddress(ctx context.Context, address *entity.Address) error {
	// Validate address
	if err := s.validateAddress(address); err != nil {
		return err
	}

	// Check if address exists and belongs to user
	existingAddr, err := s.addressRepo.GetAddressByID(ctx, address.ID, address.UserID)
	if err != nil {
		s.logger.Error("Failed to get address",
			zap.Int64("address_id", address.ID),
			zap.Int64("user_id", address.UserID),
			zap.Error(err))
		return fmt.Errorf("failed to get address: %w", err)
	}

	// If updating to set as default, clear other default addresses
	if address.IsDefault && !existingAddr.IsDefault {
		// Get existing default address
		defaultAddr, err := s.addressRepo.GetDefaultAddress(ctx, address.UserID)
		if err != nil && err.Error() != "record not found" {
			s.logger.Error("Failed to get default address",
				zap.Int64("user_id", address.UserID),
				zap.Error(err))
			return fmt.Errorf("failed to get default address: %w", err)
		}

		// If default address exists, unset it
		if defaultAddr != nil && defaultAddr.ID != address.ID {
			defaultAddr.IsDefault = false
			if err := s.addressRepo.UpdateAddress(ctx, defaultAddr); err != nil {
				s.logger.Error("Failed to update previous default address",
					zap.Int64("address_id", defaultAddr.ID),
					zap.Error(err))
				return fmt.Errorf("failed to update previous default address: %w", err)
			}
		}
	}

	// Update address
	if err := s.addressRepo.UpdateAddress(ctx, address); err != nil {
		s.logger.Error("Failed to update address",
			zap.Int64("address_id", address.ID),
			zap.Error(err))
		return fmt.Errorf("failed to update address: %w", err)
	}

	s.logger.Info("Updated address",
		zap.Int64("address_id", address.ID),
		zap.Int64("user_id", address.UserID))
	return nil
}

// DeleteAddress deletes an address by ID and user ID
func (s *AddressServiceImpl) DeleteAddress(ctx context.Context, id, userID int64) error {
	// Check if address exists and belongs to user
	address, err := s.addressRepo.GetAddressByID(ctx, id, userID)
	if err != nil {
		s.logger.Error("Failed to get address",
			zap.Int64("address_id", id),
			zap.Int64("user_id", userID),
			zap.Error(err))
		return fmt.Errorf("failed to get address: %w", err)
	}

	// Delete address
	if err := s.addressRepo.DeleteAddress(ctx, id, userID); err != nil {
		s.logger.Error("Failed to delete address",
			zap.Int64("address_id", id),
			zap.Error(err))
		return fmt.Errorf("failed to delete address: %w", err)
	}

	// If the deleted address was default, set a new default address
	if address.IsDefault {
		addresses, _, err := s.addressRepo.GetAddressList(ctx, userID, 1, 1)
		if err != nil {
			s.logger.Error("Failed to get address list",
				zap.Int64("user_id", userID),
				zap.Error(err))
			return nil // Don't return error as the deletion was successful
		}

		// If there are other addresses, set the first one as default
		if len(addresses) > 0 {
			if err := s.addressRepo.SetDefaultAddress(ctx, addresses[0].ID, userID); err != nil {
				s.logger.Error("Failed to set new default address",
					zap.Int64("address_id", addresses[0].ID),
					zap.Error(err))
				return nil // Don't return error as the deletion was successful
			}
		}
	}

	s.logger.Info("Deleted address",
		zap.Int64("address_id", id),
		zap.Int64("user_id", userID))
	return nil
}

// GetAddressByID retrieves an address by ID
func (s *AddressServiceImpl) GetAddressByID(ctx context.Context, id, userID int64) (*entity.Address, error) {
	return s.addressRepo.GetAddressByID(ctx, id, userID)
}

// GetAddressList retrieves addresses for a user with pagination
func (s *AddressServiceImpl) GetAddressList(ctx context.Context, userID int64, page, pageSize int) ([]*entity.Address, int64, error) {
	return s.addressRepo.GetAddressList(ctx, userID, page, pageSize)
}

// GetDefaultAddress gets the default address for a user
func (s *AddressServiceImpl) GetDefaultAddress(ctx context.Context, userID int64) (*entity.Address, error) {
	return s.addressRepo.GetDefaultAddress(ctx, userID)
}
