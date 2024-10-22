// database/database.go
package database

import (
	"log"
	"os"

	"github.com/phi-lani/blockchainApp/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitializeDB - Connects to the PostgreSQL database
func InitializeDB() {
	dsn := os.Getenv("DB")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate all models (this creates tables if they don't exist)
	err = DB.AutoMigrate(&models.User{}) // Add all models that need to be migrated
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
