package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"order-api-auth/config"
	"order-api-auth/handlers"
	"order-api-auth/middleware"
	"order-api-auth/service"
	"order-api-auth/storage"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize storage
	storage := storage.NewInMemoryStorage()

	// Initialize services
	smsService := service.NewMockSMSService()
	jwtService := service.NewJWTService(cfg.JWTSecret)
	authService := service.NewAuthService(storage, storage, smsService, jwtService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	purchaseHandler := handlers.NewPurchaseHandler()

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

	// Product routes (public)
	router.HandleFunc("/products", purchaseHandler.GetProducts).Methods("GET")

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
