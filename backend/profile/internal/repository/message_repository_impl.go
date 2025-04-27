package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"shop/backend/profile/internal/domain/entity"
)

// MessageRepositoryImpl implements the MessageRepository interface
type MessageRepositoryImpl struct {
	client     *mongo.Client
	database   string
	collection string
	logger     *zap.Logger
}

// NewMessageRepository creates a new MessageRepository implementation
func NewMessageRepository(
	client *mongo.Client,
	database string,
	collection string,
	logger *zap.Logger,
) MessageRepository {
	return &MessageRepositoryImpl{
		client:     client,
		database:   database,
		collection: collection,
		logger:     logger,
	}
}

// getCollection returns the MongoDB collection
func (r *MessageRepositoryImpl) getCollection() *mongo.Collection {
	return r.client.Database(r.database).Collection(r.collection)
}

// CreateMessage creates a new message
func (r *MessageRepositoryImpl) CreateMessage(ctx context.Context, message *entity.Message) error {
	coll := r.getCollection()

	// Prepare the document
	doc := bson.M{
		"user_id":      message.UserID,
		"message_type": message.MessageType,
		"subject":      message.Subject,
		"content":      message.Content,
		"file":         message.File,
		"images":       message.Images,
		"status":       message.Status,
		"created_at":   message.CreatedAt,
		"updated_at":   message.UpdatedAt,
	}

	// Insert the document
	result, err := coll.InsertOne(ctx, doc)
	if err != nil {
		r.logger.Error("Failed to create message",
			zap.Int64("user_id", message.UserID),
			zap.Error(err))
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Update the ID in the entity
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		message.ID = oid.Hex()
	}

	return nil
}

// GetMessagesByUserID retrieves messages for a user with pagination
func (r *MessageRepositoryImpl) GetMessagesByUserID(ctx context.Context, userID int64, page, pageSize int) ([]*entity.Message, int64, error) {
	coll := r.getCollection()

	// Prepare the filter
	filter := bson.M{"user_id": userID}

	// Count total documents
	total, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count messages",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	// Prepare options for pagination and sorting
	findOptions := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))

	// Find documents
	cursor, err := coll.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Error("Failed to find messages",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode documents
	messages := make([]*entity.Message, 0)
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode message", zap.Error(err))
			continue
		}

		message := &entity.Message{
			ID:          doc["_id"].(primitive.ObjectID).Hex(),
			UserID:      doc["user_id"].(int64),
			MessageType: int(doc["message_type"].(int32)),
			Subject:     doc["subject"].(string),
			Content:     doc["content"].(string),
			Status:      int(doc["status"].(int32)),
			CreatedAt:   doc["created_at"].(primitive.DateTime).Time(),
			UpdatedAt:   doc["updated_at"].(primitive.DateTime).Time(),
		}

		// Handle optional fields
		if file, ok := doc["file"].(string); ok {
			message.File = file
		}

		if images, ok := doc["images"].(primitive.A); ok {
			message.Images = make([]string, 0, len(images))
			for _, img := range images {
				if imgStr, ok := img.(string); ok {
					message.Images = append(message.Images, imgStr)
				}
			}
		}

		messages = append(messages, message)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error", zap.Error(err))
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}

	return messages, total, nil
}

// GetMessageByID retrieves a message by ID
func (r *MessageRepositoryImpl) GetMessageByID(ctx context.Context, id string) (*entity.Message, error) {
	coll := r.getCollection()

	// Convert ID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid message ID format",
			zap.String("message_id", id),
			zap.Error(err))
		return nil, fmt.Errorf("invalid message ID format: %w", err)
	}

	// Prepare the filter
	filter := bson.M{"_id": objectID}

	// Find the document
	var doc bson.M
	if err := coll.FindOne(ctx, filter).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("message not found")
		}
		r.logger.Error("Failed to find message",
			zap.String("message_id", id),
			zap.Error(err))
		return nil, fmt.Errorf("failed to find message: %w", err)
	}

	// Convert to entity
	message := &entity.Message{
		ID:          id,
		UserID:      doc["user_id"].(int64),
		MessageType: int(doc["message_type"].(int32)),
		Subject:     doc["subject"].(string),
		Content:     doc["content"].(string),
		Status:      int(doc["status"].(int32)),
		CreatedAt:   doc["created_at"].(primitive.DateTime).Time(),
		UpdatedAt:   doc["updated_at"].(primitive.DateTime).Time(),
	}

	// Handle optional fields
	if file, ok := doc["file"].(string); ok {
		message.File = file
	}

	if images, ok := doc["images"].(primitive.A); ok {
		message.Images = make([]string, 0, len(images))
		for _, img := range images {
			if imgStr, ok := img.(string); ok {
				message.Images = append(message.Images, imgStr)
			}
		}
	}

	return message, nil
}

// UpdateMessageStatus updates a message status
func (r *MessageRepositoryImpl) UpdateMessageStatus(ctx context.Context, id string, status int) error {
	coll := r.getCollection()

	// Convert ID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid message ID format",
			zap.String("message_id", id),
			zap.Error(err))
		return fmt.Errorf("invalid message ID format: %w", err)
	}

	// Prepare the filter and update
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	// Update the document
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("Failed to update message status",
			zap.String("message_id", id),
			zap.Int("status", status),
			zap.Error(err))
		return fmt.Errorf("failed to update message status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

// DeleteMessage deletes a message
func (r *MessageRepositoryImpl) DeleteMessage(ctx context.Context, id string, userID int64) error {
	coll := r.getCollection()

	// Convert ID string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid message ID format",
			zap.String("message_id", id),
			zap.Error(err))
		return fmt.Errorf("invalid message ID format: %w", err)
	}

	// Prepare the filter
	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	// Delete the document
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to delete message",
			zap.String("message_id", id),
			zap.Error(err))
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("message not found or not authorized")
	}

	return nil
}
