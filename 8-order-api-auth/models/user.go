package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	Phone     string    `json:"phone" gorm:"size:15;uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Session represents an authorization session
type Session struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	Phone     string    `json:"phone" gorm:"size:15;not null"`
	Code      string    `json:"code" gorm:"size:4;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
}

// AuthRequest represents an authorization request
type AuthRequest struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
}

// VerifyCodeRequest represents a code verification request
type VerifyCodeRequest struct {
	SessionID string `json:"sessionId" validate:"required"`
	Code      string `json:"code" validate:"required,len=4"`
}

// AuthResponse represents a response with sessionId
type AuthResponse struct {
	SessionID string `json:"sessionId"`
}

// TokenResponse represents a response with JWT token
type TokenResponse struct {
	Token string `json:"token"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
