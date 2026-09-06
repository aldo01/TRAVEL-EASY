// Package transport adapts external inputs (HTTP, NATS) to the notifier service.
package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"notification-service/internal/model"
	"notification-service/internal/notifier"
)

// HTTPHandler exposes the notifier over HTTP.
type HTTPHandler struct {
	svc *notifier.Service
}

func NewHTTPHandler(svc *notifier.Service) *HTTPHandler {
	return &HTTPHandler{svc: svc}
}

// Register wires the routes onto the given engine.
func (h *HTTPHandler) Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "notification-service"})
	})
	r.POST("/notify/booking-confirmed", h.bookingConfirmed)
}

func (h *HTTPHandler) bookingConfirmed(c *gin.Context) {
	var n model.BookingNotification
	if err := c.ShouldBindJSON(&n); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	if err := h.svc.HandleBookingConfirmed(n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to send notification"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Notification sent"})
}
