package middleware

import (
	"review-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireUser ensures the gateway forwarded an authenticated user id.
func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-Id")
		if userID == "" {
			utils.ErrorResponse(c, 401, "Unauthorized")
			c.Abort()
			return
		}
		c.Set("userId", userID)
		c.Next()
	}
}
