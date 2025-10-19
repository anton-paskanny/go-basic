package verify

import (
	"encoding/json"
	"fmt"
	"net/http"

	"validation-api/config"
)

// Handler handles email verification requests
type Handler struct {
	service *Service
}

// NewHandler creates a new verification handler
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		service: NewService(cfg),
	}
}

// SendRequest represents a request to send a verification email
type SendRequest struct {
	Email string `json:"email"`
}

// SendResponse represents a response to a send verification request
type SendResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SendHandler handles requests to send verification emails
func (h *Handler) SendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	_, err := h.service.SendVerificationEmail(req.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send verification email: %v", err), http.StatusInternalServerError)
		return
	}

	resp := SendResponse{
		Success: true,
		Message: "Verification email sent",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// VerifyHandler handles requests to verify email addresses
func (h *Handler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract hash from URL path
	hash := r.URL.Path[len("/verify/"):]
	if hash == "" {
		http.Error(w, "Invalid verification link", http.StatusBadRequest)
		return
	}

	email, verified := h.service.VerifyEmail(hash)
	if !verified {
		http.Error(w, "Invalid or expired verification link", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Email Verified</h1><p>Your email %s has been successfully verified.</p>", email)
}

// RegisterRoutes registers the verification routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/send", h.SendHandler)
	mux.HandleFunc("/verify/", h.VerifyHandler)
}
