package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"order-api-cart/middleware"
	"order-api-cart/models"
	"order-api-cart/service"
	"order-api-cart/validation"
)

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderService *service.OrderService
	validator    *validation.Validator
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(authServiceURL, productServiceURL string) *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(authServiceURL, productServiceURL),
		validator:    validation.New(),
	}
}

// CreateOrder handles POST /order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate request
	if !middleware.ValidateStruct(w, h.validator, &req) {
		return
	}

	// Create order
	order, err := h.orderService.CreateOrder(userID, &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create order: %v", err), http.StatusBadRequest)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Write response
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetOrderByID handles GET /order/{id}
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract order ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/order/")
	if path == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found", http.StatusUnauthorized)
		return
	}

	// Get order
	order, err := h.orderService.GetOrderByID(path)
	if err != nil {
		if err.Error() == "order not found" {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get order: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if the order belongs to the authenticated user
	if order.UserID != userID {
		http.Error(w, "You can only access your own orders", http.StatusForbidden)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Write response
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetMyOrders handles GET /my-orders
func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found", http.StatusUnauthorized)
		return
	}

	// Get pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get orders
	orders, err := h.orderService.GetOrdersByUserID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get orders: %v", err), http.StatusInternalServerError)
		return
	}

	// Simple pagination (in a real app, you'd implement this in the service layer)
	start := (page - 1) * limit
	end := start + limit

	if start >= len(orders) {
		orders = []models.OrderResponse{}
	} else {
		if end > len(orders) {
			end = len(orders)
		}
		orders = orders[start:end]
	}

	// Create response
	response := map[string]interface{}{
		"orders": orders,
		"page":   page,
		"limit":  limit,
		"total":  len(orders),
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Write response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status":  "ok",
		"service": "order-api-cart",
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Write response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
