package repository

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"shop/backend/profile/internal/domain/entity"
)

// UserFavModel represents the database model for user favorites
type UserFavModel struct {
	ID        int64      `gorm:"primaryKey"`
	User      int64      `gorm:"column:user;not null;index"`
	Goods     int64      `gorm:"column:goods;not null;index"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

// TableName specifies the table name for UserFavModel
func (UserFavModel) TableName() string {
	return "user_fav"
}

// UserFavRepositoryImpl implements the UserFavRepository interface
type UserFavRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewUserFavRepository creates a new UserFavRepository implementation
func NewUserFavRepository(db *gorm.DB, logger *zap.Logger) UserFavRepository {
	return &UserFavRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// toEntity converts a UserFavModel to a UserFav entity
func (r *UserFavRepositoryImpl) toEntity(model *UserFavModel) *entity.UserFav {
	return &entity.UserFav{
		ID:        model.ID,
		UserID:    model.User,
		GoodsID:   model.Goods,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// AddUserFav adds a favorite for a user
func (r *UserFavRepositoryImpl) AddUserFav(ctx context.Context, userID, goodsID int64) error {
	model := &UserFavModel{
		User:      userID,
		Goods:     goodsID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		r.logger.Error("Failed to add user favorite",
			zap.Int64("user_id", userID),
			zap.Int64("goods_id", goodsID),
			zap.Error(result.Error))
		return result.Error
	}

	return nil
}

// DeleteUserFav removes a favorite for a user
func (r *UserFavRepositoryImpl) DeleteUserFav(ctx context.Context, userID, goodsID int64) error {
	result := r.db.WithContext(ctx).
		Where("user = ? AND goods = ?", userID, goodsID).
		Delete(&UserFavModel{})

	if result.Error != nil {
		r.logger.Error("Failed to delete user favorite",
			zap.Int64("user_id", userID),
			zap.Int64("goods_id", goodsID),
			zap.Error(result.Error))
		return result.Error
	}

	return nil
}

// CheckUserFav checks if a user has favorited a goods
func (r *UserFavRepositoryImpl) CheckUserFav(ctx context.Context, userID, goodsID int64) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&UserFavModel{}).
		Where("user = ? AND goods = ?", userID, goodsID).
		Count(&count)

	if result.Error != nil {
		r.logger.Error("Failed to check user favorite",
			zap.Int64("user_id", userID),
			zap.Int64("goods_id", goodsID),
			zap.Error(result.Error))
		return false, result.Error
	}

	return count > 0, nil
}

// GetUserFavList retrieves a list of user favorites with pagination
func (r *UserFavRepositoryImpl) GetUserFavList(ctx context.Context, userID int64, page, pageSize int) ([]*entity.UserFav, int64, error) {
	var models []*UserFavModel
	var count int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&UserFavModel{}).
		Where("user = ?", userID).
		Count(&count).Error; err != nil {
		r.logger.Error("Failed to count user favorites",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return nil, 0, err
	}

	// Get data with pagination
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).
		Where("user = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		r.logger.Error("Failed to get user favorites",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return nil, 0, err
	}

	// Convert models to entities
	favorites := make([]*entity.UserFav, 0, len(models))
	for _, model := range models {
		favorites = append(favorites, r.toEntity(model))
	}

	return favorites, count, nil
}
