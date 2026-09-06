package handlers

import (
	"log"
	"time"

	"booking-service/config"
	"booking-service/internal/clients"
	"booking-service/internal/events"
	"booking-service/internal/models"
	"booking-service/internal/services"
	"booking-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateBookingRequest struct {
	LocationID string  `json:"locationId" binding:"required"`
	LockerID   string  `json:"lockerId" binding:"required"`
	StartTime  string  `json:"startTime" binding:"required"`
	Duration   float64 `json:"duration" binding:"required"`
	RateType   string  `json:"rateType" binding:"required"`
	Bags       int     `json:"bags"`
}

func CreateBooking(c *gin.Context) {
	userID := c.GetString("userId")
	userEmail := c.GetString("userEmail")

	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

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
	if startTime.Before(time.Now().Add(-5 * time.Minute)) {
		utils.ErrorResponse(c, 400, "Start time cannot be in the past")
		return
	}

	// Fetch locker + location context from location-service.
	ctx, err := clients.GetLockerContext(req.LockerID)
	if err != nil {
		utils.ErrorResponse(c, 404, "Locker not found")
		return
	}
	if ctx.LocationID != req.LocationID {
		utils.ErrorResponse(c, 400, "Locker does not belong to the specified location")
		return
	}
	if ctx.LockerStatus != "AVAILABLE" {
		utils.ErrorResponse(c, 400, "Locker is not available")
		return
	}

	bags := req.Bags
	if bags < 1 {
		bags = 1
	}

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

	totalPrice := services.CalculateBookingPrice(rateType, req.Duration, ctx.HourlyRate, ctx.DailyRate) * float64(bags)
	endTime := startTime.Add(time.Duration(req.Duration) * time.Hour)

	// Reserve the locker in location-service (atomic AVAILABLE -> RESERVED).
	if err := clients.ReserveLocker(req.LockerID); err != nil {
		utils.ErrorResponse(c, 409, "Locker is no longer available")
		return
	}

	booking := models.Booking{
		UserID:          userID,
		LocationID:      req.LocationID,
		LockerID:        req.LockerID,
		BookingNumber:   services.GenerateBookingNumber(),
		StartTime:       startTime,
		EndTime:         &endTime,
		RateType:        rateType,
		Bags:            bags,
		BasePrice:       totalPrice,
		TotalPrice:      totalPrice,
		Status:          models.BookingConfirmed,
		QRCode:          services.GenerateQRCode(),
		PaymentStatus:   models.PaymentPending,
		LocationName:    ctx.LocationName,
		LocationAddress: ctx.Address,
		LocationCity:    ctx.City,
		LockerNumber:    ctx.LockerNumber,
		LockerSize:      ctx.LockerSize,
	}

	if err := config.GetDB().Create(&booking).Error; err != nil {
		// Compensation: release the reserved locker.
		if relErr := clients.ReleaseLocker(req.LockerID); relErr != nil {
			log.Printf("[SAGA] failed to release locker after booking error: %v", relErr)
		}
		utils.ErrorResponse(c, 500, "Failed to create booking")
		return
	}

	// Publish confirmation event (async email via notification-service).
	name := ""
	email := userEmail
	phone := ""
	if u, uErr := clients.GetUser(userID); uErr == nil {
		name = u.Name
		phone = u.Phone
		if u.Email != "" {
			email = u.Email
		}
	}
	endStr := ""
	if booking.EndTime != nil {
		endStr = booking.EndTime.Format(time.RFC3339)
	}
	addr := ctx.Address
	if ctx.City != "" {
		addr = ctx.Address + ", " + ctx.City
	}
	events.PublishBookingConfirmed(events.BookingConfirmedEvent{
		BookingNumber:   booking.BookingNumber,
		UserEmail:       email,
		UserName:        name,
		UserPhone:       phone,
		LocationName:    ctx.LocationName,
		LocationAddress: addr,
		StartTime:       booking.StartTime.Format(time.RFC3339),
		EndTime:         endStr,
		Bags:            booking.Bags,
		RateType:        string(booking.RateType),
		TotalPrice:      booking.TotalPrice,
		QRCode:          booking.QRCode,
	})

	booking.Enrich()
	utils.SuccessResponse(c, 201, "Booking created successfully", booking)
}

func GetUserBookings(c *gin.Context) {
	userID := c.GetString("userId")

	var bookings []models.Booking
	query := config.GetDB().Where("user_id = ?", userID).Order("created_at DESC")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&bookings).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch bookings")
		return
	}

	for i := range bookings {
		bookings[i].Enrich()
	}

	utils.SuccessResponse(c, 200, "Bookings retrieved", gin.H{
		"count":    len(bookings),
		"bookings": bookings,
	})
}

