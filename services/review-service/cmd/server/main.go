package main

import (
	"log"
	"os"

	"review-service/config"
	"review-service/internal/events"
	"review-service/internal/handlers"
	"review-service/internal/middleware"

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

	events.Init()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "review-service"})
	})

	v1 := router.Group("/api/v1")
	{
		reviews := v1.Group("/reviews")
		{
			reviews.GET("", handlers.GetLocationReviews)
			reviews.POST("", middleware.RequireUser(), handlers.CreateReview)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	log.Printf("review-service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
