package service

import (
	"errors"
	"fmt"

	"order-api-cart/clients"
	"order-api-cart/database"
	"order-api-cart/models"

	"gorm.io/gorm"
)

// OrderService handles order business logic
type OrderService struct {
	db            *gorm.DB
	authClient    *clients.AuthServiceClient
	productClient *clients.ProductServiceClient
}

// NewOrderService creates a new order service
func NewOrderService(authServiceURL, productServiceURL string) *OrderService {
	return &OrderService{
		db:            database.GetDB(),
		authClient:    clients.NewAuthServiceClient(authServiceURL),
		productClient: clients.NewProductServiceClient(productServiceURL),
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(userID string, req *models.OrderRequest) (*models.OrderResponse, error) {
	// Validate user exists in auth service
	if err := s.authClient.ValidateUser(userID); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	// Start a transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the order
	order := &models.Order{
		UserID: userID,
		Status: "pending",
		Total:  0,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	var total float64
	var orderItems []models.OrderItem

	// Process each item in the order
	for _, itemReq := range req.Items {
		// Get product details from product service
		product, err := s.productClient.GetProductByID(itemReq.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("product not found: %s", itemReq.ProductID)
		}

		// Check quantity availability
		if product.Quantity < itemReq.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient quantity for product %s. Available: %d, Requested: %d",
				product.Name, product.Quantity, itemReq.Quantity)
		}

		// Create order item
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}

		// Update product quantity in product service
		if err := s.productClient.UpdateProductQuantity(itemReq.ProductID, -itemReq.Quantity); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update product quantity: %w", err)
		}

		orderItems = append(orderItems, orderItem)
		total += product.Price * float64(itemReq.Quantity)
	}

	// Update order total
	if err := tx.Model(order).Update("total", total).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update order total: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Load order with items for response
	var orderWithItems models.Order
	if err := s.db.Preload("OrderItems").First(&orderWithItems, "id = ?", order.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load order: %w", err)
	}

	return s.orderToResponse(&orderWithItems), nil
}

// GetOrderByID retrieves an order by ID
func (s *OrderService) GetOrderByID(orderID string) (*models.OrderResponse, error) {
	var order models.Order

	if err := s.db.Preload("OrderItems").First(&order, "id = ?", orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return s.orderToResponse(&order), nil
}

// GetOrdersByUserID retrieves all orders for a specific user
func (s *OrderService) GetOrdersByUserID(userID string) ([]models.OrderResponse, error) {
	var orders []models.Order

	if err := s.db.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	var responses []models.OrderResponse
	for _, order := range orders {
		responses = append(responses, *s.orderToResponse(&order))
	}

	return responses, nil
}

// orderToResponse converts Order model to OrderResponse
func (s *OrderService) orderToResponse(order *models.Order) *models.OrderResponse {
	var items []models.OrderItemResponse

	for _, item := range order.OrderItems {
		// Fetch product details for each item
		product, err := s.productClient.GetProductByID(item.ProductID)
		if err != nil {
			// If product fetch fails, create a minimal product response
			product = &models.ExternalProduct{
				ID:    item.ProductID,
				Name:  "Product not available",
				Price: item.Price,
			}
		}

		items = append(items, models.OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Product:   *product,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return &models.OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Status:    order.Status,
		Total:     order.Total,
		Items:     items,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}
