package database

import (
	"log"
	"os"

	"github.com/phi-lani/blockchainApp/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDB() {
	dsn := os.Getenv("DB")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.MFAToken{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
