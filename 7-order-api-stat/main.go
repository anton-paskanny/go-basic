package main

import (
	"fmt"
	"log"
	"net/http"

	"order-api-stat/config"
	"order-api-stat/database"
	"order-api-stat/handlers"
	"order-api-stat/utils"
	"order-api-stat/validation"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize validator
	validator := validation.New()
	_ = validator // Use validator in handlers

	// Initialize handlers
	productHandler := handlers.NewProductHandler()
	healthHandler := handlers.NewHealthHandler()

	// Setup routes
	mux := http.NewServeMux()

	// Product routes
	mux.HandleFunc("/products", productHandler.HandleProducts)
	mux.HandleFunc("/products/", productHandler.HandleProductByID)

	// Health check endpoint
	mux.HandleFunc("/health", healthHandler.HandleHealth)

	// Add CORS middleware
	handler := utils.CORSMiddleware(mux)

	// Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on port %d", cfg.Server.Port)
	log.Printf("Available endpoints:")
	log.Printf("  POST   /products     - Create a new product")
	log.Printf("  GET    /products      - List products (with pagination)")
	log.Printf("  GET    /products/{id} - Get a specific product")
	log.Printf("  PUT    /products/{id} - Update a product")
	log.Printf("  DELETE /products/{id} - Delete a product")
	log.Printf("  GET    /health        - Health check")

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
