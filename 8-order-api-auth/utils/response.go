package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSONResponse writes JSON response
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteErrorResponse writes JSON response with error
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]string{
		"error": message,
	}
	WriteJSONResponse(w, statusCode, response)
}

// WriteValidationErrorResponse writes JSON response with validation error
func WriteValidationErrorResponse(w http.ResponseWriter, err *ValidationError) {
	response := map[string]interface{}{
		"error":   "Validation failed",
		"field":   err.Field,
		"message": err.Message,
	}
	WriteJSONResponse(w, http.StatusBadRequest, response)
}
