package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"order-api-cart/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestData holds test data for e2e tests
type TestData struct {
	UserID     string
	ProductIDs []string
	AuthToken  string
}

// MockAuthService mocks the auth service for testing
type MockAuthService struct {
	users map[string]*models.ExternalUser
}

// MockProductService mocks the product service for testing
type MockProductService struct {
	products map[string]*models.ExternalProduct
}

// NewMockAuthService creates a new mock auth service
func NewMockAuthService() *MockAuthService {
	return &MockAuthService{
		users: make(map[string]*models.ExternalUser),
	}
}

// NewMockProductService creates a new mock product service
func NewMockProductService() *MockProductService {
	return &MockProductService{
		products: make(map[string]*models.ExternalProduct),
	}
}

// CreateTestUser creates a test user in the mock auth service
func (m *MockAuthService) CreateTestUser(t *testing.T) *models.ExternalUser {
	userID := uuid.New().String()
	user := &models.ExternalUser{
		ID:        userID,
		Phone:     "+1234567890",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.users[userID] = user
	return user
}

// GetUserByID returns user by ID
func (m *MockAuthService) GetUserByID(userID string) (*models.ExternalUser, error) {
	user, exists := m.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// CreateTestProduct creates a test product in the mock product service
func (m *MockProductService) CreateTestProduct(t *testing.T, name string, price float64, quantity int) *models.ExternalProduct {
	productID := uuid.New().String()
	product := &models.ExternalProduct{
		ID:          productID,
		Name:        name,
		Description: fmt.Sprintf("Test product: %s", name),
		Price:       price,
		Quantity:    quantity,
		Category:    "test",
		SKU:         fmt.Sprintf("TEST-%s", productID[:8]),
		Images:      []string{},
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	m.products[productID] = product
	return product
}

// GetProductByID returns product by ID
func (m *MockProductService) GetProductByID(productID string) (*models.ExternalProduct, error) {
	product, exists := m.products[productID]
	if !exists {
		return nil, fmt.Errorf("product not found")
	}
	return product, nil
}

// UpdateProductQuantity updates product quantity
func (m *MockProductService) UpdateProductQuantity(productID string, quantityChange int) error {
	product, exists := m.products[productID]
	if !exists {
		return fmt.Errorf("product not found")
	}
	product.Quantity += quantityChange
	if product.Quantity < 0 {
		return fmt.Errorf("insufficient quantity")
	}
	return nil
}

// StartMockAuthService starts a mock auth service server and registers a
// cleanup that shuts it down when the test ends, freeing the port for the
// next test.
func StartMockAuthService(t *testing.T, port string) *MockAuthService {
	mock := NewMockAuthService()

	mux := http.NewServeMux()
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.URL.Path[len("/users/"):]
		if userID == "" {
			http.Error(w, "User ID required", http.StatusBadRequest)
			return
		}

		user, err := mock.GetUserByID(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	server := &http.Server{Addr: ":" + port, Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("Mock auth service error: %v", err)
		}
	}()
	t.Cleanup(func() { server.Shutdown(context.Background()) })

	time.Sleep(100 * time.Millisecond)
	return mock
}

// StartMockProductService starts a mock product service server and registers a
// cleanup that shuts it down when the test ends, freeing the port for the
// next test.
func StartMockProductService(t *testing.T, port string) *MockProductService {
	mock := NewMockProductService()

	mux := http.NewServeMux()
	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		// Strip the "/products/" prefix, then split off any sub-path (e.g. "/quantity")
		rest := r.URL.Path[len("/products/"):]
		parts := strings.SplitN(rest, "/", 2)
		productID := parts[0]
		subpath := ""
		if len(parts) == 2 {
			subpath = parts[1]
		}

		if productID == "" {
			http.Error(w, "Product ID required", http.StatusBadRequest)
			return
		}

		switch {
		case r.Method == http.MethodGet && subpath == "":
			product, err := mock.GetProductByID(productID)
			if err != nil {
				http.Error(w, "Product not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)

		case r.Method == http.MethodPatch && subpath == "quantity":
			var req struct {
				Change int `json:"change"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			if err := mock.UpdateProductQuantity(productID, req.Change); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	server := &http.Server{Addr: ":" + port, Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("Mock product service error: %v", err)
		}
	}()
	t.Cleanup(func() { server.Shutdown(context.Background()) })

	time.Sleep(100 * time.Millisecond)
	return mock
}

// GenerateTestJWT generates a signed JWT token for testing.
// It reads JWT_SECRET from the environment (falls back to "test-secret-key")
// so the token is accepted by the auth middleware.
func GenerateTestJWT(userID string) string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test-secret-key"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(fmt.Sprintf("failed to generate test JWT: %v", err))
	}
	return tokenString
}

// CreateTestOrderRequest creates a test order request
func CreateTestOrderRequest(productIDs []string, quantities []int) *models.OrderRequest {
	if len(productIDs) != len(quantities) {
		panic("productIDs and quantities must have the same length")
	}

	var items []models.OrderItemRequest
	for i, productID := range productIDs {
		items = append(items, models.OrderItemRequest{
			ProductID: productID,
			Quantity:  quantities[i],
		})
	}

	return &models.OrderRequest{
		Items: items,
	}
}

// MakeOrderRequest makes an HTTP request to create an order
func MakeOrderRequest(t *testing.T, baseURL string, authToken string, orderReq *models.OrderRequest) (*http.Response, error) {
	jsonData, err := json.Marshal(orderReq)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", baseURL+"/api/v1/order", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

// AssertOrderResponse validates order response
func AssertOrderResponse(t *testing.T, resp *http.Response, expectedUserID string, expectedItemCount int) {
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var orderResp models.OrderResponse
	err := json.NewDecoder(resp.Body).Decode(&orderResp)
	assert.NoError(t, err)

	assert.NotEmpty(t, orderResp.ID)
	assert.Equal(t, expectedUserID, orderResp.UserID)
	assert.Equal(t, "pending", orderResp.Status)
	assert.Greater(t, orderResp.Total, 0.0)
	assert.Len(t, orderResp.Items, expectedItemCount)

	for _, item := range orderResp.Items {
		assert.NotEmpty(t, item.ID)
		assert.NotEmpty(t, item.ProductID)
		assert.Greater(t, item.Quantity, 0)
		assert.Greater(t, item.Price, 0.0)
	}
}
