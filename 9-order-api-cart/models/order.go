package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExternalUser represents a user from the auth service
type ExternalUser struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ExternalProduct represents a product from the stat service
type ExternalProduct struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Quantity    int      `json:"quantity"`
	Category    string   `json:"category"`
	SKU         string   `json:"sku"`
	Images      []string `json:"images"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// Order represents an order in the system
type Order struct {
	ID        string         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    string         `json:"user_id" gorm:"type:uuid;not null"`
	Status    string         `json:"status" gorm:"not null;default:'pending'"` // pending, confirmed, shipped, delivered, cancelled
	Total     float64        `json:"total" gorm:"not null;default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships (no foreign key constraints since tables don't exist in this service)
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
}

// BeforeCreate hook to generate UUID if not set
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return nil
}

// OrderItem represents the many-to-many relationship between Order and Product
type OrderItem struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID   string    `json:"order_id" gorm:"type:uuid;not null"`
	ProductID string    `json:"product_id" gorm:"type:uuid;not null"`
	Quantity  int       `json:"quantity" gorm:"not null;min:1"`
	Price     float64   `json:"price" gorm:"not null"` // Price at the time of order
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships (no foreign key constraints)
	Order Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// BeforeCreate hook to generate UUID if not set
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == "" {
		oi.ID = uuid.New().String()
	}
	return nil
}

// OrderRequest represents a request to create a new order
type OrderRequest struct {
	Items []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// OrderItemRequest represents an item in the order request
type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,min=1,max=1000"`
}

// OrderResponse represents a response for order operations
type OrderResponse struct {
	ID        string              `json:"id"`
	UserID    string              `json:"user_id"`
	Status    string              `json:"status"`
	Total     float64             `json:"total"`
	Items     []OrderItemResponse `json:"items"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// OrderItemResponse represents an order item in the response
type OrderItemResponse struct {
	ID        string          `json:"id"`
	ProductID string          `json:"product_id"`
	Product   ExternalProduct `json:"product"`
	Quantity  int             `json:"quantity"`
	Price     float64         `json:"price"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
