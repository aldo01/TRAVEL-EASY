package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"locker-storage-api/config"
	"locker-storage-api/internal/models"
	"locker-storage-api/internal/services"
	"locker-storage-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateLocationRequest struct {
	Name          string  `json:"name" binding:"required"`
	Type          string  `json:"type" binding:"required"`
	Address       string  `json:"address" binding:"required"`
	City          string  `json:"city" binding:"required"`
	State         *string `json:"state"`
	ZipCode       *string `json:"zipCode"`
	Country       string  `json:"country" binding:"required"`
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
	OpeningTime   string  `json:"openingTime"`
	ClosingTime   string  `json:"closingTime"`
	IsOpen24Hours bool    `json:"isOpen24Hours"`
	PhoneNumber   *string `json:"phoneNumber"`
	Email         *string `json:"email"`
	Description   *string `json:"description"`
	ImageURL      *string `json:"imageUrl"`
	Amenities     []string `json:"amenities"`
	HourlyRate    float64 `json:"hourlyRate"`
	DailyRate     float64 `json:"dailyRate"`

	Lockers struct {
		Small  int `json:"small"`
		Medium int `json:"medium"`
		Large  int `json:"large"`
		XLarge int `json:"xlarge"`
	} `json:"lockers"`
}

func CreateLocation(c *gin.Context) {
	var req CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	locationType := models.LocationType(strings.ToUpper(strings.TrimSpace(req.Type)))
	switch locationType {
	case models.LocationKiosk, models.LocationGym, models.LocationClub, models.LocationSchool,
		models.LocationHotel, models.LocationCafe, models.LocationStore, models.LocationOther:
		// ok
	default:
		utils.ErrorResponse(c, 400, "Invalid location type")
		return
	}

	amenitiesJSON := ""
	if len(req.Amenities) > 0 {
		if b, err := json.Marshal(req.Amenities); err == nil {
			amenitiesJSON = string(b)
		}
	}

	location := models.Location{
		Name:          req.Name,
		Type:          locationType,
		Address:       req.Address,
		City:          req.City,
		State:         req.State,
		ZipCode:       req.ZipCode,
		Country:       req.Country,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		OpeningTime:   req.OpeningTime,
		ClosingTime:   req.ClosingTime,
		IsOpen24Hours: req.IsOpen24Hours,
		PhoneNumber:   req.PhoneNumber,
		Email:         req.Email,
		Description:   req.Description,
		ImageURL:      req.ImageURL,
		Amenities:     amenitiesJSON,
		HourlyRate:    req.HourlyRate,
		DailyRate:     req.DailyRate,
		IsActive:      true,
		Rating:        0,
		TotalReviews:  0,
	}

	tx := config.GetDB().Begin()
	if err := tx.Create(&location).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, 500, "Failed to create location")
		return
	}

	createLockers := func(size models.LockerSize, count int, startNumber *int) error {
		for i := 0; i < count; i++ {
			*startNumber = *startNumber + 1
			locker := models.Locker{
				LocationID:     location.ID,
				LockerNumber:   strconv.Itoa(*startNumber),
				Size:           size,
				Status:         models.LockerAvailable,
				IsOperational:  true,
			}
			if err := tx.Create(&locker).Error; err != nil {
				return err
			}
		}
		return nil
	}

	lockerNumber := 0
	if err := createLockers(models.LockerSmall, req.Lockers.Small, &lockerNumber); err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, 500, "Failed to create lockers")
		return
	}
	if err := createLockers(models.LockerMedium, req.Lockers.Medium, &lockerNumber); err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, 500, "Failed to create lockers")
		return
	}
	if err := createLockers(models.LockerLarge, req.Lockers.Large, &lockerNumber); err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, 500, "Failed to create lockers")
		return
	}
	if err := createLockers(models.LockerXLarge, req.Lockers.XLarge, &lockerNumber); err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, 500, "Failed to create lockers")
		return
	}

	tx.Commit()

	utils.SuccessResponse(c, 201, "Location created successfully", gin.H{
		"location": location,
	})
}

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
