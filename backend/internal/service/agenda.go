package service

import (
	"errors"
	"strings"

	"meeting-service/internal/model"
	"meeting-service/internal/repository/postgres"
	"meeting-service/internal/validation"
)

var (
	ErrAgendaItemTitleRequired      = errors.New("agenda item title required")
	ErrAgendaItemTitleTooLong       = errors.New("agenda item title too long")
	ErrAgendaItemDescriptionTooLong = errors.New("agenda item description too long")
	ErrInvalidAgendaDuration        = errors.New("invalid agenda item duration")
)

type AgendaService struct {
	agendaRepo *postgres.AgendaRepository
}

func NewAgendaService(agendaRepo *postgres.AgendaRepository) *AgendaService {
	return &AgendaService{agendaRepo: agendaRepo}
}

type AddAgendaItemRequest struct {
	EventID         int64
	Title           string
	Description     string
	DurationMinutes *int
}

func (s *AgendaService) Add(req AddAgendaItemRequest) (*model.AgendaItem, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrAgendaItemTitleRequired
	}

	if len(title) > validation.MaxTitleLength {
		return nil, ErrAgendaItemTitleTooLong
	}

	if len(req.Description) > validation.MaxDescriptionLength {
		return nil, ErrAgendaItemDescriptionTooLong
	}

	if req.DurationMinutes != nil && *req.DurationMinutes < 0 {
		return nil, ErrInvalidAgendaDuration
	}

	return s.agendaRepo.Add(model.AgendaItem{
		EventID:         req.EventID,
		Title:           title,
		Description:     strings.TrimSpace(req.Description),
		DurationMinutes: req.DurationMinutes,
	})
}

func (s *AgendaService) ListByEventID(eventID int64) ([]model.AgendaItem, error) {
	return s.agendaRepo.ListByEventID(eventID)
}
