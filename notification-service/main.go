// notification-service — sends booking confirmation email + SMS.
//
// This main is the composition root: it selects concrete providers from
// configuration and injects them into the notifier service (Dependency
// Inversion). Adding a new provider means adding a type that satisfies
// email.Sender / sms.Sender and one branch here — no other code changes
// (Open/Closed).
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"notification-service/internal/email"
	"notification-service/internal/notifier"
	"notification-service/internal/sms"
	"notification-service/internal/transport"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("GIN_MODE") == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	svc, err := notifier.New(buildEmailSender(), buildSMSSender())
	if err != nil {
		log.Fatal("failed to init notifier:", err)
	}

	// Consume async booking.confirmed events.
	go transport.StartNATSSubscriber(svc)

	router := gin.Default()
	transport.NewHTTPHandler(svc).Register(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	log.Printf("notification-service running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("failed to start notification service:", err)
	}
}

// buildEmailSender selects the email provider from configuration.
// Priority: Brevo > SMTP > Resend > log-only.
func buildEmailSender() email.Sender {
	if key := os.Getenv("BREVO_API_KEY"); key != "" && os.Getenv("BREVO_SENDER_EMAIL") != "" {
		name := os.Getenv("BREVO_SENDER_NAME")
		if name == "" {
			name = "Travel Easy"
		}
		log.Println("[config] email provider: Brevo")
		return email.BrevoSender{APIKey: key, SenderName: name, SenderEmail: os.Getenv("BREVO_SENDER_EMAIL")}
	}
	if os.Getenv("SMTP_USER") != "" && os.Getenv("SMTP_PASS") != "" {
		log.Println("[config] email provider: SMTP")
		return email.SMTPSender{
			Host: os.Getenv("SMTP_HOST"),
			Port: os.Getenv("SMTP_PORT"),
			User: os.Getenv("SMTP_USER"),
			Pass: os.Getenv("SMTP_PASS"),
			From: os.Getenv("SMTP_FROM"),
		}
	}
	if key := os.Getenv("RESEND_API_KEY"); key != "" {
		log.Println("[config] email provider: Resend")
		return email.ResendSender{APIKey: key, From: os.Getenv("RESEND_FROM")}
	}
	log.Println("[config] email provider: none (log only)")
	return email.LogSender{}
}

// buildSMSSender selects the SMS provider from configuration.
func buildSMSSender() sms.Sender {
	if key := os.Getenv("BREVO_API_KEY"); key != "" {
		log.Println("[config] sms provider: Brevo")
		return sms.BrevoSender{APIKey: key, SenderID: os.Getenv("BREVO_SMS_SENDER")}
	}
	log.Println("[config] sms provider: none")
	return sms.NoopSender{}
}
