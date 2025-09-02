package storage

import (
	"encoding/json"
	"fmt"
	"strings"

	"demo/bin/bins"
)

const defaultStorageFile = "bins.json"

// FileIO defines minimal file operations for DI
type FileIO interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
}

// Store defines the persistence contract
type Store interface {
	SaveBins(binList *bins.BinList) error
	LoadBins() (*bins.BinList, error)
	GetStoragePath() string
}

// Storage handles persistence of bins to/from JSON files
type Storage struct {
	fileIO   FileIO
	filePath string
}

// New creates a new storage instance with injected FileIO
func New(fileIO FileIO, filePath string) *Storage {
	if filePath == "" {
		filePath = defaultStorageFile
	}
	return &Storage{fileIO: fileIO, filePath: filePath}
}

// SaveBins saves the bin list to a JSON file
func (s *Storage) SaveBins(binList *bins.BinList) error {
	data, err := json.MarshalIndent(binList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal bins to JSON: %w", err)
	}
	if s.fileIO == nil {
		return fmt.Errorf("fileIO is not initialized")
	}
	if err := s.fileIO.Write(s.filePath, data); err != nil {
		return fmt.Errorf("failed to write file %s: %w", s.filePath, err)
	}

	return nil
}

// LoadBins loads the bin list from a JSON file
func (s *Storage) LoadBins() (*bins.BinList, error) {
	if s.fileIO == nil {
		return nil, fmt.Errorf("fileIO is not initialized")
	}
	data, err := s.fileIO.Read(s.filePath)
	if err != nil {
		// Return empty list if file doesn't exist or cannot be read
		return bins.NewList(), nil
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
