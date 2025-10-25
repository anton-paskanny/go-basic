package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetDSN returns database connection string
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntWithDefault gets an environment variable as int with a default value
func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// LoadConfig loads configuration from environment variables
func LoadConfig(path string) (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables and defaults")
	} else {
		log.Println("Loaded configuration from .env file")
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnvIntWithDefault("APP_SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnvWithDefault("APP_DB_HOST", "localhost"),
			Port:     getEnvIntWithDefault("APP_DB_PORT", 5432),
			User:     getEnvWithDefault("APP_DB_USER", "postgres"),
			Password: getEnvWithDefault("APP_DB_PASSWORD", "postgres"),
			DBName:   getEnvWithDefault("APP_DB_NAME", "order_api"),
			SSLMode:  getEnvWithDefault("APP_DB_SSLMODE", "disable"),
		},
	}

	return config, nil
}
