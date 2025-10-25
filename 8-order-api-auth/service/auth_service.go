package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"order-api-auth/models"
	"order-api-auth/storage"
)

// AuthService interface for authorization
type AuthService interface {
	InitiateAuth(phone string) (*models.AuthResponse, error)
	VerifyCode(sessionID, code string) (*models.TokenResponse, error)
}

// AuthServiceImpl authorization service implementation
type AuthServiceImpl struct {
	userStorage    storage.UserStorage
	sessionStorage storage.SessionStorage
	smsService     SMSService
	jwtService     JWTService
}

// NewAuthService creates a new authorization service
func NewAuthService(
	userStorage storage.UserStorage,
	sessionStorage storage.SessionStorage,
	smsService SMSService,
	jwtService JWTService,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		userStorage:    userStorage,
		sessionStorage: sessionStorage,
		smsService:     smsService,
		jwtService:     jwtService,
	}
}

// InitiateAuth initiates authorization process
func (a *AuthServiceImpl) InitiateAuth(phone string) (*models.AuthResponse, error) {
	// Check if user exists
	user, err := a.userStorage.GetUserByPhone(phone)
	if err != nil {
		// If user not found, create new one
		user = &models.User{
			ID:        uuid.New().String(),
			Phone:     phone,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := a.userStorage.CreateUser(user); err != nil {
			return nil, err
		}
	}

	// Generate verification code
	code := a.smsService.GenerateCode()

	// Create session
	sessionID := uuid.New().String()
	session := &models.Session{
		ID:        sessionID,
		Phone:     phone,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute), // Session valid for 5 minutes
		CreatedAt: time.Now(),
		IsUsed:    false,
	}

	if err := a.sessionStorage.CreateSession(session); err != nil {
		return nil, err
	}

	// Send SMS with code
	if err := a.smsService.SendCode(phone, code); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		SessionID: sessionID,
	}, nil
}

// VerifyCode verifies confirmation code
func (a *AuthServiceImpl) VerifyCode(sessionID, code string) (*models.TokenResponse, error) {
	// Get session
	session, err := a.sessionStorage.GetSession(sessionID)
	if err != nil {
		return nil, errors.New("invalid session")
	}

	// Check if session expired
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}

	// Check if session already used
	if session.IsUsed {
		return nil, errors.New("session already used")
	}

	// Verify code
	if session.Code != code {
		return nil, errors.New("invalid code")
	}

	// Get user
	user, err := a.userStorage.GetUserByPhone(session.Phone)
	if err != nil {
		return nil, err
	}

	// Mark session as used
	if err := a.sessionStorage.MarkSessionAsUsed(sessionID); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := a.jwtService.GenerateToken(user.ID, user.Phone)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		Token: token,
	}, nil
}
