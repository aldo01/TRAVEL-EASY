package main

import (
	"log"
	"os"

	"locker-storage-api/config"
	"locker-storage-api/internal/handlers"
	"locker-storage-api/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Locker Storage API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes - Authentication
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Public routes - Locations
		locations := v1.Group("/locations")
		{
			locations.GET("", handlers.GetLocations)
			locations.GET("/nearby", handlers.GetNearbyLocations)
			locations.GET("/:id", handlers.GetLocation)
			locations.GET("/:id/lockers", handlers.GetLocationLockers)
		}

		// Public routes - Lockers
		v1.GET("/lockers/available", handlers.GetAvailableLockers)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/profile", handlers.GetProfile)

			// Bookings
			bookings := protected.Group("/bookings")
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
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
