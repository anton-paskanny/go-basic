package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator is a custom validator instance
type Validator struct {
	validator *validator.Validate
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()

	// Register custom validation tag names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validator: v,
	}
}

// Validate validates the given struct and returns validation errors
func (v *Validator) Validate(i interface{}) map[string]string {
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		errors[field] = getErrorMessage(field, tag, err.Param())
	}

	return errors
}

// getErrorMessage returns a human-readable error message for a validation error
func getErrorMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "min":
		return field + " must be at least " + param + " characters long"
	case "max":
		return field + " must be at most " + param + " characters long"
	case "gt":
		return field + " must be greater than " + param
	case "email":
		return field + " must be a valid email address"
	default:
		return field + " is invalid"
	}
}
