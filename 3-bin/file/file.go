package file

import (
	"os"
)

// ReadAll reads entire file contents
func ReadAll(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteAll writes data to file with 0644 perms
func WriteAll(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
