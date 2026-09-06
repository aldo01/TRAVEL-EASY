package email

import "notification-service/internal/brevo"

const brevoEmailURL = "https://api.brevo.com/v3/smtp/email"

// BrevoSender delivers email via the Brevo transactional email API.
type BrevoSender struct {
	APIKey      string
	SenderName  string
	SenderEmail string
}

func (s BrevoSender) Send(to, subject, htmlBody string) error {
	payload := map[string]interface{}{
		"sender":      map[string]string{"name": s.SenderName, "email": s.SenderEmail},
		"to":          []map[string]string{{"email": to}},
		"subject":     subject,
		"htmlContent": htmlBody,
	}
	return brevo.Post(s.APIKey, brevoEmailURL, payload)
}
