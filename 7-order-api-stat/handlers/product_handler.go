package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"order-api-stat/models"
	"order-api-stat/service"
	"order-api-stat/validation"
)

// ProductHandler handles HTTP requests for product operations
type ProductHandler struct {
	productService *service.ProductService
	validator      *validation.Validator
}

// NewProductHandler creates a new product handler
func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: service.NewProductService(),
		validator:      validation.New(),
	}
}

// CreateProduct handles POST /products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Validation failed", errors)
		return
	}

	// Create product
	product, err := h.productService.CreateProduct(&req)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Convert to response format
	response := models.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Category:    product.Category,
		SKU:         product.SKU,
		Images:      product.Images,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	h.sendJSONResponse(w, http.StatusCreated, response)
}

// GetProduct handles GET /products/{id}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	// Get product
	product, err := h.productService.GetProduct(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendErrorResponse(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Convert to response format
	response := models.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Category:    product.Category,
		SKU:         product.SKU,
		Images:      product.Images,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	h.sendJSONResponse(w, http.StatusOK, response)
}

// ListProducts handles GET /products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	page := 1
	limit := 10
	category := ""

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if cat := r.URL.Query().Get("category"); cat != "" {
		category = cat
	}

	// List products
	response, err := h.productService.ListProducts(page, limit, category)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.sendJSONResponse(w, http.StatusOK, response)
}

// UpdateProduct handles PUT /products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Validation failed", errors)
		return
	}

	// Update product
	product, err := h.productService.UpdateProduct(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendErrorResponse(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			h.sendErrorResponse(w, http.StatusConflict, err.Error(), nil)
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Convert to response format
	response := models.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Category:    product.Category,
		SKU:         product.SKU,
		Images:      product.Images,
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	h.sendJSONResponse(w, http.StatusOK, response)
}

// DeleteProduct handles DELETE /products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	id, err := h.extractIDFromPath(r.URL.Path)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid product ID", nil)
		return
	}

	// Delete product
	if err := h.productService.DeleteProduct(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.sendErrorResponse(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper methods

// extractIDFromPath extracts ID from URL path like /products/123
func (h *ProductHandler) extractIDFromPath(path string) (uint, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 || parts[0] != "products" {
		return 0, fmt.Errorf("invalid path format")
	}

	id, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format")
	}

	return uint(id), nil
}

// sendJSONResponse sends a JSON response
func (h *ProductHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends an error response
func (h *ProductHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string, details map[string]string) {
	response := models.ErrorResponse{
		Error:   message,
		Details: details,
	}
	h.sendJSONResponse(w, statusCode, response)
}

// HandleProducts handles all product-related routes
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateProduct(w, r)
	case http.MethodGet:
		h.ListProducts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleProductByID handles product routes with ID parameter
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProduct(w, r)
	case http.MethodPut:
		h.UpdateProduct(w, r)
	case http.MethodDelete:
		h.DeleteProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
