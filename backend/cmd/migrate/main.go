package main

import (
	"log"

	"locker-storage-api/config"
	"locker-storage-api/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := config.GetDB()

	log.Println("Running database migrations...")

	// Auto migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Location{},
		&models.Locker{},
		&models.Booking{},
		&models.AccessLog{},
	)

	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("✓ Migrations completed successfully")
}
