package verify

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"sync"
	"time"

	"validation-api/config"

	"github.com/jordan-wright/email"
)

// Service handles email verification logic
type Service struct {
	cfg           *config.Config
	verifications map[string]verificationData
	mu            sync.RWMutex
}

type verificationData struct {
	email     string
	createdAt time.Time
	verified  bool
}

// NewService creates a new verification service
func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg:           cfg,
		verifications: make(map[string]verificationData),
	}
}

// SendVerificationEmail sends a verification email and returns the verification hash
func (s *Service) SendVerificationEmail(emailAddr string) (string, error) {
	// Validate email format
	if err := validateEmail(emailAddr); err != nil {
		return "", err
	}

	// Generate random verification hash
	hashBytes := make([]byte, 16)
	if _, err := rand.Read(hashBytes); err != nil {
		return "", fmt.Errorf("failed to generate verification hash: %w", err)
	}
	hash := hex.EncodeToString(hashBytes)

	// Store verification data
	s.mu.Lock()
	s.verifications[hash] = verificationData{
		email:     emailAddr,
		createdAt: time.Now(),
		verified:  false,
	}
	s.mu.Unlock()

	// Create verification link
	verificationLink := fmt.Sprintf("http://localhost%s/verify/%s", s.cfg.Server.Address, hash)

	// Create email
	e := email.NewEmail()
	e.From = fmt.Sprintf("Verification <no-reply@%s>", s.cfg.Email.Host)
	e.To = []string{emailAddr}
	e.Subject = "Email Verification"
	e.HTML = []byte(fmt.Sprintf(`
		<h1>Email Verification</h1>
		<p>Please click the link below to verify your email address:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>This link will expire in 24 hours.</p>
	`, verificationLink))

	// Send email
	auth := smtp.PlainAuth("", s.cfg.Email.Address, s.cfg.Email.Password, s.cfg.Email.Host)
	err := e.Send(
		fmt.Sprintf("%s:%d", s.cfg.Email.Host, s.cfg.Email.Port),
		auth,
	)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return hash, nil
}

// VerifyEmail verifies an email using the provided hash
func (s *Service) VerifyEmail(hash string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.verifications[hash]
	if !exists {
		return "", false
	}

	// Check if verification has expired (24 hours)
	if time.Since(data.createdAt) > 24*time.Hour {
		delete(s.verifications, hash)
		return "", false
	}

	// Mark as verified
	data.verified = true
	s.verifications[hash] = data

	return data.email, true
}
