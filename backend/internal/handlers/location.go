package handlers

import (
	"strconv"

	"locker-storage-api/config"
	"locker-storage-api/internal/models"
	"locker-storage-api/internal/services"
	"locker-storage-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetLocations(c *gin.Context) {
	var locations []models.Location

	query := config.GetDB().Where("is_active = ?", true)

	// Filter by city if provided
	if city := c.Query("city"); city != "" {
		query = query.Where("city ILIKE ?", "%"+city+"%")
	}

	// Filter by type if provided
	if locationType := c.Query("type"); locationType != "" {
		query = query.Where("type = ?", locationType)
	}

	if err := query.Find(&locations).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch locations")
		return
	}

	utils.SuccessResponse(c, 200, "Locations retrieved", gin.H{
		"count":     len(locations),
		"locations": locations,
	})
}

func GetNearbyLocations(c *gin.Context) {
	// Get query parameters
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	radiusStr := c.DefaultQuery("radius", "10") // Default 10km

	if latStr == "" || lonStr == "" {
		utils.ErrorResponse(c, 400, "Latitude and longitude are required")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid latitude")
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid longitude")
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid radius")
		return
	}

	// Get all active locations
	var locations []models.Location
	if err := config.GetDB().Where("is_active = ?", true).Find(&locations).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch locations")
		return
	}

	// Find nearby locations
	nearby := services.FindNearbyLocations(lat, lon, radius, locations)

	// Get available locker count for each location
	for i := range nearby {
		var count int64
		config.GetDB().Model(&models.Locker{}).
			Where("location_id = ? AND status = ? AND is_operational = ?",
				nearby[i].ID, models.LockerAvailable, true).
			Count(&count)
		nearby[i].AvailableLockers = int(count)
	}

	utils.SuccessResponse(c, 200, "Nearby locations found", gin.H{
		"count":     len(nearby),
		"locations": nearby,
	})
}

func GetLocation(c *gin.Context) {
	locationID := c.Param("id")

	var location models.Location
	if err := config.GetDB().First(&location, "id = ?", locationID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Location not found")
		return
	}

	// Get available locker count
	var availableCount int64
	config.GetDB().Model(&models.Locker{}).
		Where("location_id = ? AND status = ? AND is_operational = ?",
			locationID, models.LockerAvailable, true).
		Count(&availableCount)

	utils.SuccessResponse(c, 200, "Location retrieved", gin.H{
		"location":         location,
		"availableLockers": availableCount,
	})
}

func GetLocationLockers(c *gin.Context) {
	locationID := c.Param("id")

	var lockers []models.Locker
	query := config.GetDB().Where("location_id = ?", locationID)

	// Filter by size if provided
	if size := c.Query("size"); size != "" {
		query = query.Where("size = ?", size)
	}

	// Filter by status (default to AVAILABLE)
	status := c.DefaultQuery("status", string(models.LockerAvailable))
	query = query.Where("status = ?", status).Where("is_operational = ?", true)

	if err := query.Find(&lockers).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch lockers")
		return
	}

	utils.SuccessResponse(c, 200, "Lockers retrieved", gin.H{
		"count":   len(lockers),
		"lockers": lockers,
	})
}

func GetAvailableLockers(c *gin.Context) {
	var lockers []models.Locker
	query := config.GetDB().
		Where("status = ? AND is_operational = ?", models.LockerAvailable, true).
		Preload("Location")

	// Filter by size if provided
	if size := c.Query("size"); size != "" {
		query = query.Where("size = ?", size)
	}

	// Filter by location ID if provided
	if locationID := c.Query("locationId"); locationID != "" {
		query = query.Where("location_id = ?", locationID)
	}

	if err := query.Find(&lockers).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch lockers")
		return
	}

	utils.SuccessResponse(c, 200, "Available lockers retrieved", gin.H{
		"count":   len(lockers),
		"lockers": lockers,
	})
}
