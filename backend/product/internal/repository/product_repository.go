package repository

import (
	"context"

	"shop/product/internal/domain/entity"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	// Create adds a new product
	Create(ctx context.Context, product *entity.Product) error

	// Update modifies an existing product
	Update(ctx context.Context, product *entity.Product) error

	// Delete removes a product by ID
	Delete(ctx context.Context, id int64) error

	// GetByID retrieves a product by ID
	GetByID(ctx context.Context, id int64) (*entity.Product, error)

	// GetProductWithRelations retrieves a product with its relations loaded
	GetProductWithRelations(ctx context.Context, id int64) (*entity.Product, error)

	// ListProducts retrieves products based on filters
	ListProducts(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*entity.Product, int64, error)

	// ListByIDs retrieves products by a slice of IDs
	ListByIDs(ctx context.Context, ids []int64) ([]*entity.Product, error)

	// AddProductImage adds an image to a product
	AddProductImage(ctx context.Context, productID int64, imageURL string) error

	// DeleteProductImage removes an image from a product
	DeleteProductImage(ctx context.Context, imageID int64) error

	// GetProductImages retrieves all images for a product
	GetProductImages(ctx context.Context, productID int64) ([]string, error)
}
