package config

import (
	"os"
)

// Config представляет конфигурацию приложения
type Config struct {
	Port      string
	JWTSecret string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	config := &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}

	return config
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
