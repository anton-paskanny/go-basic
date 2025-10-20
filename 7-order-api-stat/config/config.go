package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// GetDSN returns database connection string
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(path string) (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables and defaults")
	} else {
		log.Println("Loaded configuration from .env file")
	}

	// Setup viper
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Enable viper to read environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP") // will be uppercased automatically

	// Map environment variables to config keys
	viper.BindEnv("server.port", "APP_SERVER_PORT")
	viper.BindEnv("database.host", "APP_DB_HOST")
	viper.BindEnv("database.port", "APP_DB_PORT")
	viper.BindEnv("database.user", "APP_DB_USER")
	viper.BindEnv("database.password", "APP_DB_PASSWORD")
	viper.BindEnv("database.dbname", "APP_DB_NAME")
	viper.BindEnv("database.sslmode", "APP_DB_SSLMODE")

	// Set default values
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "order_api")
	viper.SetDefault("database.sslmode", "disable")

	// Try to read config file, but don't return error if not found
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables and defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}
