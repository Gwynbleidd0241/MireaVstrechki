package service

import (
	"errors"

	"meeting-service/internal/model"
	"meeting-service/internal/repository/postgres"
)

var (
	ErrParticipantUserRequired = errors.New("participant user required")
	ErrInvalidParticipantRole  = errors.New("invalid participant role")
)

type ParticipantService struct {
	participantRepo *postgres.ParticipantRepository
}

func NewParticipantService(participantRepo *postgres.ParticipantRepository) *ParticipantService {
	return &ParticipantService{participantRepo: participantRepo}
}

type AddParticipantRequest struct {
	EventID int64
	UserID  int64
	Role    string
}

func (s *ParticipantService) Add(req AddParticipantRequest) (*model.Participant, error) {
	if req.UserID <= 0 {
		return nil, ErrParticipantUserRequired
	}

	role := req.Role
	if role == "" {
		role = "participant"
	}

	if role != "participant" && role != "responsible" {
		return nil, ErrInvalidParticipantRole
	}

	return s.participantRepo.Add(model.Participant{
		EventID: req.EventID,
		UserID:  req.UserID,
		Role:    role,
	})
}

func (s *ParticipantService) ListByEventID(eventID int64) ([]model.Participant, error) {
	return s.participantRepo.ListByEventID(eventID)
}
