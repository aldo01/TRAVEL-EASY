package handlers

import (
	"location-service/config"
	"location-service/internal/models"
	"location-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// LockerContext is returned to booking-service so it can price and snapshot a booking.
type LockerContext struct {
	LockerID     string  `json:"lockerId"`
	LockerNumber string  `json:"lockerNumber"`
	LockerSize   string  `json:"lockerSize"`
	LockerStatus string  `json:"lockerStatus"`
	LocationID   string  `json:"locationId"`
	LocationName string  `json:"locationName"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	HourlyRate   float64 `json:"hourlyRate"`
	DailyRate    float64 `json:"dailyRate"`
}

// GetLockerContext (internal) returns locker + location details for a booking.
func GetLockerContext(c *gin.Context) {
	lockerID := c.Param("id")

	var locker models.Locker
	if err := config.GetDB().First(&locker, "id = ?", lockerID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Locker not found")
		return
	}

	var location models.Location
	if err := config.GetDB().First(&location, "id = ?", locker.LocationID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Location not found")
		return
	}

	utils.SuccessResponse(c, 200, "Locker context", LockerContext{
		LockerID:     locker.ID,
		LockerNumber: locker.LockerNumber,
		LockerSize:   string(locker.Size),
		LockerStatus: string(locker.Status),
		LocationID:   location.ID,
		LocationName: location.Name,
		Address:      location.Address,
		City:         location.City,
		HourlyRate:   location.HourlyRate,
		DailyRate:    location.DailyRate,
	})
}

// ReserveLocker (internal) atomically moves an AVAILABLE locker to RESERVED.
func ReserveLocker(c *gin.Context) {
	lockerID := c.Param("id")

	res := config.GetDB().Model(&models.Locker{}).
		Where("id = ? AND status = ? AND is_operational = ?", lockerID, models.LockerAvailable, true).
		Update("status", models.LockerReserved)

	if res.Error != nil {
		utils.ErrorResponse(c, 500, "Failed to reserve locker")
		return
	}
	if res.RowsAffected == 0 {
		utils.ErrorResponse(c, 409, "Locker is not available")
		return
	}

	utils.SuccessResponse(c, 200, "Locker reserved", nil)
}

// setLockerStatus is a helper for release/occupy transitions.
func setLockerStatus(c *gin.Context, status models.LockerStatus) {
	lockerID := c.Param("id")

	res := config.GetDB().Model(&models.Locker{}).
		Where("id = ?", lockerID).
		Update("status", status)

	if res.Error != nil {
		utils.ErrorResponse(c, 500, "Failed to update locker status")
		return
	}
	if res.RowsAffected == 0 {
		utils.ErrorResponse(c, 404, "Locker not found")
		return
	}

	utils.SuccessResponse(c, 200, "Locker status updated", nil)
}

// ReleaseLocker (internal) frees a locker back to AVAILABLE.
func ReleaseLocker(c *gin.Context) {
	setLockerStatus(c, models.LockerAvailable)
}

// OccupyLocker (internal) marks a locker OCCUPIED.
func OccupyLocker(c *gin.Context) {
	setLockerStatus(c, models.LockerOccupied)
}
