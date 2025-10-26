package database

import (
	"order-api-auth/models"

	"gorm.io/gorm"
)

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB) error {
	// Auto-migrate the schema
	err := db.AutoMigrate(
		&models.User{},
		&models.Session{},
	)
	if err != nil {
		return err
	}

	return nil
}
