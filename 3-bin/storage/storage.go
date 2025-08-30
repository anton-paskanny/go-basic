package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"demo/bin/bins"
)

const defaultStorageFile = "bins.json"

// Storage handles persistence of bins to/from JSON files
type Storage struct {
	filePath string
}

// New creates a new storage instance
func New(filePath string) *Storage {
	if filePath == "" {
		filePath = defaultStorageFile
	}
	return &Storage{filePath: filePath}
}

// SaveBins saves the bin list to a JSON file
func (s *Storage) SaveBins(binList *bins.BinList) error {
	data, err := json.MarshalIndent(binList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal bins to JSON: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", s.filePath, err)
	}

	return nil
}

// LoadBins loads the bin list from a JSON file
func (s *Storage) LoadBins() (*bins.BinList, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty list if file doesn't exist
			return bins.NewList(), nil
		}
		return nil, fmt.Errorf("failed to read file %s: %w", s.filePath, err)
	}

	var binList bins.BinList
	if err := json.Unmarshal(data, &binList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", s.filePath, err)
	}

	return &binList, nil
}

// IsJSONFile checks if the given file has a JSON extension
func IsJSONFile(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".json")
}

// GetStoragePath returns the current storage file path
func (s *Storage) GetStoragePath() string {
	return s.filePath
}
