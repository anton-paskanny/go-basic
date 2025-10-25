package utils

import (
	"regexp"
	"strings"
)

// ValidatePhone validates phone number correctness
func ValidatePhone(phone string) error {
	if strings.TrimSpace(phone) == "" {
		return &ValidationError{Field: "phone", Message: "Phone number is required"}
	}

	// Remove all non-digit characters
	cleanPhone := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Check length (should be 10-15 digits)
	if len(cleanPhone) < 10 || len(cleanPhone) > 15 {
		return &ValidationError{Field: "phone", Message: "Phone number must be 10-15 digits"}
	}

	return nil
}

// ValidateCode validates confirmation code correctness
func ValidateCode(code string) error {
	if strings.TrimSpace(code) == "" {
		return &ValidationError{Field: "code", Message: "Code is required"}
	}

	// Check that code contains only digits
	if !regexp.MustCompile(`^\d{4}$`).MatchString(code) {
		return &ValidationError{Field: "code", Message: "Code must be 4 digits"}
	}

	return nil
}

// ValidateSessionID validates session ID correctness
func ValidateSessionID(sessionID string) error {
	if strings.TrimSpace(sessionID) == "" {
		return &ValidationError{Field: "sessionId", Message: "Session ID is required"}
	}

	// Check basic UUID length
	if len(sessionID) < 10 {
		return &ValidationError{Field: "sessionId", Message: "Invalid session ID format"}
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
