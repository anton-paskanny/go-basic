package verify

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

// Error messages for validation
const (
	ErrInvalidEmailFormat = "invalid email format"
	ErrEmptyEmail         = "email cannot be empty"
)

// validateEmail validates an email address format
func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("%s", ErrEmptyEmail)
	}

	// Basic validation using net/mail
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrInvalidEmailFormat, err)
	}

	// Additional validation using regex
	// This regex checks for:
	// - At least one character before @
	// - At least one character after @
	// - At least one dot after @
	// - At least one character after the last dot
	// - No consecutive dots
	// - No spaces
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("%s", ErrInvalidEmailFormat)
	}

	// Check for common invalid patterns
	if strings.Count(email, "@") != 1 {
		return fmt.Errorf("%s", ErrInvalidEmailFormat)
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("%s", ErrInvalidEmailFormat)
	}

	domain := parts[1]
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("%s", ErrInvalidEmailFormat)
	}

	return nil
}
