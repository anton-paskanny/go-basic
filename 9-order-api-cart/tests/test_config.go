package tests

import (
	"os"
	"testing"

	"order-api-cart/config"
	"order-api-cart/database"
)

// TestConfig holds test configuration
type TestConfig struct {
	*config.Config
	TestDBName string
}

// LoadTestConfig loads test configuration
func LoadTestConfig() *TestConfig {
	// Set test environment variables
	os.Setenv("DB_NAME", "order_cart_test_db")
	os.Setenv("SERVER_PORT", "8083")
	os.Setenv("AUTH_SERVICE_URL", "http://localhost:8084")
	os.Setenv("PRODUCT_SERVICE_URL", "http://localhost:8085")
	os.Setenv("JWT_SECRET", "test-secret-key")

	cfg := config.LoadConfig()
	return &TestConfig{
		Config:     cfg,
		TestDBName: "order_cart_test_db",
	}
}

// SetupTestDB sets up test database
func SetupTestDB(t *testing.T) *TestConfig {
	cfg := LoadTestConfig()

	// Connect to test database
	if err := database.Connect(cfg.Config); err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return cfg
}

// CleanupTestDB cleans up test database
func CleanupTestDB(t *testing.T) {
	db := database.GetDB()
	if db != nil {
		// Drop all tables
		db.Exec("DROP SCHEMA public CASCADE")
		db.Exec("CREATE SCHEMA public")
		db.Exec("GRANT ALL ON SCHEMA public TO postgres")
		db.Exec("GRANT ALL ON SCHEMA public TO public")
	}
}
