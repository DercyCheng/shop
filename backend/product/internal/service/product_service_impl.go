package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shop/product/internal/domain/entity"
	"shop/product/internal/repository"
)

// ProductServiceImpl implements ProductService interface
type ProductServiceImpl struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	brandRepo    repository.BrandRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, brandRepo repository.BrandRepository) ProductService {
	return &ProductServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
	}
}

// GetProductByID retrieves a product by ID
func (s *ProductServiceImpl) GetProductByID(ctx context.Context, id int64) (*entity.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

// GetProductWithRelations retrieves a product with its relations loaded
func (s *ProductServiceImpl) GetProductWithRelations(ctx context.Context, id int64) (*entity.Product, error) {
	return s.productRepo.GetProductWithRelations(ctx, id)
}

// ListProducts retrieves products based on filters
func (s *ProductServiceImpl) ListProducts(ctx context.Context, filter *ProductFilter) ([]*entity.Product, int64, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	// Construct filters map from ProductFilter struct
	filters := make(map[string]interface{})

	if filter.CategoryID > 0 {
		filters["category_id"] = filter.CategoryID
	}

	if filter.BrandID > 0 {
		filters["brands_id"] = filter.BrandID
	}

	if filter.Keyword != "" {
		filters["keyword"] = filter.Keyword
	}

	if filter.OnSale != nil {
		filters["on_sale"] = *filter.OnSale
	}

	if filter.ShipFree != nil {
		filters["ship_free"] = *filter.ShipFree
	}

	if filter.IsNew != nil {
		filters["is_new"] = *filter.IsNew
	}

	if filter.IsHot != nil {
		filters["is_hot"] = *filter.IsHot
	}

	if filter.MinPrice > 0 {
		filters["min_price"] = filter.MinPrice
	}

	if filter.MaxPrice > 0 {
		filters["max_price"] = filter.MaxPrice
	}

	// Set ordering
	orderBy := "id DESC" // Default ordering
	if filter.OrderBy != "" {
		order := "DESC"
		if filter.Order == "asc" {
			order = "ASC"
		}
		orderBy = fmt.Sprintf("%s %s", filter.OrderBy, order)
	}
	filters["order_by"] = orderBy

	return s.productRepo.ListProducts(ctx, filter.Page, filter.PageSize, filters)
}

// GetProductsByIDs retrieves products by a slice of IDs
func (s *ProductServiceImpl) GetProductsByIDs(ctx context.Context, ids []int64) ([]*entity.Product, error) {
	if len(ids) == 0 {
		return []*entity.Product{}, nil
	}
	return s.productRepo.ListByIDs(ctx, ids)
}

// CreateProduct adds a new product
func (s *ProductServiceImpl) CreateProduct(ctx context.Context, product *entity.Product, images []string) (*entity.Product, error) {
	// Validate required fields
	if product.Name == "" {
		return nil, errors.New("product name is required")
	}

	if product.CategoryID <= 0 {
		return nil, errors.New("category ID is required")
	}

	if product.BrandID <= 0 {
		return nil, errors.New("brand ID is required")
	}

	// Verify category exists
	_, err := s.categoryRepo.GetByID(ctx, product.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %v", err)
	}

	// Verify brand exists
	_, err = s.brandRepo.GetByID(ctx, product.BrandID)
	if err != nil {
		return nil, fmt.Errorf("invalid brand: %v", err)
	}

	// Set timestamps
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	// Create product
	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	// Add product images if any
	for _, imageURL := range images {
		if err := s.productRepo.AddProductImage(ctx, product.ID, imageURL); err != nil {
			// Log error but continue
			fmt.Printf("Error adding product image: %v\n", err)
		}
	}

	// Load product with relations (including the images we just added)
	return s.productRepo.GetProductWithRelations(ctx, product.ID)
}

// UpdateProduct modifies an existing product
func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, product *entity.Product, images []string) error {
	// Verify product exists
	existingProduct, err := s.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return fmt.Errorf("product not found: %v", err)
	}

	// If category changed, verify new category exists
	if product.CategoryID > 0 && product.CategoryID != existingProduct.CategoryID {
		_, err := s.categoryRepo.GetByID(ctx, product.CategoryID)
		if err != nil {
			return fmt.Errorf("invalid category: %v", err)
		}
	} else {
		// Keep existing category if not specified
		product.CategoryID = existingProduct.CategoryID
	}

	// If brand changed, verify new brand exists
	if product.BrandID > 0 && product.BrandID != existingProduct.BrandID {
		_, err := s.brandRepo.GetByID(ctx, product.BrandID)
		if err != nil {
			return fmt.Errorf("invalid brand: %v", err)
		}
	} else {
		// Keep existing brand if not specified
		product.BrandID = existingProduct.BrandID
	}

	// Update timestamp
	product.UpdatedAt = time.Now()

	// Update product
	if err := s.productRepo.Update(ctx, product); err != nil {
		return err
	}

	// If new images provided, update them
	if images != nil && len(images) > 0 {
		// Get existing images to delete
		existingImages, err := s.productRepo.GetProductImages(ctx, product.ID)
		if err == nil {
			// TODO: Delete existing images or implement a more sophisticated image management system
		}

		// Add new images
		for _, imageURL := range images {
			if err := s.productRepo.AddProductImage(ctx, product.ID, imageURL); err != nil {
				// Log error but continue
				fmt.Printf("Error adding product image: %v\n", err)
			}
		}
	}

	return nil
}

// DeleteProduct removes a product by ID
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id int64) error {
	return s.productRepo.Delete(ctx, id)
}

// AddProductImage adds an image to a product
func (s *ProductServiceImpl) AddProductImage(ctx context.Context, productID int64, imageURL string) error {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %v", err)
	}

	return s.productRepo.AddProductImage(ctx, productID, imageURL)
}

// DeleteProductImage removes an image from a product
func (s *ProductServiceImpl) DeleteProductImage(ctx context.Context, imageID int64) error {
	return s.productRepo.DeleteProductImage(ctx, imageID)
}

// GetProductImages retrieves all images for a product
func (s *ProductServiceImpl) GetProductImages(ctx context.Context, productID int64) ([]string, error) {
	return s.productRepo.GetProductImages(ctx, productID)
}

// UpdateProductStatus updates product status fields like on_sale, is_new, etc.
func (s *ProductServiceImpl) UpdateProductStatus(ctx context.Context, id int64, field string, value bool) error {
	// Get existing product
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %v", err)
	}

	// Update specified field
	switch field {
	case "on_sale":
		product.OnSale = value
	case "ship_free":
		product.ShipFree = value
	case "is_new":
		product.IsNew = value
	case "is_hot":
		product.IsHot = value
	default:
		return fmt.Errorf("invalid status field: %s", field)
	}

	// Update timestamp
	product.UpdatedAt = time.Now()

	// Save changes
	return s.productRepo.Update(ctx, product)
}
