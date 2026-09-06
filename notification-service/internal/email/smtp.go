package email

import (
	"bytes"
	"fmt"
	"net/smtp"
)

// SMTPSender delivers email over SMTP (STARTTLS on 587), e.g. Gmail app password.
type SMTPSender struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func (s SMTPSender) Send(to, subject, htmlBody string) error {
	host := s.Host
	if host == "" {
		host = "smtp.gmail.com"
	}
	port := s.Port
	if port == "" {
		port = "587"
	}
	from := s.From
	if from == "" {
		from = s.User
	}

	var msg bytes.Buffer
	fmt.Fprintf(&msg, "From: %s\r\n", from)
	fmt.Fprintf(&msg, "To: %s\r\n", to)
	fmt.Fprintf(&msg, "Subject: %s\r\n", subject)
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	auth := smtp.PlainAuth("", s.User, s.Pass, host)
	return smtp.SendMail(host+":"+port, auth, s.User, []string{to}, msg.Bytes())
}
