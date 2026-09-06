package handlers

import (
	"strings"

	"locker-storage-api/config"
	"locker-storage-api/internal/models"
	"locker-storage-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required"`
	Comment string `json:"comment"`
}

// GetLocationReviews returns all reviews for a location (public).
func GetLocationReviews(c *gin.Context) {
	locationID := c.Param("id")

	var reviews []models.Review
	if err := config.GetDB().
		Where("location_id = ?", locationID).
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch reviews")
		return
	}

	utils.SuccessResponse(c, 200, "Reviews retrieved", gin.H{
		"count":   len(reviews),
		"reviews": reviews,
	})
}

// CreateReview lets an authenticated user review a location once (protected).
func CreateReview(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	locationID := c.Param("id")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		utils.ErrorResponse(c, 400, "Rating must be between 1 and 5")
		return
	}

	// Location must exist.
	var location models.Location
	if err := config.GetDB().First(&location, "id = ?", locationID).Error; err != nil {
		utils.ErrorResponse(c, 404, "Location not found")
		return
	}

	// Require the user to have a booking at this location.
	var bookingCount int64
	config.GetDB().Model(&models.Booking{}).
		Where("user_id = ? AND location_id = ?", user.ID, locationID).
		Count(&bookingCount)
	if bookingCount == 0 {
		utils.ErrorResponse(c, 403, "You can only review locations you have booked")
		return
	}

	// Prevent duplicate reviews from the same user.
	var existing models.Review
	if err := config.GetDB().
		Where("user_id = ? AND location_id = ?", user.ID, locationID).
		First(&existing).Error; err == nil {
		utils.ErrorResponse(c, 409, "You have already reviewed this location")
		return
	}

	review := models.Review{
		LocationID: locationID,
		UserID:     user.ID,
		UserName:   user.Name,
		Rating:     req.Rating,
		Comment:    strings.TrimSpace(req.Comment),
	}

	tx := config.GetDB().Begin()
	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, 500, "Failed to create review")
		return
	}

	// Recompute the location's rating aggregate.
	var stats struct {
		Avg   float64
		Count int64
	}
	tx.Model(&models.Review{}).
		Select("COALESCE(AVG(rating), 0) as avg, COUNT(*) as count").
		Where("location_id = ?", locationID).
		Scan(&stats)

	if err := tx.Model(&models.Location{}).
		Where("id = ?", locationID).
		Updates(map[string]interface{}{
			"rating":        stats.Avg,
			"total_reviews": stats.Count,
		}).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(c, 500, "Failed to update location rating")
		return
	}

	tx.Commit()

	utils.SuccessResponse(c, 201, "Review submitted", gin.H{
		"review": review,
	})
}
