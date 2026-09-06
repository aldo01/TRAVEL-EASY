package handlers

import (
	"auth-service/config"
	"auth-service/internal/models"
	"auth-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetUserInternal (service-to-service) returns a user by id.
func GetUserInternal(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.GetDB().First(&user, "id = ?", id).Error; err != nil {
		utils.ErrorResponse(c, 404, "User not found")
		return
	}

	utils.SuccessResponse(c, 200, "User retrieved", user)
}
