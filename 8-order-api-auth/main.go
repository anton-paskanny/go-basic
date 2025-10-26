package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"order-api-auth/config"
	"order-api-auth/database"
	"order-api-auth/handlers"
	"order-api-auth/middleware"
	"order-api-auth/service"
	"order-api-auth/storage"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize storage
	storage := storage.NewPostgreSQLStorage(db)

	// Initialize services
	smsService := service.NewMockSMSService()
	jwtService := service.NewJWTService(cfg.JWTSecret)
	authService := service.NewAuthService(storage, storage, smsService, jwtService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	purchaseHandler := handlers.NewPurchaseHandler(cfg.ProductServiceURL)

	// Initialize middleware
	corsMiddleware := middleware.NewCORSMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Create router
	router := mux.NewRouter()

	// Apply CORS middleware to all routes
	router.Use(corsMiddleware.CORS)

	// Auth routes (public)
	router.HandleFunc("/auth/initiate", authHandler.InitiateAuth).Methods("POST")
	router.HandleFunc("/auth/verify", authHandler.VerifyCode).Methods("POST")

	// Protected routes (require JWT)
	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(authMiddleware.RequireAuth)
	protectedRouter.HandleFunc("/purchase", purchaseHandler.PurchaseProduct).Methods("POST")

	// Start expired sessions cleanup in background
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			storage.CleanupExpiredSessions()
		}
	}()

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
