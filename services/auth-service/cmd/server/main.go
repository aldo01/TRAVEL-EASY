package main

import (
	"log"
	"os"

	"auth-service/config"
	"auth-service/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "auth-service"})
	})

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}
		v1.GET("/profile", handlers.GetProfile)
	}

	// Internal service-to-service routes (not exposed by the gateway).
	internal := router.Group("/internal")
	{
		internal.GET("/users/:id", handlers.GetUserInternal)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("auth-service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
