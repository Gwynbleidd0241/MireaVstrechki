package postgres

import (
	"database/sql"
	"time"
)

// ReminderEvent содержит данные события и список email-адресов всех участников.
type ReminderEvent struct {
	EventID    int64
	Title      string
	StartTime  time.Time
	Location   string
	MeetingURL string
	Emails     []string // email организатора + всех участников
}

type ReminderRepository struct {
	db *sql.DB
}

func NewReminderRepository(db *sql.DB) *ReminderRepository {
	return &ReminderRepository{db: db}
}

// FindUpcoming возвращает события, начинающиеся в окне [from, to],
// для которых напоминание ещё не было отправлено.
func (r *ReminderRepository) FindUpcoming(from, to time.Time) ([]ReminderEvent, error) {
	rows, err := r.db.Query(
		`SELECT DISTINCT e.id, e.title, e.start_time, e.location, e.meeting_url
		 FROM events e
		 WHERE e.start_time >= $1
		   AND e.start_time <  $2
		   AND e.status = 'scheduled'
		   AND NOT EXISTS (
		       SELECT 1 FROM event_reminders er WHERE er.event_id = e.id
		   )`,
		from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []ReminderEvent
	for rows.Next() {
		var ev ReminderEvent
		if err := rows.Scan(&ev.EventID, &ev.Title, &ev.StartTime, &ev.Location, &ev.MeetingURL); err != nil {
			return nil, err
		}
		events = append(events, ev)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Для каждого события собираем email-адреса: организатор + участники
	for i, ev := range events {
		emails, err := r.collectEmails(ev.EventID)
		if err != nil {
			return nil, err
		}
		events[i].Emails = emails
	}

	return events, nil
}

// collectEmails возвращает дедуплицированный список email-адресов
// организатора события и всех его участников.
func (r *ReminderRepository) collectEmails(eventID int64) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT DISTINCT u.email
		 FROM users u
		 WHERE u.id = (SELECT creator_id FROM events WHERE id = $1)
		    OR u.id IN (SELECT user_id FROM event_participants WHERE event_id = $1)`,
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, rows.Err()
}

// MarkSent фиксирует факт отправки напоминания для события.
// Повторный вызов для того же event_id игнорируется (ON CONFLICT DO NOTHING).
func (r *ReminderRepository) MarkSent(eventID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO event_reminders (event_id) VALUES ($1) ON CONFLICT DO NOTHING`,
		eventID,
	)
	return err
}
