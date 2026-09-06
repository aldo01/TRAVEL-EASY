package main

import (
	"log"

	"auth-service/config"
	"auth-service/internal/models"
	"auth-service/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	strPtr := func(s string) *string { return &s }

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := config.GetDB()
	log.Println("Seeding auth-service database...")

	createUser := func(email, name, password string, phone *string) {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			log.Println("✗ Failed to hash password for", email, ":", err)
			return
		}
		user := models.User{
			Email:       email,
			Password:    hashedPassword,
			Name:        name,
			PhoneNumber: phone,
			IsActive:    true,
		}
		if err := db.FirstOrCreate(&user, models.User{Email: user.Email}).Error; err != nil {
			log.Println("✗ Failed to create user:", email, ":", err)
			return
		}
		log.Println("✓ Ensured demo user:", email)
	}

	createUser("test@example.com", "John Doe", "password123", strPtr("+4512345678"))
	createUser("alice@example.com", "Alice Johnson", "password123", strPtr("+14155550101"))
	createUser("bob@example.com", "Bob Singh", "password123", strPtr("+14155550102"))
	createUser("neha@demo.com", "Neha Sharma", "password123", strPtr("+919900000001"))
	createUser("rahul@demo.com", "Rahul Verma", "password123", strPtr("+919900000002"))
	createUser("host@demo.com", "Demo Host", "password123", strPtr("+919900000003"))

	log.Println("✓ auth-service seeding completed")
	log.Println("Test credentials: test@example.com / password123")
}
