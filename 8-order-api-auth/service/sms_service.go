package service

import (
	"fmt"
	"math/rand"
	"time"
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
	// In real project this would be SMS provider API call
	fmt.Printf("🔐 SMS sent to %s with code: %s\n", phone, code)
	return nil
}

// GenerateCode generates a 4-digit verification code
func (s *MockSMSService) GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}
