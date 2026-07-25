package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type BookingNotification struct {
	BookingNumber string  `json:"bookingNumber" binding:"required"`
	UserEmail     string  `json:"userEmail" binding:"required,email"`
	UserName      string  `json:"userName" binding:"required"`
	LocationName  string  `json:"locationName" binding:"required"`
	LocationAddr  string  `json:"locationAddress"`
	StartTime     string  `json:"startTime" binding:"required"`
	EndTime       string  `json:"endTime"`
	Bags          int     `json:"bags"`
	RateType      string  `json:"rateType"`
	TotalPrice    float64 `json:"totalPrice"`
	QRCode        string  `json:"qrCode"`
}

const emailTemplate = `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; color: #333;">
  <div style="text-align: center; padding: 20px 0; border-bottom: 2px solid #2563eb;">
    <h1 style="color: #2563eb; margin: 0;">Travel Easy</h1>
    <p style="color: #666; margin: 4px 0 0;">Luggage storage & locker booking</p>
  </div>

  <div style="padding: 24px 0;">
    <h2 style="color: #111; margin: 0 0 4px;">Booking Confirmed!</h2>
    <p style="color: #666; margin: 0;">Hi {{.UserName}}, your luggage storage booking is confirmed.</p>
  </div>

  <div style="background: #f9fafb; border: 1px solid #e5e7eb; border-radius: 8px; padding: 20px; margin-bottom: 20px;">
    <table style="width: 100%; border-collapse: collapse;">
      <tr>
        <td style="padding: 8px 0; color: #666; width: 40%;">Booking ID</td>
        <td style="padding: 8px 0; font-weight: bold; color: #111;">{{.BookingNumber}}</td>
      </tr>
      <tr>
        <td style="padding: 8px 0; color: #666;">Location</td>
        <td style="padding: 8px 0; color: #111;">{{.LocationName}}</td>
      </tr>
      {{if .LocationAddr}}
      <tr>
        <td style="padding: 8px 0; color: #666;">Address</td>
        <td style="padding: 8px 0; color: #111;">{{.LocationAddr}}</td>
      </tr>
      {{end}}
      <tr>
        <td style="padding: 8px 0; color: #666;">Drop-off</td>
        <td style="padding: 8px 0; color: #111;">{{.StartTime}}</td>
      </tr>
      {{if .EndTime}}
      <tr>
        <td style="padding: 8px 0; color: #666;">Pick-up</td>
        <td style="padding: 8px 0; color: #111;">{{.EndTime}}</td>
      </tr>
      {{end}}
      <tr>
        <td style="padding: 8px 0; color: #666;">Bags</td>
        <td style="padding: 8px 0; color: #111;">{{.Bags}}</td>
      </tr>
      <tr>
        <td style="padding: 8px 0; color: #666;">Rate type</td>
        <td style="padding: 8px 0; color: #111;">{{.RateType}}</td>
      </tr>
      <tr style="border-top: 1px solid #e5e7eb;">
        <td style="padding: 12px 0 8px; font-weight: bold; color: #111;">Total</td>
        <td style="padding: 12px 0 8px; font-weight: bold; color: #111; font-size: 18px;">₹{{printf "%.2f" .TotalPrice}}</td>
      </tr>
    </table>
  </div>

  {{if .QRCode}}
  <div style="background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 8px; padding: 16px; margin-bottom: 20px; text-align: center;">
    <p style="margin: 0 0 4px; font-weight: bold; color: #1e40af;">Your QR Code</p>
    <p style="margin: 0; color: #3b82f6; font-family: monospace; font-size: 14px;">{{.QRCode}}</p>
    <p style="margin: 8px 0 0; color: #666; font-size: 12px;">Show this code at the location to unlock your locker.</p>
  </div>
  {{end}}

  <div style="padding: 16px 0; border-top: 1px solid #e5e7eb; color: #999; font-size: 12px; text-align: center;">
    <p>Thank you for using Travel Easy.</p>
    <p>If you need to cancel, please do so before midnight the day before your drop-off.</p>
  </div>
</body>
</html>`

func sendEmail(to, subject, htmlBody string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	fromEmail := os.Getenv("RESEND_FROM")

	if apiKey == "" {
		log.Printf("[NOTIFICATION] RESEND_API_KEY not set — logging email instead")
		log.Printf("[NOTIFICATION] To: %s | Subject: %s", to, subject)
		return nil
	}

	if fromEmail == "" {
		fromEmail = "Travel Easy <onboarding@resend.dev>"
	}

	payload := map[string]interface{}{
		"from":    fromEmail,
		"to":      []string{to},
		"subject": subject,
		"html":    htmlBody,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	log.Printf("[NOTIFICATION] Resend response: %s", string(respBody))
	return nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "notification-service"})
	})

	router.POST("/notify/booking-confirmed", func(c *gin.Context) {
		var req BookingNotification
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"success": false, "error": err.Error()})
			return
		}

		tmpl, err := template.New("email").Parse(emailTemplate)
		if err != nil {
			log.Printf("[NOTIFICATION] Template error: %v", err)
			c.JSON(500, gin.H{"success": false, "error": "Failed to render email"})
			return
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, req); err != nil {
			log.Printf("[NOTIFICATION] Template exec error: %v", err)
			c.JSON(500, gin.H{"success": false, "error": "Failed to render email"})
			return
		}

		subject := fmt.Sprintf("Booking Confirmed — %s", req.BookingNumber)

		if err := sendEmail(req.UserEmail, subject, body.String()); err != nil {
			log.Printf("[NOTIFICATION] Email send error: %v", err)
			c.JSON(500, gin.H{"success": false, "error": "Failed to send email"})
			return
		}

		log.Printf("[NOTIFICATION] Booking confirmation sent to %s for %s", req.UserEmail, req.BookingNumber)
		c.JSON(200, gin.H{"success": true, "message": "Notification sent"})
	})

	log.Printf("Notification service running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start notification service:", err)
	}
}
