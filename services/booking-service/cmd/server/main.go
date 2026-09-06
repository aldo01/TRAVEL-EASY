package main

import (
	"log"
	"os"

	"booking-service/config"
	"booking-service/internal/events"
	"booking-service/internal/handlers"
	"booking-service/internal/middleware"

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

	// Connect the async event publisher (booking.confirmed).
	events.Init()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "booking-service"})
	})

	v1 := router.Group("/api/v1")
	{
		bookings := v1.Group("/bookings")
		bookings.Use(middleware.RequireUser())
		{
			bookings.POST("", handlers.CreateBooking)
			bookings.GET("", handlers.GetUserBookings)
			bookings.GET("/:id", handlers.GetBooking)
			bookings.POST("/:id/unlock", handlers.UnlockLocker)
			bookings.POST("/:id/checkin", handlers.CheckIn)
			bookings.POST("/:id/checkout", handlers.CheckOut)
			bookings.POST("/:id/cancel", handlers.CancelBooking)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	log.Printf("booking-service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
