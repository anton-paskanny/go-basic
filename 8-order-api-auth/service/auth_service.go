package service

import (
	"crypto/subtle"
	"errors"
	"time"

	"github.com/google/uuid"

	"order-api-auth/models"
	"order-api-auth/storage"
)

const maxCodeAttempts = 5

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
	// Check if user exists; only create one when record is genuinely absent
	user, err := a.userStorage.GetUserByPhone(phone)
	if err != nil {
		if err.Error() != "user not found" {
			return nil, err
		}
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

	// Generate verification code and send SMS before persisting the session,
	// so a failed SMS does not leave an orphaned session in the DB.
	code := a.smsService.GenerateCode()
	if err := a.smsService.SendCode(phone, code); err != nil {
		return nil, err
	}

	sessionID := uuid.New().String()
	session := &models.Session{
		ID:        sessionID,
		Phone:     phone,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
		IsUsed:    false,
	}
	if err := a.sessionStorage.CreateSession(session); err != nil {
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

	// Enforce brute-force limit before checking the code
	if session.Attempts >= maxCodeAttempts {
		return nil, errors.New("too many incorrect attempts, session is locked")
	}

	// Constant-time comparison prevents timing attacks
	if subtle.ConstantTimeCompare([]byte(session.Code), []byte(code)) != 1 {
		_ = a.sessionStorage.IncrementAttempts(sessionID)
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
