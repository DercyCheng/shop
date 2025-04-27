package repository

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"shop/backend/profile/internal/domain/entity"
)

// AddressModel represents the database model for addresses
type AddressModel struct {
	ID           int64      `gorm:"primaryKey"`
	User         int64      `gorm:"column:user;not null;index"`
	Province     string     `gorm:"column:province;not null"`
	City         string     `gorm:"column:city;not null"`
	District     string     `gorm:"column:district;not null"`
	Address      string     `gorm:"column:address;not null"`
	SignerName   string     `gorm:"column:signer_name;not null"`
	SignerMobile string     `gorm:"column:signer_mobile;not null"`
	IsDefault    bool       `gorm:"column:is_default;default:0"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

// TableName specifies the table name for AddressModel
func (AddressModel) TableName() string {
	return "address"
}

// AddressRepositoryImpl implements the AddressRepository interface
type AddressRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewAddressRepository creates a new AddressRepository implementation
func NewAddressRepository(db *gorm.DB, logger *zap.Logger) AddressRepository {
	return &AddressRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// toEntity converts an AddressModel to an Address entity
func (r *AddressRepositoryImpl) toEntity(model *AddressModel) *entity.Address {
	return &entity.Address{
		ID:           model.ID,
		UserID:       model.User,
		Province:     model.Province,
		City:         model.City,
		District:     model.District,
		Address:      model.Address,
		SignerName:   model.SignerName,
		SignerMobile: model.SignerMobile,
		IsDefault:    model.IsDefault,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

// toModel converts an Address entity to an AddressModel
func (r *AddressRepositoryImpl) toModel(entity *entity.Address) *AddressModel {
	return &AddressModel{
		ID:           entity.ID,
		User:         entity.UserID,
		Province:     entity.Province,
		City:         entity.City,
		District:     entity.District,
		Address:      entity.Address,
		SignerName:   entity.SignerName,
		SignerMobile: entity.SignerMobile,
		IsDefault:    entity.IsDefault,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}

// CreateAddress creates a new address for a user
func (r *AddressRepositoryImpl) CreateAddress(ctx context.Context, address *entity.Address) (*entity.Address, error) {
	now := time.Now()
	address.CreatedAt = now
	address.UpdatedAt = now

	model := r.toModel(address)
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		r.logger.Error("Failed to create address",
			zap.Int64("user_id", address.UserID),
			zap.Error(result.Error))
		return nil, result.Error
	}

	// Update the ID in the entity
	address.ID = model.ID
	return address, nil
}

// UpdateAddress updates an existing address
func (r *AddressRepositoryImpl) UpdateAddress(ctx context.Context, address *entity.Address) error {
	address.UpdatedAt = time.Now()
	model := r.toModel(address)

	result := r.db.WithContext(ctx).
		Where("id = ? AND user = ?", address.ID, address.UserID).
		Updates(model)

	if result.Error != nil {
		r.logger.Error("Failed to update address",
			zap.Int64("address_id", address.ID),
			zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("address not found")
	}

	return nil
}

// DeleteAddress deletes an address by ID and user ID
func (r *AddressRepositoryImpl) DeleteAddress(ctx context.Context, id, userID int64) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user = ?", id, userID).
		Delete(&AddressModel{})

	if result.Error != nil {
		r.logger.Error("Failed to delete address",
			zap.Int64("address_id", id),
			zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("address not found")
	}

	return nil
}

// GetAddressByID retrieves an address by ID
func (r *AddressRepositoryImpl) GetAddressByID(ctx context.Context, id, userID int64) (*entity.Address, error) {
	var model AddressModel
	result := r.db.WithContext(ctx).
		Where("id = ? AND user = ?", id, userID).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("address not found")
		}
		r.logger.Error("Failed to get address",
			zap.Int64("address_id", id),
			zap.Error(result.Error))
		return nil, result.Error
	}

	return r.toEntity(&model), nil
}

// GetAddressList retrieves addresses for a user with pagination
func (r *AddressRepositoryImpl) GetAddressList(ctx context.Context, userID int64, page, pageSize int) ([]*entity.Address, int64, error) {
	var models []*AddressModel
	var count int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&AddressModel{}).
		Where("user = ?", userID).
		Count(&count).Error; err != nil {
		r.logger.Error("Failed to count addresses",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return nil, 0, err
	}

	// Get data with pagination
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Where("user = ?", userID).
		Order("is_default DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&models).Error; err != nil {
		r.logger.Error("Failed to get addresses",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return nil, 0, err
	}

	// Convert models to entities
	addresses := make([]*entity.Address, 0, len(models))
	for _, model := range models {
		addresses = append(addresses, r.toEntity(model))
	}

	return addresses, count, nil
}

// GetDefaultAddress gets the default address for a user
func (r *AddressRepositoryImpl) GetDefaultAddress(ctx context.Context, userID int64) (*entity.Address, error) {
	var model AddressModel
	result := r.db.WithContext(ctx).
		Where("user = ? AND is_default = ?", userID, true).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("default address not found")
		}
		r.logger.Error("Failed to get default address",
			zap.Int64("user_id", userID),
			zap.Error(result.Error))
		return nil, result.Error
	}

	return r.toEntity(&model), nil
}

// SetDefaultAddress sets an address as the default address for a user
func (r *AddressRepositoryImpl) SetDefaultAddress(ctx context.Context, id, userID int64) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(tx.Error))
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// First, unset all default addresses for the user
	if err := tx.Model(&AddressModel{}).
		Where("user = ? AND is_default = ?", userID, true).
		Update("is_default", false).Error; err != nil {
		tx.Rollback()
		r.logger.Error("Failed to unset default addresses",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return err
	}

	// Then, set the specified address as default
	result := tx.Model(&AddressModel{}).
		Where("id = ? AND user = ?", id, userID).
		Update("is_default", true)

	if result.Error != nil {
		tx.Rollback()
		r.logger.Error("Failed to set default address",
			zap.Int64("address_id", id),
			zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("address not found")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		r.logger.Error("Failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}
