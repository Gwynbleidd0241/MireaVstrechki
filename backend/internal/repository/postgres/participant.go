package postgres

import (
	"database/sql"

	"meeting-service/internal/model"
)

type ParticipantRepository struct {
	db *sql.DB
}

func NewParticipantRepository(db *sql.DB) *ParticipantRepository {
	return &ParticipantRepository{db: db}
}

func (r *ParticipantRepository) Add(participant model.Participant) (*model.Participant, error) {
	var created model.Participant

	err := r.db.QueryRow(
		`INSERT INTO event_participants (event_id, user_id, role)
		 VALUES ($1, $2, $3)
		 RETURNING id, event_id, user_id, role, created_at`,
		participant.EventID,
		participant.UserID,
		participant.Role,
	).Scan(
		&created.ID,
		&created.EventID,
		&created.UserID,
		&created.Role,
		&created.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *ParticipantRepository) ListByEventID(eventID int64) ([]model.Participant, error) {
	rows, err := r.db.Query(
		`SELECT id, event_id, user_id, role, created_at
		 FROM event_participants
		 WHERE event_id = $1
		 ORDER BY created_at ASC`,
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []model.Participant

	for rows.Next() {
		var participant model.Participant

		if err := rows.Scan(
			&participant.ID,
			&participant.EventID,
			&participant.UserID,
			&participant.Role,
			&participant.CreatedAt,
		); err != nil {
			return nil, err
		}

		participants = append(participants, participant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}
