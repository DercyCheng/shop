package service

import (
	"context"

	"shop/backend/profile/internal/domain/entity"
)

// MessageService defines the interface for message operations
type MessageService interface {
	// CreateMessage creates a new message
	CreateMessage(ctx context.Context, message *entity.Message) error

	// GetMessagesByUserID retrieves messages for a user with pagination
	GetMessagesByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*entity.Message, int64, error)

	// GetMessageByID retrieves a message by ID
	GetMessageByID(ctx context.Context, id string) (*entity.Message, error)

	// UpdateMessageStatus updates a message status
	UpdateMessageStatus(ctx context.Context, id string, status int) error
}
