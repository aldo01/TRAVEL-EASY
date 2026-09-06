package sms

import "notification-service/internal/brevo"

const brevoSMSURL = "https://api.brevo.com/v3/transactionalSMS/sms"

// BrevoSender delivers transactional SMS via the Brevo API.
type BrevoSender struct {
	APIKey   string
	SenderID string
}

func (s BrevoSender) Send(to, text string) error {
	if to == "" {
		return nil
	}
	sender := s.SenderID
	if sender == "" {
		sender = "TravelEasy"
	}
	payload := map[string]interface{}{
		"type":      "transactional",
		"sender":    sender,
		"recipient": to,
		"content":   text,
	}
	return brevo.Post(s.APIKey, brevoSMSURL, payload)
}
