package notification

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

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

	header := strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", strings.Join(to, ", ")),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	if s.cfg.Port == "465" {
		return s.sendTLS(addr, auth, from, to, []byte(header))
	}

	return smtp.SendMail(addr, auth, from, to, []byte(header))
}

func (s *EmailSender) sendTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	tlsCfg := &tls.Config{
		ServerName: s.cfg.Host,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail: %w", err)
	}
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return fmt.Errorf("smtp rcpt %s: %w", addr, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	return w.Close()
}
