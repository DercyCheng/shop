package service

import (
	"context"

	"shop/product/internal/domain/entity"
)

// ProductFilter defines the filter options for product listing
type ProductFilter struct {
	Page       int
	PageSize   int
	CategoryID int64
	BrandID    int64
	Keyword    string
	OnSale     *bool
	ShipFree   *bool
	IsNew      *bool
	IsHot      *bool
	MinPrice   float64
	MaxPrice   float64
	OrderBy    string
	Order      string
}

// ProductService defines the interface for product business logic
type ProductService interface {
	// GetProductByID retrieves a product by ID
	GetProductByID(ctx context.Context, id int64) (*entity.Product, error)

	// GetProductWithRelations retrieves a product with its relations loaded
	GetProductWithRelations(ctx context.Context, id int64) (*entity.Product, error)

	// ListProducts retrieves products based on filters
	ListProducts(ctx context.Context, filter *ProductFilter) ([]*entity.Product, int64, error)

	// GetProductsByIDs retrieves products by a slice of IDs
	GetProductsByIDs(ctx context.Context, ids []int64) ([]*entity.Product, error)

	// CreateProduct adds a new product
	CreateProduct(ctx context.Context, product *entity.Product, images []string) (*entity.Product, error)

	// UpdateProduct modifies an existing product
	UpdateProduct(ctx context.Context, product *entity.Product, images []string) error

	// DeleteProduct removes a product by ID
	DeleteProduct(ctx context.Context, id int64) error

	// AddProductImage adds an image to a product
	AddProductImage(ctx context.Context, productID int64, imageURL string) error

	// DeleteProductImage removes an image from a product
	DeleteProductImage(ctx context.Context, imageID int64) error

	// GetProductImages retrieves all images for a product
	GetProductImages(ctx context.Context, productID int64) ([]string, error)

	// UpdateProductStatus updates product status fields like on_sale, is_new, etc.
	UpdateProductStatus(ctx context.Context, id int64, field string, value bool) error
}