func GetBooking(c *gin.Context) {
	userID := c.GetString("userId")
	bookingID := c.Param("id")

	var booking models.Booking
	if err := config.GetDB().
		Preload("AccessLogs").
		First(&booking, "id = ? AND user_id = ?", bookingID, userID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}

	booking.Enrich()
	utils.SuccessResponse(c, 200, "Booking retrieved", booking)
}

func UnlockLocker(c *gin.Context) {
	userID := c.GetString("userId")
	bookingID := c.Param("id")

	type UnlockRequest struct {
		QRCode string `json:"qrCode" binding:"required"`
	}
	var req UnlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	var booking models.Booking
	if err := config.GetDB().First(&booking, "id = ? AND user_id = ?", bookingID, userID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}
	if booking.QRCode != req.QRCode {
		utils.ErrorResponse(c, 401, "Invalid QR code")
		return
	}

	now := time.Now()
	if booking.Status != models.BookingConfirmed && booking.Status != models.BookingActive {
		utils.ErrorResponse(c, 400, "Booking is not active")
		return
	}
	if now.Before(booking.StartTime) {
		utils.ErrorResponse(c, 400, "Booking has not started yet")
		return
	}
	if booking.EndTime != nil && now.After(*booking.EndTime) {
		utils.ErrorResponse(c, 400, "Booking has expired")
		return
	}

	method := "QR_SCAN"
	accessLog := models.AccessLog{
		BookingID: booking.ID,
		Action:    models.AccessUnlock,
		Timestamp: now,
		Method:    &method,
	}
	if err := config.GetDB().Create(&accessLog).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to log access")
		return
	}

	if booking.Status == models.BookingConfirmed {
		config.GetDB().Model(&booking).Updates(map[string]interface{}{
			"status":        models.BookingActive,
			"check_in_time": now,
		})
		if err := clients.OccupyLocker(booking.LockerID); err != nil {
			log.Printf("[LOCKER] failed to mark occupied: %v", err)
		}
	}

	utils.SuccessResponse(c, 200, "Locker unlocked successfully", gin.H{
		"action":    "UNLOCK",
		"timestamp": now,
	})
}

func CheckIn(c *gin.Context) {
	userID := c.GetString("userId")
	bookingID := c.Param("id")

	var booking models.Booking
	if err := config.GetDB().First(&booking, "id = ? AND user_id = ?", bookingID, userID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}

	now := time.Now()
	config.GetDB().Model(&booking).Updates(map[string]interface{}{
		"status":        models.BookingActive,
		"check_in_time": now,
	})
	if err := clients.OccupyLocker(booking.LockerID); err != nil {
		log.Printf("[LOCKER] failed to mark occupied: %v", err)
	}

	utils.SuccessResponse(c, 200, "Check-in successful", gin.H{"checkInTime": now})
}

func CheckOut(c *gin.Context) {
	userID := c.GetString("userId")
	bookingID := c.Param("id")

	var booking models.Booking
	if err := config.GetDB().First(&booking, "id = ? AND user_id = ?", bookingID, userID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}

	now := time.Now()
	config.GetDB().Model(&booking).Updates(map[string]interface{}{
		"status":         models.BookingCompleted,
		"check_out_time": now,
	})
	if err := clients.ReleaseLocker(booking.LockerID); err != nil {
		log.Printf("[LOCKER] failed to release: %v", err)
	}

	utils.SuccessResponse(c, 200, "Check-out successful", gin.H{"checkOutTime": now})
}

func CancelBooking(c *gin.Context) {
	userID := c.GetString("userId")
	bookingID := c.Param("id")

	type CancelRequest struct {
		Reason string `json:"reason"`
	}
	var req CancelRequest
	_ = c.ShouldBindJSON(&req)

	var booking models.Booking
	if err := config.GetDB().First(&booking, "id = ? AND user_id = ?", bookingID, userID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}
	if booking.Status == models.BookingCompleted || booking.Status == models.BookingCancelled {
		utils.ErrorResponse(c, 400, "Booking cannot be cancelled")
		return
	}

	config.GetDB().Model(&booking).Updates(map[string]interface{}{
		"status":              models.BookingCancelled,
		"cancellation_reason": req.Reason,
	})
	if err := clients.ReleaseLocker(booking.LockerID); err != nil {
		log.Printf("[LOCKER] failed to release: %v", err)
	}

	utils.SuccessResponse(c, 200, "Booking cancelled successfully", nil)
}
