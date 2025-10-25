package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"order-api-auth/models"
	"order-api-auth/utils"
)

// PurchaseHandler handler for purchase operations
type PurchaseHandler struct {
	products map[string]*models.Product
}

// NewPurchaseHandler creates a new purchase handler
func NewPurchaseHandler() *PurchaseHandler {
	// Initialize with some sample products
	products := map[string]*models.Product{
		"1": {
			ID:          "1",
			Name:        "Laptop",
			Description: "High-performance laptop",
			Price:       999.99,
			Stock:       10,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		"2": {
			ID:          "2",
			Name:        "Smartphone",
			Description: "Latest smartphone model",
			Price:       699.99,
			Stock:       25,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		"3": {
			ID:          "3",
			Name:        "Headphones",
			Description: "Wireless noise-canceling headphones",
			Price:       199.99,
			Stock:       50,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	return &PurchaseHandler{
		products: products,
	}
}

// GetProducts returns list of available products
func (h *PurchaseHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Convert map to slice
	productList := make([]*models.Product, 0, len(h.products))
	for _, product := range h.products {
		productList = append(productList, product)
	}

	utils.WriteJSONResponse(w, http.StatusOK, productList)
}

// PurchaseProduct handles product purchase
func (h *PurchaseHandler) PurchaseProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.ProductID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	if req.Quantity <= 0 {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	// Get product
	product, exists := h.products[req.ProductID]
	if !exists {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// Check stock
	if product.Stock < req.Quantity {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Insufficient stock")
		return
	}

	// Get user ID from context (set by middleware)
	userID := r.Context().Value("user_id").(string)

	// Calculate total
	total := product.Price * float64(req.Quantity)

	// Create purchase
	purchaseID := uuid.New().String()
	purchase := &models.Purchase{
		ID:        purchaseID,
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Total:     total,
		Status:    "completed",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Update stock
	product.Stock -= req.Quantity
	product.UpdatedAt = time.Now()

	// Return response
	response := models.PurchaseResponse{
		PurchaseID: purchase.ID,
		Total:      purchase.Total,
		Status:     purchase.Status,
		Message:    "Purchase completed successfully",
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}
