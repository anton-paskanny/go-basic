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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found", http.StatusUnauthorized)
		return
	}

	authToken := r.Header.Get("Authorization")

	var req models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if !middleware.ValidateStruct(w, h.validator, &req) {
		return
	}

	order, err := h.orderService.CreateOrder(userID, &req, authToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create order: %v", err), http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

// GetOrderByID handles GET /order/{id}
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/order/")
	if path == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found", http.StatusUnauthorized)
		return
	}

	authToken := r.Header.Get("Authorization")

	order, err := h.orderService.GetOrderByID(path, authToken)
	if err != nil {
		if err.Error() == "order not found" {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get order: %v", err), http.StatusInternalServerError)
		return
	}

	if order.UserID != userID {
		http.Error(w, "You can only access your own orders", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetMyOrders handles GET /my-orders
func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found", http.StatusUnauthorized)
		return
	}

	authToken := r.Header.Get("Authorization")

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

	orders, total, err := h.orderService.GetOrdersByUserID(userID, page, limit, authToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get orders: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"orders": orders,
		"page":   page,
		"limit":  limit,
		"total":  total,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status":  "ok",
		"service": "order-api-cart",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
