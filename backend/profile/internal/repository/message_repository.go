package repository

import (
	"context"

	"shop/backend/profile/internal/domain/entity"
)

// MessageRepository defines the interface for Message repository operations
type MessageRepository interface {
	// CreateMessage creates a new message
	CreateMessage(ctx context.Context, message *entity.Message) error

	// GetMessagesByUserID retrieves messages for a user with pagination
	GetMessagesByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*entity.Message, int64, error)

	// GetMessageByID retrieves a message by ID
	GetMessageByID(ctx context.Context, id string) (*entity.Message, error)

	// UpdateMessageStatus updates a message status
	UpdateMessageStatus(ctx context.Context, id string, status int) error

	// DeleteMessage deletes a message
	DeleteMessage(ctx context.Context, id string, userID int64) error
}
