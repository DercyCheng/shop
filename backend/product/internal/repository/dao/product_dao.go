package dao

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"shop/product/internal/domain/entity"
)

// ProductDAO handles data access operations for products
type ProductDAO struct {
	db *gorm.DB
}

// NewProductDAO creates a new product data access object
func NewProductDAO(db *gorm.DB) *ProductDAO {
	return &ProductDAO{db: db}
}

// Create adds a new product to the database
func (d *ProductDAO) Create(ctx context.Context, product *entity.Product) error {
	return d.db.WithContext(ctx).Create(product).Error
}

// Update modifies an existing product
func (d *ProductDAO) Update(ctx context.Context, product *entity.Product) error {
	return d.db.WithContext(ctx).Updates(product).Error
}

// Delete soft-deletes a product by ID
func (d *ProductDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}

// GetByID retrieves a product by ID
func (d *ProductDAO) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	var product entity.Product
	result := d.db.WithContext(ctx).First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found: %v", result.Error)
		}
		return nil, result.Error
	}
	return &product, nil
}

// GetProductWithRelations retrieves a product with its relations loaded
func (d *ProductDAO) GetProductWithRelations(ctx context.Context, id int64) (*entity.Product, error) {
	var product entity.Product
	result := d.db.WithContext(ctx).First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found: %v", result.Error)
		}
		return nil, result.Error
	}

	// Load category
	var category entity.Category
	if err := d.db.WithContext(ctx).First(&category, product.CategoryID).Error; err == nil {
		product.Category = &category
	}

	// Load brand
	var brand entity.Brand
	if err := d.db.WithContext(ctx).First(&brand, product.BrandID).Error; err == nil {
		product.Brand = &brand
	}

	// Load images
	var productImages []entity.ProductImage
	if err := d.db.WithContext(ctx).Where("goods_id = ?", id).Find(&productImages).Error; err == nil {
		images := make([]string, len(productImages))
		for i, img := range productImages {
			images[i] = img.Image
		}
		product.Images = images
	}

	return &product, nil
}

// ListProducts retrieves products based on filters with pagination
func (d *ProductDAO) ListProducts(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var total int64
	offset := (page - 1) * pageSize

	query := d.db.WithContext(ctx).Model(&entity.Product{})

	// Apply filters
	for key, value := range filters {
		if key == "keyword" {
			query = query.Where("name LIKE ? OR goods_brief LIKE ?", "%"+value.(string)+"%", "%"+value.(string)+"%")
		} else if key == "min_price" {
			query = query.Where("shop_price >= ?", value)
		} else if key == "max_price" {
			query = query.Where("shop_price <= ?", value)
		} else {
			query = query.Where(key+" = ?", value)
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Order by id desc by default
	orderBy := "id DESC"
	if orderByVal, exists := filters["order_by"]; exists {
		orderBy = orderByVal.(string)
	}

	// Get paginated results
	if err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ListByIDs retrieves products by a slice of IDs
func (d *ProductDAO) ListByIDs(ctx context.Context, ids []int64) ([]*entity.Product, error) {
	var products []*entity.Product
	if err := d.db.WithContext(ctx).Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// AddProductImage adds an image to a product
func (d *ProductDAO) AddProductImage(ctx context.Context, productID int64, imageURL string) error {
	productImage := entity.ProductImage{
		GoodsID: productID,
		Image:   imageURL,
	}
	return d.db.WithContext(ctx).Create(&productImage).Error
}

// DeleteProductImage removes an image from a product
func (d *ProductDAO) DeleteProductImage(ctx context.Context, imageID int64) error {
	return d.db.WithContext(ctx).Delete(&entity.ProductImage{}, imageID).Error
}

// GetProductImages retrieves all images for a product
func (d *ProductDAO) GetProductImages(ctx context.Context, productID int64) ([]string, error) {
	var productImages []entity.ProductImage
	if err := d.db.WithContext(ctx).Where("goods_id = ?", productID).Find(&productImages).Error; err != nil {
		return nil, err
	}

	images := make([]string, len(productImages))
	for i, img := range productImages {
		images[i] = img.Image
	}
	return images, nil
}
