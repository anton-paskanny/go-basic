package verify

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// DefaultStorageFileName is the default file name for storing verification data
	DefaultStorageFileName = "verification_data.json"

	// DefaultStorageDir is the default directory for storing verification data
	DefaultStorageDir = "data"
)

// StorableVerificationData represents verification data that can be stored in JSON
type StorableVerificationData struct {
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Verified  bool      `json:"verified"`
}

// VerificationStorage handles saving and loading verification data
type VerificationStorage struct {
	filePath      string
	verifications map[string]StorableVerificationData
	mu            sync.RWMutex
}

// GetDefaultStoragePath returns the default path for the verification data file
func GetDefaultStoragePath() string {
	return filepath.Join(DefaultStorageDir, DefaultStorageFileName)
}

// NewVerificationStorage creates a new verification storage
func NewVerificationStorage(filePath string) *VerificationStorage {
	if filePath == "" {
		filePath = GetDefaultStoragePath()
	}

	return &VerificationStorage{
		filePath:      filePath,
		verifications: make(map[string]StorableVerificationData),
	}
}

// Load loads verification data from the JSON file
func (s *VerificationStorage) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// File doesn't exist, initialize empty map
		s.verifications = make(map[string]StorableVerificationData)
		return nil
	}

	// Read file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read verification data file: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &s.verifications); err != nil {
		return fmt.Errorf("failed to parse verification data: %w", err)
	}

	return nil
}

// Save saves verification data to the JSON file
func (s *VerificationStorage) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for verification data: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(s.verifications, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal verification data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write verification data: %w", err)
	}

	return nil
}

// Get gets verification data for a hash
func (s *VerificationStorage) Get(hash string) (StorableVerificationData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.verifications[hash]
	return data, exists
}

// Set sets verification data for a hash
func (s *VerificationStorage) Set(hash string, data StorableVerificationData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.verifications[hash] = data
}

// Delete deletes verification data for a hash
func (s *VerificationStorage) Delete(hash string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.verifications, hash)
}

// GetAll gets all verification data
func (s *VerificationStorage) GetAll() map[string]StorableVerificationData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy to avoid concurrent access issues
	result := make(map[string]StorableVerificationData, len(s.verifications))
	for k, v := range s.verifications {
		result[k] = v
	}

	return result
}
