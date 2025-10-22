package main

import (
	"fmt"
	"net/http"

	"order-api-stat/config"
	"order-api-stat/database"
	"order-api-stat/handlers"
	"order-api-stat/utils"
	"order-api-stat/validation"

	"github.com/sirupsen/logrus"
)

func main() {
	// Configure logrus to use JSON format
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
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

	// Add middleware chain: logging -> CORS
	handler := utils.LoggingMiddleware(utils.CORSMiddleware(mux))

	// Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.WithField("port", cfg.Server.Port).Info("Server starting")
	logrus.Info("Available endpoints:")
	logrus.Info("  POST   /products     - Create a new product")
	logrus.Info("  GET    /products      - List products (with pagination)")
	logrus.Info("  GET    /products/{id} - Get a specific product")
	logrus.Info("  PUT    /products/{id} - Update a product")
	logrus.Info("  DELETE /products/{id} - Delete a product")
	logrus.Info("  GET    /health        - Health check")

	if err := http.ListenAndServe(addr, handler); err != nil {
		logrus.Fatalf("Server failed to start: %v", err)
	}
}
