package main

import (
	"log"
	"net/http"
	"strings"

	"order-api-cart/config"
	"order-api-cart/database"
	"order-api-cart/handlers"
	"order-api-cart/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Create handlers with service URLs from config
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

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, protectedHandler); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
