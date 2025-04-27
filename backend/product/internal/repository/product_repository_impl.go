package repository

import (
	"context"

	"shop/product/internal/domain/entity"
	"shop/product/internal/repository/dao"
)

// ProductRepositoryImpl implements ProductRepository interface
type ProductRepositoryImpl struct {
	productDAO *dao.ProductDAO
}

// NewProductRepository creates a new product repository
func NewProductRepository(productDAO *dao.ProductDAO) ProductRepository {
	return &ProductRepositoryImpl{
		productDAO: productDAO,
	}
}

// Create adds a new product
func (r *ProductRepositoryImpl) Create(ctx context.Context, product *entity.Product) error {
	return r.productDAO.Create(ctx, product)
}

// Update modifies an existing product
func (r *ProductRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {
	return r.productDAO.Update(ctx, product)
}

// Delete removes a product by ID
func (r *ProductRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.productDAO.Delete(ctx, id)
}

// GetByID retrieves a product by ID
func (r *ProductRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	return r.productDAO.GetByID(ctx, id)
}

// GetProductWithRelations retrieves a product with its relations loaded
func (r *ProductRepositoryImpl) GetProductWithRelations(ctx context.Context, id int64) (*entity.Product, error) {
	return r.productDAO.GetProductWithRelations(ctx, id)
}

// ListProducts retrieves products based on filters
func (r *ProductRepositoryImpl) ListProducts(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*entity.Product, int64, error) {
	return r.productDAO.ListProducts(ctx, page, pageSize, filters)
}

// ListByIDs retrieves products by a slice of IDs
func (r *ProductRepositoryImpl) ListByIDs(ctx context.Context, ids []int64) ([]*entity.Product, error) {
	return r.productDAO.ListByIDs(ctx, ids)
}

// AddProductImage adds an image to a product
func (r *ProductRepositoryImpl) AddProductImage(ctx context.Context, productID int64, imageURL string) error {
	return r.productDAO.AddProductImage(ctx, productID, imageURL)
}

// DeleteProductImage removes an image from a product
func (r *ProductRepositoryImpl) DeleteProductImage(ctx context.Context, imageID int64) error {
	return r.productDAO.DeleteProductImage(ctx, imageID)
}

// GetProductImages retrieves all images for a product
func (r *ProductRepositoryImpl) GetProductImages(ctx context.Context, productID int64) ([]string, error) {
	return r.productDAO.GetProductImages(ctx, productID)
}
