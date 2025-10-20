package database

import (
	"log"

	"gorm.io/gorm"

	"order-api-stat/models"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Add all models to migrate here
	err := db.AutoMigrate(
		&models.Product{},
	)

	if err != nil {
		log.Printf("Error running migrations: %v", err)
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}
