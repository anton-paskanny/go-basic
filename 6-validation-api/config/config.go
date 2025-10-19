package config

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
			Host: "smtp.example.com",
			Port: 587,
		},
		Server: ServerConfig{
			Address: ":8080",
		},
	}
}
