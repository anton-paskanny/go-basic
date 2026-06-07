package service

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"order-api-stat/models"
)

// ProductService handles product-related business logic
type ProductService struct {
	db *gorm.DB
}

// NewProductService creates a new product service
func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Category:    req.Category,
		SKU:         req.SKU,
		Images:      req.Images,
	}

	if err := s.db.Create(product).Error; err != nil {
		if isUniqueConstraintError(err) {
			return nil, fmt.Errorf("product with SKU '%s' already exists", req.SKU)
		}
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(id uint) (*models.Product, error) {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// ListProducts retrieves a list of products with pagination
func (s *ProductService) ListProducts(page, limit int, category string) (*models.ProductListResponse, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{})

	// Filter by category if provided
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	// Convert to response format
	productResponses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = models.ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
			Category:    product.Category,
			SKU:         product.SKU,
			Images:      product.Images,
			CreatedAt:   product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &models.ProductListResponse{
		Products: productResponses,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(id uint, req *models.UpdateProductRequest) (*models.Product, error) {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.SKU != nil {
		updates["sku"] = *req.SKU
	}
	if req.Images != nil {
		updates["images"] = req.Images
	}

	if err := s.db.Model(&product).Updates(updates).Error; err != nil {
		if req.SKU != nil && isUniqueConstraintError(err) {
			return nil, fmt.Errorf("product with SKU '%s' already exists", *req.SKU)
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Fetch updated product
	if err := s.db.First(&product, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated product: %w", err)
	}

	return &product, nil
}

// isUniqueConstraintError reports whether err is a PostgreSQL unique constraint violation.
func isUniqueConstraintError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "23505") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint")
}

// DeleteProduct soft deletes a product
func (s *ProductService) DeleteProduct(id uint) error {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("product with ID %d not found", id)
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	if err := s.db.Delete(&product).Error; err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
