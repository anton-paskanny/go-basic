package handlers

import (
	"encoding/json"
	"net/http"

	"order-api-auth/models"
	"order-api-auth/service"
	"order-api-auth/utils"
)

// AuthHandler handler for authorization
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new authorization handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// InitiateAuth initiates authorization process
func (h *AuthHandler) InitiateAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate phone
	if err := utils.ValidatePhone(req.Phone); err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.WriteValidationErrorResponse(w, validationErr)
			return
		}
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Initiate authorization
	response, err := h.authService.InitiateAuth(req.Phone)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// VerifyCode verifies confirmation code
func (h *AuthHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate session ID
	if err := utils.ValidateSessionID(req.SessionID); err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.WriteValidationErrorResponse(w, validationErr)
			return
		}
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate code
	if err := utils.ValidateCode(req.Code); err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.WriteValidationErrorResponse(w, validationErr)
			return
		}
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Verify code
	response, err := h.authService.VerifyCode(req.SessionID, req.Code)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}
