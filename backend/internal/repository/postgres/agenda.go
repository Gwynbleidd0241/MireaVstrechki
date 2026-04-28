package postgres

import (
	"database/sql"

	"meeting-service/internal/model"
)

type AgendaRepository struct {
	db *sql.DB
}

func NewAgendaRepository(db *sql.DB) *AgendaRepository {
	return &AgendaRepository{db: db}
}

func (r *AgendaRepository) Add(item model.AgendaItem) (*model.AgendaItem, error) {
	var created model.AgendaItem

	err := r.db.QueryRow(
		`INSERT INTO agenda_items (event_id, position, title, description, duration_minutes)
		 VALUES (
		   $1,
		   COALESCE((SELECT MAX(position) FROM agenda_items WHERE event_id = $1), 0) + 1,
		   $2, $3, $4
		 )
		 RETURNING id, event_id, position, title, description, duration_minutes, is_done, created_at`,
		item.EventID,
		item.Title,
		item.Description,
		item.DurationMinutes,
	).Scan(
		&created.ID,
		&created.EventID,
		&created.Position,
		&created.Title,
		&created.Description,
		&created.DurationMinutes,
		&created.IsDone,
		&created.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *AgendaRepository) ListByEventID(eventID int64) ([]model.AgendaItem, error) {
	rows, err := r.db.Query(
		`SELECT id, event_id, position, title, description, duration_minutes, is_done, created_at
		 FROM agenda_items
		 WHERE event_id = $1
		 ORDER BY position ASC`,
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.AgendaItem

	for rows.Next() {
		var item model.AgendaItem

		if err := rows.Scan(
			&item.ID,
			&item.EventID,
			&item.Position,
			&item.Title,
			&item.Description,
			&item.DurationMinutes,
			&item.IsDone,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
