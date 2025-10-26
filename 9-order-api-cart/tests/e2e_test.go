package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"order-api-cart/config"
	"order-api-cart/database"
	"order-api-cart/handlers"
	"order-api-cart/middleware"
	"order-api-cart/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Setup test environment
	os.Setenv("DB_NAME", "order_cart_test_db")
	os.Setenv("SERVER_PORT", "8083")
	os.Setenv("AUTH_SERVICE_URL", "http://localhost:8084")
	os.Setenv("PRODUCT_SERVICE_URL", "http://localhost:8085")
	os.Setenv("JWT_SECRET", "test-secret-key")

	// Run tests
	code := m.Run()

	// Cleanup
	os.Exit(code)
}

func TestCreateOrderE2E(t *testing.T) {
	// Setup test database
	cfg := LoadTestConfig()
	defer CleanupTestDB(t)

	// Connect to test database
	err := database.Connect(cfg.Config)
	require.NoError(t, err)

	// Run migrations
	err = database.Migrate()
	require.NoError(t, err)

	// Start mock services
	mockAuth := StartMockAuthService(t, "8084")
	mockProduct := StartMockProductService(t, "8085")

	// Create test data
	testUser := mockAuth.CreateTestUser(t)
	testProduct1 := mockProduct.CreateTestProduct(t, "Test Product 1", 10.50, 100)
	testProduct2 := mockProduct.CreateTestProduct(t, "Test Product 2", 25.00, 50)

	// Generate test JWT token
	authToken := GenerateTestJWT(testUser.ID)

	// Start the main application server
	server := startTestServer(t, cfg.Config)
	defer server.Shutdown(context.Background())

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Test cases
	t.Run("CreateOrder_Success", func(t *testing.T) {
		// Create order request
		orderReq := CreateTestOrderRequest(
			[]string{testProduct1.ID, testProduct2.ID},
			[]int{2, 1},
		)

		// Make request
		resp, err := MakeOrderRequest(t, "http://localhost:8083", authToken, orderReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response
		AssertOrderResponse(t, resp, testUser.ID, 2)

		// Verify order was created in database
		var order models.Order
		db := database.GetDB()
		err = db.Preload("OrderItems").First(&order, "user_id = ?", testUser.ID).Error
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, order.UserID)
		assert.Equal(t, "pending", order.Status)
		assert.Equal(t, 46.0, order.Total) // (10.50 * 2) + (25.00 * 1)
		assert.Len(t, order.OrderItems, 2)
	})

	t.Run("CreateOrder_InvalidProduct", func(t *testing.T) {
		// Create order request with non-existent product
		orderReq := CreateTestOrderRequest(
			[]string{"non-existent-product-id"},
			[]int{1},
		)

		// Make request
		resp, err := MakeOrderRequest(t, "http://localhost:8083", authToken, orderReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return error
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("CreateOrder_InsufficientQuantity", func(t *testing.T) {
		// Create order request with quantity exceeding available stock
		orderReq := CreateTestOrderRequest(
			[]string{testProduct1.ID},
			[]int{1000}, // More than available (100)
		)

		// Make request
		resp, err := MakeOrderRequest(t, "http://localhost:8083", authToken, orderReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return error
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("CreateOrder_Unauthorized", func(t *testing.T) {
		// Create order request without auth token
		orderReq := CreateTestOrderRequest(
			[]string{testProduct1.ID},
			[]int{1},
		)

		// Make request without auth token
		resp, err := MakeOrderRequest(t, "http://localhost:8083", "", orderReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("CreateOrder_InvalidRequest", func(t *testing.T) {
		// Create invalid order request (empty items)
		orderReq := &models.OrderRequest{
			Items: []models.OrderItemRequest{},
		}

		// Make request
		resp, err := MakeOrderRequest(t, "http://localhost:8083", authToken, orderReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return validation error
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGetOrderByIDE2E(t *testing.T) {
	// Setup test database
	cfg := LoadTestConfig()
	defer CleanupTestDB(t)

	// Connect to test database
	err := database.Connect(cfg.Config)
	require.NoError(t, err)

	// Run migrations
	err = database.Migrate()
	require.NoError(t, err)

	// Start mock services
	mockAuth := StartMockAuthService(t, "8084")
	mockProduct := StartMockProductService(t, "8085")

	// Create test data
	testUser := mockAuth.CreateTestUser(t)
	testProduct := mockProduct.CreateTestProduct(t, "Test Product", 15.00, 100)

	// Generate test JWT token
	authToken := GenerateTestJWT(testUser.ID)

	// Start the main application server
	server := startTestServer(t, cfg.Config)
	defer server.Shutdown(context.Background())

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Create an order first
	orderReq := CreateTestOrderRequest(
		[]string{testProduct.ID},
		[]int{3},
	)

	resp, err := MakeOrderRequest(t, "http://localhost:8083", authToken, orderReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	var orderResp models.OrderResponse
	err = json.NewDecoder(resp.Body).Decode(&orderResp)
	require.NoError(t, err)

	// Test getting order by ID
	t.Run("GetOrderByID_Success", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8083/api/v1/order/%s", orderResp.ID), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authToken)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var retrievedOrder models.OrderResponse
		err = json.NewDecoder(resp.Body).Decode(&retrievedOrder)
		require.NoError(t, err)

		assert.Equal(t, orderResp.ID, retrievedOrder.ID)
		assert.Equal(t, testUser.ID, retrievedOrder.UserID)
		assert.Equal(t, "pending", retrievedOrder.Status)
		assert.Equal(t, 45.0, retrievedOrder.Total) // 15.00 * 3
	})

	t.Run("GetOrderByID_NotFound", func(t *testing.T) {
		nonExistentID := "non-existent-order-id"
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8083/api/v1/order/%s", nonExistentID), nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authToken)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestGetMyOrdersE2E(t *testing.T) {
	// Setup test database
	cfg := LoadTestConfig()
	defer CleanupTestDB(t)

	// Connect to test database
	err := database.Connect(cfg.Config)
	require.NoError(t, err)

	// Run migrations
	err = database.Migrate()
	require.NoError(t, err)

	// Start mock services
	mockAuth := StartMockAuthService(t, "8084")
	mockProduct := StartMockProductService(t, "8085")

	// Create test data
	testUser := mockAuth.CreateTestUser(t)
	testProduct1 := mockProduct.CreateTestProduct(t, "Test Product 1", 10.00, 100)
	testProduct2 := mockProduct.CreateTestProduct(t, "Test Product 2", 20.00, 100)

	// Generate test JWT token
	authToken := GenerateTestJWT(testUser.ID)

	// Start the main application server
	server := startTestServer(t, cfg.Config)
	defer server.Shutdown(context.Background())

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Create multiple orders
	for i := 0; i < 3; i++ {
		orderReq := CreateTestOrderRequest(
			[]string{testProduct1.ID, testProduct2.ID},
			[]int{i + 1, i + 1},
		)

		resp, err := MakeOrderRequest(t, "http://localhost:8083", authToken, orderReq)
		require.NoError(t, err)
		resp.Body.Close()
	}

	// Test getting user's orders
	t.Run("GetMyOrders_Success", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://localhost:8083/api/v1/my-orders", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authToken)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		orders, ok := response["orders"].([]interface{})
		require.True(t, ok)
		assert.Len(t, orders, 3)

		page, ok := response["page"].(float64)
		require.True(t, ok)
		assert.Equal(t, float64(1), page)

		limit, ok := response["limit"].(float64)
		require.True(t, ok)
		assert.Equal(t, float64(10), limit)
	})
}

// startTestServer starts the test server
func startTestServer(t *testing.T, cfg *config.Config) *http.Server {
	// Create handlers
	orderHandler := handlers.NewOrderHandler(cfg.Services.AuthServiceURL, cfg.Services.ProductServiceURL)

	// Create mux
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", handlers.HealthCheck)

	// API routes
	mux.HandleFunc("/api/v1/order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			orderHandler.CreateOrder(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Order by ID endpoint
	mux.HandleFunc("/api/v1/order/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			orderHandler.GetOrderByID(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// My orders endpoint
	mux.HandleFunc("/api/v1/my-orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			orderHandler.GetMyOrders(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Apply middleware
	handler := middleware.CORSMiddleware()(mux)

	// Apply auth middleware to protected routes
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is for a protected endpoint
		if strings.HasPrefix(r.URL.Path, "/api/v1/order") || strings.HasPrefix(r.URL.Path, "/api/v1/my-orders") {
			// Apply auth middleware
			authMiddleware := middleware.AuthMiddleware(cfg)
			authMiddleware(handler).ServeHTTP(w, r)
		} else {
			// Serve unprotected routes directly
			handler.ServeHTTP(w, r)
		}
	})

	// Create server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: protectedHandler,
	}

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	return server
}
