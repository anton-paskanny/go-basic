package verify

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"

	"validation-api/config"

	"github.com/jordan-wright/email"
)

// Service handles email verification logic
type Service struct {
	cfg     *config.Config
	storage *VerificationStorage
	mu      sync.RWMutex
}

// NewService creates a new verification service
func NewService(cfg *config.Config) *Service {
	// Create storage with default path
	storage := NewVerificationStorage("")

	// Load existing verification data
	if err := storage.Load(); err != nil {
		log.Printf("Warning: Failed to load verification data: %v", err)
	}

	return &Service{
		cfg:     cfg,
		storage: storage,
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
	s.storage.Set(hash, StorableVerificationData{
		Email:     emailAddr,
		CreatedAt: time.Now(),
		Verified:  false,
	})

	// Save to JSON file
	if err := s.storage.Save(); err != nil {
		log.Printf("Warning: Failed to save verification data: %v", err)
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

	// Set up SMTP server address
	smtpServer := fmt.Sprintf("%s:%d", s.cfg.Email.Host, s.cfg.Email.Port)

	// Try to send the email with timeout and error handling
	emailSent := false
	var err error

	// Create a channel to handle timeout
	done := make(chan error, 1)

	// Try to send the email in a separate goroutine
	go func() {
		sendErr := e.Send(smtpServer, auth)
		done <- sendErr
	}()

	// Wait for the email to be sent or timeout
	select {
	case err = <-done:
		if err == nil {
			emailSent = true
		}
	case <-time.After(10 * time.Second):
		err = fmt.Errorf("timeout connecting to SMTP server %s", smtpServer)
	}

	if !emailSent {
		log.Printf("SMTP Error: Failed to send verification email: %v", err)
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return hash, nil
}

// VerifyEmail verifies an email using the provided hash
func (s *Service) VerifyEmail(hash string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.storage.Get(hash)
	if !exists {
		return "", false
	}

	// Check if verification has expired (24 hours)
	if time.Since(data.CreatedAt) > 24*time.Hour {
		s.storage.Delete(hash)
		if err := s.storage.Save(); err != nil {
			log.Printf("Warning: Failed to save verification data after expiration: %v", err)
		}
		return "", false
	}

	// Get email before deleting
	email := data.Email

	// Delete the record after successful verification
	s.storage.Delete(hash)
	if err := s.storage.Save(); err != nil {
		log.Printf("Warning: Failed to save verification data after verification: %v", err)
	}

	return email, true
}
