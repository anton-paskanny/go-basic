package models

import (
	"time"
)

// Purchase represents a purchase transaction
type Purchase struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"` // pending, completed, cancelled
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PurchaseRequest represents a purchase request
type PurchaseRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// PurchaseResponse represents a purchase response
type PurchaseResponse struct {
	PurchaseID string  `json:"purchase_id"`
	Total      float64 `json:"total"`
	Status     string  `json:"status"`
	Message    string  `json:"message"`
}
