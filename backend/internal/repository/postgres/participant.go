package postgres

import (
	"database/sql"
	"errors"

	"meeting-service/internal/model"
)

type ParticipantRepository struct {
	db *sql.DB
}

func NewParticipantRepository(db *sql.DB) *ParticipantRepository {
	return &ParticipantRepository{db: db}
}

func (r *ParticipantRepository) GetByID(id int64) (*model.Participant, error) {
	var participant model.Participant

	err := r.db.QueryRow(
		`SELECT id, event_id, user_id, role, rsvp_status, created_at
		 FROM event_participants
		 WHERE id = $1`,
		id,
	).Scan(
		&participant.ID,
		&participant.EventID,
		&participant.UserID,
		&participant.Role,
		&participant.RSVPStatus,
		&participant.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &participant, nil
}

func (r *ParticipantRepository) UpdateRole(id int64, role string) (*model.Participant, error) {
	var updated model.Participant

	err := r.db.QueryRow(
		`UPDATE event_participants
		 SET role = $1
		 WHERE id = $2
		 RETURNING id, event_id, user_id, role, rsvp_status, created_at`,
		role,
		id,
	).Scan(
		&updated.ID,
		&updated.EventID,
		&updated.UserID,
		&updated.Role,
		&updated.RSVPStatus,
		&updated.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *ParticipantRepository) UpdateRSVP(id int64, rsvpStatus string) (*model.Participant, error) {
	var updated model.Participant

	err := r.db.QueryRow(
		`UPDATE event_participants
		 SET rsvp_status = $1
		 WHERE id = $2
		 RETURNING id, event_id, user_id, role, rsvp_status, created_at`,
		rsvpStatus,
		id,
	).Scan(
		&updated.ID,
		&updated.EventID,
		&updated.UserID,
		&updated.Role,
		&updated.RSVPStatus,
		&updated.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *ParticipantRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM event_participants WHERE id = $1`, id)
	return err
}

func (r *ParticipantRepository) Add(participant model.Participant) (*model.Participant, error) {
	var created model.Participant

	err := r.db.QueryRow(
		`INSERT INTO event_participants (event_id, user_id, role, rsvp_status)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, event_id, user_id, role, rsvp_status, created_at`,
		participant.EventID,
		participant.UserID,
		participant.Role,
		participant.RSVPStatus,
	).Scan(
		&created.ID,
		&created.EventID,
		&created.UserID,
		&created.Role,
		&created.RSVPStatus,
		&created.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *ParticipantRepository) ListByEventID(eventID int64) ([]model.Participant, error) {
	rows, err := r.db.Query(
		`SELECT id, event_id, user_id, role, rsvp_status, created_at
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
			&participant.RSVPStatus,
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
