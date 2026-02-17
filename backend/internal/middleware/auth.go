package middleware

import (
	"strings"

	"locker-storage-api/config"
	"locker-storage-api/internal/models"
	"locker-storage-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, 401, "Authorization header required")
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.ErrorResponse(c, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		// Get user from database
		var user models.User
		if err := config.GetDB().First(&user, "id = ?", claims.UserID).Error; err != nil {
			utils.ErrorResponse(c, 401, "User not found")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
