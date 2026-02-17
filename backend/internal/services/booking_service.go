package services

import (
	"fmt"
	"math"
	"time"

	"locker-storage-api/internal/models"

	"github.com/google/uuid"
)

func GenerateBookingNumber() string {
	return fmt.Sprintf("BK%d%s", time.Now().Unix(), uuid.New().String()[:8])
}

func CalculateBookingPrice(rateType models.RateType, duration float64, hourlyRate, dailyRate float64) float64 {
	switch rateType {
	case models.RateHourly:
		return math.Ceil(duration) * hourlyRate
	case models.RateDaily:
		days := math.Ceil(duration / 24)
		return days * dailyRate
	case models.RateMultiDay:
		days := math.Ceil(duration / 24)
		if days > 1 {
			// 10% discount for multi-day bookings
			return (dailyRate * days) * 0.9
		}
		return dailyRate
	default:
		return duration * hourlyRate
	}
}

func GenerateQRCode() string {
	return uuid.New().String()
}
