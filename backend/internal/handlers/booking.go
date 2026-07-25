package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"locker-storage-api/config"
	"locker-storage-api/internal/models"
	"locker-storage-api/internal/services"
	"locker-storage-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateBookingRequest struct {
	LocationID string  `json:"locationId" binding:"required"`
	LockerID   string  `json:"lockerId" binding:"required"`
	StartTime  string  `json:"startTime" binding:"required"`
	Duration   float64 `json:"duration" binding:"required"` // in hours
	RateType   string  `json:"rateType" binding:"required"` // "HOURLY" or "DAILY"
	Bags       int     `json:"bags"`                         // number of bags (price multiplier)
}

func CreateBooking(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// Parse start time
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid start time format. Use RFC3339")
		return
	}

	if req.Duration <= 0 {
		utils.ErrorResponse(c, 400, "Duration must be greater than 0")
		return
	}

	if req.Duration < 3 {
		utils.ErrorResponse(c, 400, "Minimum booking duration is 3 hours")
		return
	}

	// Disallow bookings starting in the past (small grace for clock skew)
	if startTime.Before(time.Now().Add(-5 * time.Minute)) {
		utils.ErrorResponse(c, 400, "Start time cannot be in the past")
		return
	}

	// Validate locker availability
	var locker models.Locker
	if err := config.GetDB().First(&locker, "id = ?", req.LockerID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Locker not found")
		return
	}

	if locker.Status != models.LockerAvailable {
		utils.ErrorResponse(c, 400, "Locker is not available")
		return
	}

	// Get location for pricing
	var location models.Location
	if err := config.GetDB().First(&location, "id = ?", req.LocationID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Location not found")
		return
	}

	bags := req.Bags
	if bags < 1 {
		bags = 1
	}

	// Calculate price
	var rateType models.RateType
	switch req.RateType {
	case "HOURLY":
		rateType = models.RateHourly
	case "DAILY":
		rateType = models.RateDaily
	case "MULTI_DAY":
		rateType = models.RateMultiDay
	default:
		rateType = models.RateHourly
	}

	totalPrice := services.CalculateBookingPrice(rateType, req.Duration, location.HourlyRate, location.DailyRate) * float64(bags)

	// Calculate end time
	endTime := startTime.Add(time.Duration(req.Duration) * time.Hour)

	// Create booking
	booking := models.Booking{
		UserID:        user.ID,
		LocationID:    req.LocationID,
		LockerID:      req.LockerID,
		BookingNumber: services.GenerateBookingNumber(),
		StartTime:     startTime,
		EndTime:       &endTime,
		RateType:      rateType,
		Bags:          bags,
		BasePrice:     totalPrice,
		TotalPrice:    totalPrice,
		Status:        models.BookingConfirmed,
		QRCode:        services.GenerateQRCode(),
		PaymentStatus: models.PaymentPending,
	}

	// Start transaction
	tx := config.GetDB().Begin()

	// Create booking
	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, 500, "Failed to create booking")
		return
	}

	// Update locker status
	if err := tx.Model(&locker).Update("status", models.LockerReserved).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, 500, "Failed to update locker status")
		return
	}

	tx.Commit()

	// Load relationships
	config.GetDB().Preload("Location").Preload("Locker").First(&booking, "id = ?", booking.ID)

	// Send booking confirmation notification (fire-and-forget)
	go sendBookingNotification(user, booking)

	utils.SuccessResponse(c, 201, "Booking created successfully", booking)
}

func sendBookingNotification(user models.User, booking models.Booking) {
	notificationURL := "http://notification-service:8090/notify/booking-confirmed"

	endTimeStr := ""
	if booking.EndTime != nil {
		endTimeStr = booking.EndTime.Format(time.RFC3339)
	}

	locationAddr := ""
	if booking.Location.Address != "" {
		locationAddr = booking.Location.Address + ", " + booking.Location.City
	}

	payload := map[string]interface{}{
		"bookingNumber":   booking.BookingNumber,
		"userEmail":       user.Email,
		"userName":        user.Name,
		"locationName":    booking.Location.Name,
		"locationAddress": locationAddr,
		"startTime":       booking.StartTime.Format(time.RFC3339),
		"endTime":         endTimeStr,
		"bags":            booking.Bags,
		"rateType":        string(booking.RateType),
		"totalPrice":      booking.TotalPrice,
		"qrCode":          booking.QRCode,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[NOTIFY] Failed to marshal notification payload: %v", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(notificationURL, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("[NOTIFY] Failed to send notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("[NOTIFY] Notification service returned status %d", resp.StatusCode)
	} else {
		log.Printf("[NOTIFY] Booking confirmation notification sent for %s", booking.BookingNumber)
	}
}

