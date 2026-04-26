package postgres

import (
	"database/sql"

	"meeting-service/internal/model"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event model.Event) (*model.Event, error) {
	var created model.Event

	err := r.db.QueryRow(
		`INSERT INTO events (title, description, start_time, end_time, creator_id)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, title, description, start_time, end_time, creator_id, created_at`,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.CreatorID,
	).Scan(
		&created.ID,
		&created.Title,
		&created.Description,
		&created.StartTime,
		&created.EndTime,
		&created.CreatorID,
		&created.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *EventRepository) List() ([]model.Event, error) {
	rows, err := r.db.Query(
		`SELECT id, title, description, start_time, end_time, creator_id, created_at
		 FROM events
		 ORDER BY start_time ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event

	for rows.Next() {
		var event model.Event

		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.CreatorID,
			&event.CreatedAt,
		); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
