package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"order-api-cart/models"
	"order-api-cart/validation"
)

// ValidationMiddleware creates a middleware for validating request bodies
func ValidationMiddleware(validator *validation.Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only validate POST, PUT, PATCH requests
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				// Read the request body
				body := r.Body
				defer body.Close()

				// Parse JSON into a generic map to check basic structure
				var requestData map[string]interface{}
				if err := json.NewDecoder(body).Decode(&requestData); err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)

					errorResponse := models.ErrorResponse{
						Error:   "Invalid JSON",
						Message: "Request body must be valid JSON",
					}

					json.NewEncoder(w).Encode(errorResponse)
					return
				}

				// Re-encode the data back to JSON for further processing
				jsonData, err := json.Marshal(requestData)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)

					errorResponse := models.ErrorResponse{
						Error:   "Invalid request data",
						Message: "Unable to process request data",
					}

					json.NewEncoder(w).Encode(errorResponse)
					return
				}

				// Store the parsed data in the request context for handlers to use
				ctx := r.Context()
				ctx = context.WithValue(ctx, "requestData", jsonData)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ValidateStruct validates a struct and returns appropriate HTTP response
func ValidateStruct(w http.ResponseWriter, validator *validation.Validator, data interface{}) bool {
	if err := validator.ValidateStruct(data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		errorResponse := models.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(errorResponse)
		return false
	}
	return true
}