func GetUserBookings(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var bookings []models.Booking
	query := config.GetDB().Where("user_id = ?", user.ID).
		Preload("Location").
		Preload("Locker").
		Order("created_at DESC")

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&bookings).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch bookings")
		return
	}

	utils.SuccessResponse(c, 200, "Bookings retrieved", gin.H{
		"count":    len(bookings),
		"bookings": bookings,
	})
}

func GetBooking(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	bookingID := c.Param("id")

	var booking models.Booking
	if err := config.GetDB().
		Preload("Location").
		Preload("Locker").
		Preload("AccessLogs").
		First(&booking, "id = ? AND user_id = ?", bookingID, user.ID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}

	utils.SuccessResponse(c, 200, "Booking retrieved", booking)
}

func UnlockLocker(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	bookingID := c.Param("id")

	type UnlockRequest struct {
		QRCode string `json:"qrCode" binding:"required"`
	}

	var req UnlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// Get booking
	var booking models.Booking
	if err := config.GetDB().First(&booking, "id = ? AND user_id = ?", bookingID, user.ID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}

	// Validate QR code
	if booking.QRCode != req.QRCode {
		utils.ErrorResponse(c, 401, "Invalid QR code")
		return
	}

	// Check if booking is active
	now := time.Now()
	if booking.Status != models.BookingConfirmed && booking.Status != models.BookingActive {
		utils.ErrorResponse(c, 400, "Booking is not active")
		return
	}

	// Check if within booking time
	if now.Before(booking.StartTime) {
		utils.ErrorResponse(c, 400, "Booking has not started yet")
		return
	}

	if booking.EndTime != nil && now.After(*booking.EndTime) {
		utils.ErrorResponse(c, 400, "Booking has expired")
		return
	}

	// Log access
	accessLog := models.AccessLog{
		BookingID: booking.ID,
		Action:    models.AccessUnlock,
		Timestamp: now,
	}

	method := "QR_SCAN"
	accessLog.Method = &method

	if err := config.GetDB().Create(&accessLog).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to log access")
		return
	}

	// Update booking status to ACTIVE if first time
	if booking.Status == models.BookingConfirmed {
		config.GetDB().Model(&booking).Updates(map[string]interface{}{
			"status":        models.BookingActive,
			"check_in_time": now,
		})

		// Update locker status
		config.GetDB().Model(&models.Locker{}).
			Where("id = ?", booking.LockerID).
			Update("status", models.LockerOccupied)
	}

	utils.SuccessResponse(c, 200, "Locker unlocked successfully", gin.H{
		"action":    "UNLOCK",
		"timestamp": now,
	})
}

func CheckIn(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	bookingID := c.Param("id")

	var booking models.Booking
	if err := config.GetDB().First(&booking, "id = ? AND user_id = ?", bookingID, user.ID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}

	now := time.Now()
	config.GetDB().Model(&booking).Updates(map[string]interface{}{
		"status":        models.BookingActive,
		"check_in_time": now,
	})

	config.GetDB().Model(&models.Locker{}).
		Where("id = ?", booking.LockerID).
		Update("status", models.LockerOccupied)

	utils.SuccessResponse(c, 200, "Check-in successful", gin.H{
		"checkInTime": now,
	})
}

func CheckOut(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	bookingID := c.Param("id")

	var booking models.Booking
	if err := config.GetDB().First(&booking, "id = ? AND user_id = ?", bookingID, user.ID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}

	now := time.Now()
	config.GetDB().Model(&booking).Updates(map[string]interface{}{
		"status":         models.BookingCompleted,
		"check_out_time": now,
	})

	config.GetDB().Model(&models.Locker{}).
		Where("id = ?", booking.LockerID).
		Update("status", models.LockerAvailable)

	utils.SuccessResponse(c, 200, "Check-out successful", gin.H{
		"checkOutTime": now,
	})
}

func CancelBooking(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	bookingID := c.Param("id")

	type CancelRequest struct {
		Reason string `json:"reason"`
	}

	var req CancelRequest
	c.ShouldBindJSON(&req)

	var booking models.Booking
	if err := config.GetDB().First(&booking, "id = ? AND user_id = ?", bookingID, user.ID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}

	if booking.Status == models.BookingCompleted || booking.Status == models.BookingCancelled {
		utils.ErrorResponse(c, 400, "Booking cannot be cancelled")
		return
	}

	// Update booking
	config.GetDB().Model(&booking).Updates(map[string]interface{}{
		"status":              models.BookingCancelled,
		"cancellation_reason": req.Reason,
	})

	// Free up locker
	config.GetDB().Model(&models.Locker{}).
		Where("id = ?", booking.LockerID).
		Update("status", models.LockerAvailable)

	utils.SuccessResponse(c, 200, "Booking cancelled successfully", nil)
}
