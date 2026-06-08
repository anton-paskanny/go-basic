package service

import (
	"errors"
	"fmt"
	"log"

	"order-api-cart/clients"
	"order-api-cart/database"
	"order-api-cart/models"

	"github.com/google/uuid"
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

// CreateOrder creates a new order.
//
// Design note: all external HTTP calls (auth + product) are made outside the
// DB transaction so that a transaction rollback never leaves partially-decremented
// inventory in the product service. Inventory is decremented only after the DB
// transaction commits successfully.
func (s *OrderService) CreateOrder(userID string, req *models.OrderRequest, authToken string) (*models.OrderResponse, error) {
	// Validate user exists in auth service
	if err := s.authClient.ValidateUser(userID, authToken); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	// Pre-fetch all products and validate quantities before opening a transaction.
	// Keeping external calls outside the transaction prevents inventory from being
	// decremented in the product service when the DB transaction later rolls back.
	type itemData struct {
		req     models.OrderItemRequest
		product *models.ExternalProduct
	}
	itemsData := make([]itemData, 0, len(req.Items))
	for _, itemReq := range req.Items {
		product, err := s.productClient.GetProductByID(itemReq.ProductID, authToken)
		if err != nil {
			return nil, fmt.Errorf("product not found: %s", itemReq.ProductID)
		}
		if product.Quantity < itemReq.Quantity {
			return nil, fmt.Errorf("insufficient quantity for product %s. Available: %d, Requested: %d",
				product.Name, product.Quantity, itemReq.Quantity)
		}
		itemsData = append(itemsData, itemData{req: itemReq, product: product})
	}

	// DB-only transaction — no external service calls inside
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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
	for _, item := range itemsData {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.req.ProductID,
			Quantity:  item.req.Quantity,
			Price:     item.product.Price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
		total += item.product.Price * float64(item.req.Quantity)
	}

	if err := tx.Model(order).Update("total", total).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update order total: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Decrement inventory after a successful DB commit.
	// If an update fails at this point the order is already persisted; log a
	// warning so ops can reconcile manually. A full solution would use a saga
	// or an outbox pattern.
	for _, item := range itemsData {
		if err := s.productClient.UpdateProductQuantity(item.req.ProductID, -item.req.Quantity, authToken); err != nil {
			log.Printf("WARNING: inventory update failed for product %s after order %s was committed: %v",
				item.req.ProductID, order.ID, err)
		}
	}

	var orderWithItems models.Order
	if err := s.db.Preload("OrderItems").First(&orderWithItems, "id = ?", order.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load order: %w", err)
	}

	return s.orderToResponse(&orderWithItems, authToken), nil
}

// GetOrderByID retrieves an order by ID
func (s *OrderService) GetOrderByID(orderID, authToken string) (*models.OrderResponse, error) {
	if _, err := uuid.Parse(orderID); err != nil {
		return nil, errors.New("order not found")
	}

	var order models.Order

	if err := s.db.Preload("OrderItems").First(&order, "id = ?", orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return s.orderToResponse(&order, authToken), nil
}

// GetOrdersByUserID retrieves paginated orders for a user.
// Pagination is performed at the DB level to avoid loading all orders into memory.
func (s *OrderService) GetOrdersByUserID(userID string, page, limit int, authToken string) ([]models.OrderResponse, int64, error) {
	var total int64
	if err := s.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count user orders: %w", err)
	}

	var orders []models.Order
	offset := (page - 1) * limit
	if err := s.db.Preload("OrderItems").Where("user_id = ?", userID).
		Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get user orders: %w", err)
	}

	responses := make([]models.OrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, *s.orderToResponse(&order, authToken))
	}

	return responses, total, nil
}

// orderToResponse converts an Order model to an OrderResponse, enriching each
// item with product details from the product service. If a product fetch fails,
// a stub is used and a warning is logged rather than failing the whole request.
func (s *OrderService) orderToResponse(order *models.Order, authToken string) *models.OrderResponse {
	var items []models.OrderItemResponse

	for _, item := range order.OrderItems {
		product, err := s.productClient.GetProductByID(item.ProductID, authToken)
		if err != nil {
			log.Printf("WARNING: failed to fetch product %s for order %s response: %v",
				item.ProductID, order.ID, err)
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
