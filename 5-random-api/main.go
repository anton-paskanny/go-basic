package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

// RandomResponse represents the JSON response structure
type RandomResponse struct {
	Number int `json:"number"`
}

// randomHandler generates a random number between 1-6 and returns it as JSON
func randomHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random number from 1 to 6
	randomNumber := rand.Intn(6) + 1

	// Create response
	response := RandomResponse{Number: randomNumber}

	// Set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Send JSON response
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Register handler
	http.HandleFunc("/random", randomHandler)

	// Start server
	fmt.Println("Server started at http://localhost:8080")
	fmt.Println("To get a random number, send a GET request to http://localhost:8080/random")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
