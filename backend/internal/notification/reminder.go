package notification

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"meeting-service/internal/repository/postgres"
)

const (
	reminderCheckInterval = 1 * time.Minute
	reminderWindowStart   = 55 * time.Minute
	reminderWindowEnd     = 65 * time.Minute
)

type ReminderScheduler struct {
	repo   *postgres.ReminderRepository
	sender *EmailSender
	logger *zap.Logger
}

func NewReminderScheduler(
	repo *postgres.ReminderRepository,
	sender *EmailSender,
	logger *zap.Logger,
) *ReminderScheduler {
	return &ReminderScheduler{
		repo:   repo,
		sender: sender,
		logger: logger,
	}
}

func (s *ReminderScheduler) Start(ctx context.Context) {
	if !s.sender.Enabled() {
		s.logger.Info("reminder scheduler: SMTP not configured, reminders disabled")
		return
	}

	s.logger.Info("reminder scheduler started")

	ticker := time.NewTicker(reminderCheckInterval)
	defer ticker.Stop()

	s.run()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("reminder scheduler stopped")
			return
		case <-ticker.C:
			s.run()
		}
	}
}

func (s *ReminderScheduler) run() {
	now := time.Now().UTC()
	from := now.Add(reminderWindowStart)
	to := now.Add(reminderWindowEnd)

	events, err := s.repo.FindUpcoming(from, to)
	if err != nil {
		s.logger.Error("reminder: failed to query upcoming events", zap.Error(err))
		return
	}

	for _, ev := range events {
		if len(ev.Emails) == 0 {
			continue
		}

		subject := fmt.Sprintf("Напоминание: «%s» начинается через час", ev.Title)
		body := buildBody(ev)

		if err := s.sender.Send(ev.Emails, subject, body); err != nil {
			s.logger.Error("reminder: failed to send email",
				zap.Int64("event_id", ev.EventID),
				zap.Error(err),
			)
			continue
		}

		if err := s.repo.MarkSent(ev.EventID); err != nil {
			s.logger.Error("reminder: failed to mark as sent",
				zap.Int64("event_id", ev.EventID),
				zap.Error(err),
			)
		} else {
			s.logger.Info("reminder sent",
				zap.Int64("event_id", ev.EventID),
				zap.String("title", ev.Title),
				zap.Int("recipients", len(ev.Emails)),
			)
		}
	}
}

func buildBody(ev postgres.ReminderEvent) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Встреча «%s» начинается в %s.\n\n",
		ev.Title,
		ev.StartTime.Local().Format("02.01.2006 в 15:04"),
	))

	if ev.Location != "" {
		sb.WriteString(fmt.Sprintf("📍 Место: %s\n", ev.Location))
	}

	if ev.MeetingURL != "" {
		sb.WriteString(fmt.Sprintf("Ссылка: %s\n", ev.MeetingURL))
	}

	sb.WriteString("\nЭто автоматическое напоминание от Meeting Service.")

	return sb.String()
}
