package main

import (
	"log"
	"os"

	"location-service/config"
	"location-service/internal/events"
	"location-service/internal/handlers"
	"location-service/internal/middleware"

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

	// Keep location ratings in sync with review-service via async events.
	go events.SubscribeReviewEvents()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "location-service"})
	})

	v1 := router.Group("/api/v1")
	{
		locations := v1.Group("/locations")
		{
			locations.GET("", handlers.GetLocations)
			locations.GET("/nearby", handlers.GetNearbyLocations)
			locations.GET("/:id", handlers.GetLocation)
			locations.GET("/:id/lockers", handlers.GetLocationLockers)
			locations.POST("", middleware.RequireUser(), handlers.CreateLocation)
		}

		v1.GET("/lockers/available", handlers.GetAvailableLockers)
		v1.GET("/geocode/suggest", handlers.SuggestPlaces)
	}

	// Internal service-to-service routes (not exposed by the gateway).
	internal := router.Group("/internal")
	{
		internal.GET("/lockers/:id/context", handlers.GetLockerContext)
		internal.POST("/lockers/:id/reserve", handlers.ReserveLocker)
		internal.POST("/lockers/:id/release", handlers.ReleaseLocker)
		internal.POST("/lockers/:id/occupy", handlers.OccupyLocker)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("location-service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
