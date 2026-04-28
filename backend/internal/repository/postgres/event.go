package postgres

import (
	"database/sql"
	"errors"

	"meeting-service/internal/model"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetByID(id int64) (*model.Event, error) {
	var event model.Event

	err := r.db.QueryRow(
		`SELECT id, title, description, start_time, end_time, creator_id, created_at
		 FROM events
		 WHERE id = $1`,
		id,
	).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.CreatorID,
		&event.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventRepository) Update(event model.Event) (*model.Event, error) {
	var updated model.Event

	err := r.db.QueryRow(
		`UPDATE events
		 SET title = $1, description = $2, start_time = $3, end_time = $4
		 WHERE id = $5
		 RETURNING id, title, description, start_time, end_time, creator_id, created_at`,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
		event.ID,
	).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Description,
		&updated.StartTime,
		&updated.EndTime,
		&updated.CreatorID,
		&updated.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *EventRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM events WHERE id = $1`, id)
	return err
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

func (r *EventRepository) ListForUser(userID int64) ([]model.Event, error) {
	rows, err := r.db.Query(
		`SELECT DISTINCT e.id, e.title, e.description, e.start_time, e.end_time, e.creator_id, e.created_at
		 FROM events e
		 LEFT JOIN event_participants p ON p.event_id = e.id
		 WHERE e.creator_id = $1 OR p.user_id = $1
		 ORDER BY e.start_time ASC`,
		userID,
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
