// Package notifier orchestrates notification delivery. It depends only on the
// email.Sender and sms.Sender interfaces (Dependency Inversion), so providers
// can be swapped without changing this logic.
package notifier

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"notification-service/internal/email"
	"notification-service/internal/model"
	"notification-service/internal/sms"
)

// Service renders and dispatches booking notifications.
type Service struct {
	email email.Sender
	sms   sms.Sender
	tmpl  *template.Template
}

// New wires the orchestrator with its dependencies.
func New(emailSender email.Sender, smsSender sms.Sender) (*Service, error) {
	tmpl, err := template.New("booking-email").Parse(bookingEmailTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse email template: %w", err)
	}
	return &Service{email: emailSender, sms: smsSender, tmpl: tmpl}, nil
}

// HandleBookingConfirmed sends the confirmation email and (best effort) SMS.
func (s *Service) HandleBookingConfirmed(n model.BookingNotification) error {
	var body bytes.Buffer
	if err := s.tmpl.Execute(&body, n); err != nil {
		return fmt.Errorf("render email: %w", err)
	}

	subject := fmt.Sprintf("Booking Confirmed — %s", n.BookingNumber)
	if err := s.email.Send(n.UserEmail, subject, body.String()); err != nil {
		return err
	}
	log.Printf("[notifier] email sent to %s for %s", n.UserEmail, n.BookingNumber)

	// SMS is best-effort: a failure must not fail the whole notification.
	if n.UserPhone != "" {
		text := fmt.Sprintf(
			"Travel Easy: Booking %s confirmed at %s. Drop-off %s. Total INR %.2f. QR: %s",
			n.BookingNumber, n.LocationName, n.StartTime, n.TotalPrice, n.QRCode,
		)
		if err := s.sms.Send(n.UserPhone, text); err != nil {
			log.Printf("[notifier] sms failed for %s: %v", n.BookingNumber, err)
		} else {
			log.Printf("[notifier] sms sent to %s for %s", n.UserPhone, n.BookingNumber)
		}
	}
	return nil
}
