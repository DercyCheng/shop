package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"shop/backend/profile/internal/domain/entity"
	"shop/backend/profile/internal/repository"
)

// MessageServiceImpl implements the MessageService interface
type MessageServiceImpl struct {
	messageRepo repository.MessageRepository
	logger      *zap.Logger
}

// NewMessageService creates a new MessageService implementation
func NewMessageService(
	messageRepo repository.MessageRepository,
	logger *zap.Logger,
) MessageService {
	return &MessageServiceImpl{
		messageRepo: messageRepo,
		logger:      logger,
	}
}

// validateMessage validates message data
func (s *MessageServiceImpl) validateMessage(message *entity.Message) error {
	if message.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if message.Subject == "" {
		return fmt.Errorf("subject is required")
	}

	if message.Content == "" {
		return fmt.Errorf("message content is required")
	}

	// Validate message type
	if _, ok := entity.MessageTypeMap[message.MessageType]; !ok {
		return fmt.Errorf("invalid message type")
	}

	return nil
}

// CreateMessage creates a new message
func (s *MessageServiceImpl) CreateMessage(ctx context.Context, message *entity.Message) error {
	// Validate message
	if err := s.validateMessage(message); err != nil {
		return err
	}

	// Set defaults
	message.Status = 0 // Unprocessed
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	// Create message
	if err := s.messageRepo.CreateMessage(ctx, message); err != nil {
		s.logger.Error("Failed to create message",
			zap.Int64("user_id", message.UserID),
			zap.Error(err))
		return fmt.Errorf("failed to create message: %w", err)
	}

	s.logger.Info("Created message",
		zap.String("message_id", message.ID),
		zap.Int64("user_id", message.UserID))
	return nil
}

// GetMessagesByUserID retrieves messages for a user with pagination
func (s *MessageServiceImpl) GetMessagesByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*entity.Message, int64, error) {
	messages, total, err := s.messageRepo.GetMessagesByUserID(ctx, userID, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to get messages",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, total, nil
}

// GetMessageByID retrieves a message by ID
func (s *MessageServiceImpl) GetMessageByID(ctx context.Context, id string) (*entity.Message, error) {
	message, err := s.messageRepo.GetMessageByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get message",
			zap.String("message_id", id),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return message, nil
}

// UpdateMessageStatus updates a message status
func (s *MessageServiceImpl) UpdateMessageStatus(ctx context.Context, id string, status int) error {
	// Validate status
	if _, ok := entity.MessageStatusMap[status]; !ok {
		return fmt.Errorf("invalid message status")
	}

	// Check if message exists
	_, err := s.messageRepo.GetMessageByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get message",
			zap.String("message_id", id),
			zap.Error(err))
		return fmt.Errorf("failed to get message: %w", err)
	}

	// Update status
	if err := s.messageRepo.UpdateMessageStatus(ctx, id, status); err != nil {
		s.logger.Error("Failed to update message status",
			zap.String("message_id", id),
			zap.Int("status", status),
			zap.Error(err))
		return fmt.Errorf("failed to update message status: %w", err)
	}

	s.logger.Info("Updated message status",
		zap.String("message_id", id),
		zap.Int("status", status))
	return nil
}
