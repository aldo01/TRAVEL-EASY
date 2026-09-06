// Package sms defines the SMS delivery abstraction and its providers.
package sms

import "log"

// Sender delivers a single SMS. Implementations are interchangeable.
type Sender interface {
	Send(to, text string) error
}

// NoopSender is used when no SMS provider is configured.
type NoopSender struct{}

func (NoopSender) Send(to, _ string) error {
	if to != "" {
		log.Printf("[sms] no provider configured — skipped. to=%s", to)
	}
	return nil
}
