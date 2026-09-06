// Package brevo is a thin client for the Brevo (Sendinblue) REST API,
// shared by the email and sms sender implementations.
package brevo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Post sends an authenticated JSON POST to a Brevo endpoint.
func Post(apiKey, url string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("brevo API error (status %d): %s", resp.StatusCode, string(respBody))
	}
	return nil
}
