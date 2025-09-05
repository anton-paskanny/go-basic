package config

import (
	"log"
	"os"
)

// Config holds application configuration
type Config struct {
	APIKey string
}

// Load reads configuration from environment variables
func Load() *Config {
	apiKey := os.Getenv("KEY")
	if apiKey == "" {
		log.Fatal("KEY environment variable is required")
	}

	return &Config{
		APIKey: apiKey,
	}
}
