package models

// CreateProductRequest represents the request payload for creating a product
type CreateProductRequest struct {
	Name        string   `json:"name" validate:"required,min=3,max=255"`
	Description string   `json:"description" validate:"omitempty,max=1000"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	Quantity    int      `json:"quantity" validate:"min=0"`
	Category    string   `json:"category" validate:"omitempty,max=100"`
	SKU         string   `json:"sku" validate:"required,min=3,max=50"`
	Images      []string `json:"images"`
}

// UpdateProductRequest represents the request payload for updating a product
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=1000"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Quantity    *int     `json:"quantity,omitempty" validate:"omitempty,min=0"`
	Category    *string  `json:"category,omitempty" validate:"omitempty,max=100"`
	SKU         *string  `json:"sku,omitempty" validate:"omitempty,min=3,max=50"`
	Images      []string `json:"images,omitempty"`
}

// ProductResponse represents the response payload for product operations
type ProductResponse struct {
	ID          uint     `json:"id"`
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

// ProductListResponse represents the response for listing products
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}
