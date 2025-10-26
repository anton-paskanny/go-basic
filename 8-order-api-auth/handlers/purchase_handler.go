package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"order-api-auth/models"
	"order-api-auth/utils"
)

// ExternalProduct represents a product from the product service
type ExternalProduct struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Category    string  `json:"category"`
	SKU         string  `json:"sku"`
}

// PurchaseHandler handler for purchase operations
type PurchaseHandler struct {
	productServiceURL string
}

// NewPurchaseHandler creates a new purchase handler
func NewPurchaseHandler(productServiceURL string) *PurchaseHandler {
	return &PurchaseHandler{
		productServiceURL: productServiceURL,
	}
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

	// Get product from external service
	product, err := h.getProductFromService(req.ProductID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// Check stock
	if product.Quantity < req.Quantity {
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

	// Return response
	response := models.PurchaseResponse{
		PurchaseID: purchase.ID,
		Total:      purchase.Total,
		Status:     purchase.Status,
		Message:    "Purchase completed successfully",
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// getProductFromService fetches product information from the product service
func (h *PurchaseHandler) getProductFromService(productID string) (*ExternalProduct, error) {
	url := fmt.Sprintf("%s/products/%s", h.productServiceURL, productID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product not found")
	}

	var product ExternalProduct
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}
