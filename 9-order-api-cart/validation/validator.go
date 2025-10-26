package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground validator
type Validator struct {
	validator *validator.Validate
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()

	// Register custom validators if needed
	// v.RegisterValidation("custom_tag", customValidationFunction)

	return &Validator{
		validator: v,
	}
}

// ValidateStruct validates a struct and returns validation errors
func (v *Validator) ValidateStruct(s interface{}) error {
	err := v.validator.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors ValidationErrors

	for _, err := range err.(validator.ValidationErrors) {
		validationError := ValidationError{
			Field: err.Field(),
			Tag:   err.Tag(),
			Value: fmt.Sprintf("%v", err.Value()),
		}

		// Generate human-readable error message
		validationError.Message = v.getErrorMessage(err)
		validationErrors.Errors = append(validationErrors.Errors, validationError)
	}

	return validationErrors
}

// getErrorMessage generates a human-readable error message
func (v *Validator) getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", fe.Field(), fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fe.Field())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", fe.Field())
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", fe.Field())
	case "numeric":
		return fmt.Sprintf("%s must contain only numeric characters", fe.Field())
	case "alphaunicode":
		return fmt.Sprintf("%s must contain only unicode alphabetic characters", fe.Field())
	case "alphanumunicode":
		return fmt.Sprintf("%s must contain only unicode alphanumeric characters", fe.Field())
	case "boolean":
		return fmt.Sprintf("%s must be a boolean value", fe.Field())
	case "datetime":
		return fmt.Sprintf("%s must be a valid datetime", fe.Field())
	case "json":
		return fmt.Sprintf("%s must be valid JSON", fe.Field())
	case "file":
		return fmt.Sprintf("%s must be a valid file path", fe.Field())
	case "uri":
		return fmt.Sprintf("%s must be a valid URI", fe.Field())
	case "base64":
		return fmt.Sprintf("%s must be valid base64 encoded string", fe.Field())
	case "base64url":
		return fmt.Sprintf("%s must be valid base64url encoded string", fe.Field())
	case "contains":
		return fmt.Sprintf("%s must contain '%s'", fe.Field(), fe.Param())
	case "containsany":
		return fmt.Sprintf("%s must contain at least one of: %s", fe.Field(), fe.Param())
	case "containsrune":
		return fmt.Sprintf("%s must contain the rune '%s'", fe.Field(), fe.Param())
	case "excludes":
		return fmt.Sprintf("%s must not contain '%s'", fe.Field(), fe.Param())
	case "excludesall":
		return fmt.Sprintf("%s must not contain any of: %s", fe.Field(), fe.Param())
	case "excludesrune":
		return fmt.Sprintf("%s must not contain the rune '%s'", fe.Field(), fe.Param())
	case "startswith":
		return fmt.Sprintf("%s must start with '%s'", fe.Field(), fe.Param())
	case "endswith":
		return fmt.Sprintf("%s must end with '%s'", fe.Field(), fe.Param())
	case "startsnotwith":
		return fmt.Sprintf("%s must not start with '%s'", fe.Field(), fe.Param())
	case "endsnotwith":
		return fmt.Sprintf("%s must not end with '%s'", fe.Field(), fe.Param())
	case "ip":
		return fmt.Sprintf("%s must be a valid IP address", fe.Field())
	case "ipv4":
		return fmt.Sprintf("%s must be a valid IPv4 address", fe.Field())
	case "ipv6":
		return fmt.Sprintf("%s must be a valid IPv6 address", fe.Field())
	case "cidr":
		return fmt.Sprintf("%s must be a valid CIDR", fe.Field())
	case "cidrv4":
		return fmt.Sprintf("%s must be a valid IPv4 CIDR", fe.Field())
	case "cidrv6":
		return fmt.Sprintf("%s must be a valid IPv6 CIDR", fe.Field())
	case "mac":
		return fmt.Sprintf("%s must be a valid MAC address", fe.Field())
	case "hostname":
		return fmt.Sprintf("%s must be a valid hostname", fe.Field())
	case "hostname_port":
		return fmt.Sprintf("%s must be a valid hostname with port", fe.Field())
	case "fqdn":
		return fmt.Sprintf("%s must be a valid FQDN", fe.Field())
	case "unique":
		return fmt.Sprintf("%s must contain unique values", fe.Field())
	case "iscolor":
		return fmt.Sprintf("%s must be a valid color", fe.Field())
	case "oneofcaseinsensitive":
		return fmt.Sprintf("%s must be one of (case insensitive): %s", fe.Field(), fe.Param())
	case "printascii":
		return fmt.Sprintf("%s must contain only printable ASCII characters", fe.Field())
	case "printunicode":
		return fmt.Sprintf("%s must contain only printable unicode characters", fe.Field())
	case "multibyte":
		return fmt.Sprintf("%s must contain multibyte characters", fe.Field())
	case "datauri":
		return fmt.Sprintf("%s must be a valid data URI", fe.Field())
	case "latitude":
		return fmt.Sprintf("%s must be a valid latitude", fe.Field())
	case "longitude":
		return fmt.Sprintf("%s must be a valid longitude", fe.Field())
	case "ssn":
		return fmt.Sprintf("%s must be a valid SSN", fe.Field())
	case "sin":
		return fmt.Sprintf("%s must be a valid SIN", fe.Field())
	case "luhn":
		return fmt.Sprintf("%s must pass the Luhn algorithm check", fe.Field())
	case "credit_card":
		return fmt.Sprintf("%s must be a valid credit card number", fe.Field())
	case "ean":
		return fmt.Sprintf("%s must be a valid EAN", fe.Field())
	case "gtin":
		return fmt.Sprintf("%s must be a valid GTIN", fe.Field())
	case "isbn":
		return fmt.Sprintf("%s must be a valid ISBN", fe.Field())
	case "isbn10":
		return fmt.Sprintf("%s must be a valid ISBN-10", fe.Field())
	case "isbn13":
		return fmt.Sprintf("%s must be a valid ISBN-13", fe.Field())
	case "issn":
		return fmt.Sprintf("%s must be a valid ISSN", fe.Field())
	case "uuid3":
		return fmt.Sprintf("%s must be a valid UUID v3", fe.Field())
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID v4", fe.Field())
	case "uuid5":
		return fmt.Sprintf("%s must be a valid UUID v5", fe.Field())
	case "ulid":
		return fmt.Sprintf("%s must be a valid ULID", fe.Field())
	case "cron":
		return fmt.Sprintf("%s must be a valid cron expression", fe.Field())
	case "dive":
		return fmt.Sprintf("%s validation failed", fe.Field())
	case "keys":
		return fmt.Sprintf("%s keys validation failed", fe.Field())
	case "endkeys":
		return fmt.Sprintf("%s end keys validation failed", fe.Field())
	case "required_with":
		return fmt.Sprintf("%s is required when %s is present", fe.Field(), fe.Param())
	case "required_with_all":
		return fmt.Sprintf("%s is required when all of %s are present", fe.Field(), fe.Param())
	case "required_without":
		return fmt.Sprintf("%s is required when %s is not present", fe.Field(), fe.Param())
	case "required_without_all":
		return fmt.Sprintf("%s is required when none of %s are present", fe.Field(), fe.Param())
	case "excluded_with":
		return fmt.Sprintf("%s must be excluded when %s is present", fe.Field(), fe.Param())
	case "excluded_with_all":
		return fmt.Sprintf("%s must be excluded when all of %s are present", fe.Field(), fe.Param())
	case "excluded_without":
		return fmt.Sprintf("%s must be excluded when %s is not present", fe.Field(), fe.Param())
	case "excluded_without_all":
		return fmt.Sprintf("%s must be excluded when none of %s are present", fe.Field(), fe.Param())
	case "isdefault":
		return fmt.Sprintf("%s must be the default value", fe.Field())
	case "eq":
		return fmt.Sprintf("%s must be equal to %s", fe.Field(), fe.Param())
	case "ne":
		return fmt.Sprintf("%s must not be equal to %s", fe.Field(), fe.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", fe.Field(), fe.Param())
	case "eqcsfield":
		return fmt.Sprintf("%s must be equal to %s", fe.Field(), fe.Param())
	case "necsfield":
		return fmt.Sprintf("%s must not be equal to %s", fe.Field(), fe.Param())
	case "gtcsfield":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "gtecsfield":
		return fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "ltcsfield":
		return fmt.Sprintf("%s must be less than %s", fe.Field(), fe.Param())
	case "ltecsfield":
		return fmt.Sprintf("%s must be less than or equal to %s", fe.Field(), fe.Param())
	case "nefield":
		return fmt.Sprintf("%s must not be equal to %s", fe.Field(), fe.Param())
	case "gtefield":
		return fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "ltefield":
		return fmt.Sprintf("%s must be less than or equal to %s", fe.Field(), fe.Param())
	case "gtfield":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "ltfield":
		return fmt.Sprintf("%s must be less than %s", fe.Field(), fe.Param())
	case "fieldcontains":
		return fmt.Sprintf("%s must contain %s", fe.Field(), fe.Param())
	case "fieldexcludes":
		return fmt.Sprintf("%s must not contain %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is not valid", fe.Field())
	}
}
