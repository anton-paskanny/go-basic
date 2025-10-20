package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Email  EmailConfig
	Server ServerConfig
}

// EmailConfig holds email configuration
type EmailConfig struct {
	Address  string
	Password string
	Host     string
	Port     int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Address string
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		Email: EmailConfig{
			Address:  getEnv("EMAIL_ADDRESS", ""),
			Password: getEnv("EMAIL_PASSWORD", ""),
			Host:     getEnv("EMAIL_HOST", "smtp.example.com"),
			Port:     getEnvAsInt("EMAIL_PORT", 587),
		},
		Server: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ":8080"),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
