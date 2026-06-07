package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// SMSService interface for sending SMS
type SMSService interface {
	SendCode(phone, code string) error
	GenerateCode() string
}

// MockSMSService mock service for sending SMS (for testing)
type MockSMSService struct{}

// NewMockSMSService creates a new mock SMS service
func NewMockSMSService() *MockSMSService {
	return &MockSMSService{}
}

// SendCode sends SMS with code (in real project this would be SMS provider integration)
func (s *MockSMSService) SendCode(phone, code string) error {
	// In real project this would be an SMS provider API call.
	// The code is intentionally not logged to avoid leaking OTPs.
	fmt.Printf("SMS sent to %s\n", phone)
	return nil
}

// GenerateCode generates a cryptographically random 4-digit verification code.
func (s *MockSMSService) GenerateCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		panic("failed to generate secure random code: " + err.Error())
	}
	return fmt.Sprintf("%04d", n.Int64())
}
