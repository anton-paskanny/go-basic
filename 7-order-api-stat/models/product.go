package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Product represents a product in the system
type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:255;not null" validate:"required,min=3,max=255"`
	Description string         `json:"description" gorm:"type:text" validate:"omitempty,max=1000"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null" validate:"required,gt=0"`
	Quantity    int            `json:"quantity" gorm:"not null;default:0" validate:"min=0"`
	Category    string         `json:"category" gorm:"size:100" validate:"omitempty,max=100"`
	SKU         string         `json:"sku" gorm:"size:50;uniqueIndex" validate:"required,min=3,max=50"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}
