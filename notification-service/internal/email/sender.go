// Package email defines the email delivery abstraction and its providers.
package email

import "log"

// Sender delivers a single HTML email. Implementations are interchangeable
// (Liskov), so the notifier depends only on this interface (DIP).
type Sender interface {
	Send(to, subject, htmlBody string) error
}

// LogSender is the fallback used when no provider is configured; it logs
// instead of sending, so the app still runs in local/dev environments.
type LogSender struct{}

func (LogSender) Send(to, subject, _ string) error {
	log.Printf("[email] no provider configured — logged only. to=%s subject=%q", to, subject)
	return nil
}
