package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const resendURL = "https://api.resend.com/emails"

// ResendSender delivers email via the Resend API.
type ResendSender struct {
	APIKey string
	From   string
}

func (s ResendSender) Send(to, subject, htmlBody string) error {
	from := s.From
	if from == "" {
		from = "Travel Easy <onboarding@resend.dev>"
	}

	payload := map[string]interface{}{
		"from":    from,
		"to":      []string{to},
		"subject": subject,
		"html":    htmlBody,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resendURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend API error (status %d): %s", resp.StatusCode, string(respBody))
	}
	return nil
}
