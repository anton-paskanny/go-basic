package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadAll reads entire file contents
func ReadAll(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteAll writes data to file with 0644 perms
func WriteAll(path string, data []byte) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	return os.WriteFile(path, data, 0644)
}

// IsJSONFile checks if the given file has a JSON extension
func IsJSONFile(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".json")
}

// ValidateJSONFile checks if a file is a valid JSON file
func ValidateJSONFile(path string) error {
	if !IsJSONFile(path) {
		return fmt.Errorf("file %s does not have a .json extension", path)
	}

	data, err := ReadAll(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Try to parse as JSON to validate
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("file %s contains invalid JSON: %w", path, err)
	}

	return nil
}

// ReadJSONFile reads and validates a JSON file, returning the raw bytes
func ReadJSONFile(path string) ([]byte, error) {
	if err := ValidateJSONFile(path); err != nil {
		return nil, err
	}

	return ReadAll(path)
}
