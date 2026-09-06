package handlers

import (
	"strings"

	"review-service/config"
	"review-service/internal/clients"
	"review-service/internal/events"
	"review-service/internal/models"
	"review-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateReviewRequest struct {
	LocationID string `json:"locationId" binding:"required"`
	Rating     int    `json:"rating" binding:"required"`
	Comment    string `json:"comment"`
}

// GetLocationReviews returns reviews for a location (public).
// GET /api/v1/reviews?locationId=...
func GetLocationReviews(c *gin.Context) {
	locationID := c.Query("locationId")
	if locationID == "" {
		utils.ErrorResponse(c, 400, "locationId is required")
		return
	}

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
// POST /api/v1/reviews
func CreateReview(c *gin.Context) {
	userID := c.GetString("userId")

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}
	if req.Rating < 1 || req.Rating > 5 {
		utils.ErrorResponse(c, 400, "Rating must be between 1 and 5")
		return
	}

	// One review per user per location.
	var existing models.Review
	if err := config.GetDB().
		Where("user_id = ? AND location_id = ?", userID, req.LocationID).
		First(&existing).Error; err == nil {
		utils.ErrorResponse(c, 409, "You have already reviewed this location")
		return
	}

	// Best-effort display name from auth-service.
	userName := "Traveler"
	if u, err := clients.GetUser(userID); err == nil && u.Name != "" {
		userName = u.Name
	}

	review := models.Review{
		LocationID: req.LocationID,
		UserID:     userID,
		UserName:   userName,
		Rating:     req.Rating,
		Comment:    strings.TrimSpace(req.Comment),
	}
	if err := config.GetDB().Create(&review).Error; err != nil {
		utils.ErrorResponse(c, 500, "Failed to create review")
		return
	}

	// Recompute aggregate for this location.
	var stats struct {
		Avg   float64
		Count int64
	}
	config.GetDB().Model(&models.Review{}).
		Select("COALESCE(AVG(rating), 0) as avg, COUNT(*) as count").
		Where("location_id = ?", req.LocationID).
		Scan(&stats)

	// Notify location-service to update its rating (async).
	events.PublishReviewCreated(events.ReviewCreatedEvent{
		LocationID:    req.LocationID,
		AverageRating: stats.Avg,
		TotalReviews:  int(stats.Count),
	})

	utils.SuccessResponse(c, 201, "Review submitted", gin.H{
		"review": review,
	})
}
