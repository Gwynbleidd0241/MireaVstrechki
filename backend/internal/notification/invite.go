package notification

import (
	"fmt"
	"time"
)

type InviteData struct {
	RecipientEmail string
	EventTitle     string
	StartTime      time.Time
	Location       string
	MeetingURL     string
}

func (s *EmailSender) SendInvite(data InviteData) error {
	subject := fmt.Sprintf("Вас добавили во встречу «%s»", data.EventTitle)
	body := buildInviteBody(data)
	return s.Send([]string{data.RecipientEmail}, subject, body)
}

func buildInviteBody(data InviteData) string {
	body := fmt.Sprintf(
		"Вы добавлены как участник встречи «%s».\n\nДата и время: %s",
		data.EventTitle,
		data.StartTime.Local().Format("02.01.2006 в 15:04"),
	)

	if data.Location != "" {
		body += fmt.Sprintf("\n📍 Место: %s", data.Location)
	}

	if data.MeetingURL != "" {
		body += fmt.Sprintf("\n Ссылка: %s", data.MeetingURL)
	}

	body += "\n\nЭто автоматическое уведомление от MireaVstrechki."
	return body
}
