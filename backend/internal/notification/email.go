package notification

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"

	"meeting-service/internal/config"
)

type EmailSender struct {
	cfg config.SMTP
}

func NewEmailSender(cfg config.SMTP) *EmailSender {
	return &EmailSender{cfg: cfg}
}

func (s *EmailSender) Enabled() bool {
	return s.cfg.Host != "" && s.cfg.Port != ""
}

func (s *EmailSender) Send(to []string, subject, body string) error {
	if !s.Enabled() {
		return nil
	}

	from := s.cfg.From
	if from == "" {
		from = s.cfg.Username
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain; charset=UTF-8", body)

	port := 587
	if s.cfg.Port == "465" {
		port = 465
	}

	d := gomail.NewDialer(s.cfg.Host, port, s.cfg.Username, s.cfg.Password)
	d.SSL = port == 465
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
