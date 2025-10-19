package main

import (
	"fmt"
	"log"
	"net/http"

	"validation-api/config"
	"validation-api/verify"
)

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Create a new router
	mux := http.NewServeMux()

	// Initialize verification handler
	verifyHandler := verify.NewHandler(cfg)
	verifyHandler.RegisterRoutes(mux)

	// Start the server
	serverAddr := cfg.Server.Address
	fmt.Printf("Starting server on %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, mux))
}
