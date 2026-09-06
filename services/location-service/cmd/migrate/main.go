package main

import (
	"log"

	"location-service/config"
	"location-service/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Running location-service migrations...")
	if err := models.AutoMigrate(config.GetDB()); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("✓ Migrations completed successfully")
}
