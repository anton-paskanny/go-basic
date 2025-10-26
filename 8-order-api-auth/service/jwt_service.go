package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService interface for working with JWT tokens
type JWTService interface {
	GenerateToken(userID, phone string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// Claims represents claims in JWT token
type Claims struct {
	UserID string `json:"user_id"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

// JWTServiceImpl JWT service implementation
type JWTServiceImpl struct {
	secretKey []byte
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string) *JWTServiceImpl {
	return &JWTServiceImpl{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken generates JWT token for user
func (j *JWTServiceImpl) GenerateToken(userID, phone string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Phone:  phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken validates JWT token and returns claims
func (j *JWTServiceImpl) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
