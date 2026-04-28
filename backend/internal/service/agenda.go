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
	ErrAgendaItemNotFound           = errors.New("agenda item not found")
)

type AgendaService struct {
	agendaRepo *postgres.AgendaRepository
	eventRepo  *postgres.EventRepository
}

func NewAgendaService(agendaRepo *postgres.AgendaRepository, eventRepo *postgres.EventRepository) *AgendaService {
	return &AgendaService{
		agendaRepo: agendaRepo,
		eventRepo:  eventRepo,
	}
}

type AddAgendaItemRequest struct {
	EventID         int64
	Title           string
	Description     string
	DurationMinutes *int
	ActingUserID    int64
	ActingUserRole  string
}

type UpdateAgendaItemRequest struct {
	EventID         int64
	ItemID          int64
	Title           string
	Description     string
	DurationMinutes *int
	IsDone          bool
	ActingUserID    int64
	ActingUserRole  string
}

func (s *AgendaService) Add(req AddAgendaItemRequest) (*model.AgendaItem, error) {
	if err := s.checkManagePermission(req.EventID, req.ActingUserID, req.ActingUserRole); err != nil {
		return nil, err
	}

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

func (s *AgendaService) Update(req UpdateAgendaItemRequest) (*model.AgendaItem, error) {
	item, err := s.agendaRepo.GetByID(req.ItemID)
	if err != nil {
		return nil, err
	}

	if item == nil || item.EventID != req.EventID {
		return nil, ErrAgendaItemNotFound
	}

	if err := s.checkManagePermission(req.EventID, req.ActingUserID, req.ActingUserRole); err != nil {
		return nil, err
	}

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

	return s.agendaRepo.Update(model.AgendaItem{
		ID:              req.ItemID,
		Title:           title,
		Description:     strings.TrimSpace(req.Description),
		DurationMinutes: req.DurationMinutes,
		IsDone:          req.IsDone,
	})
}

func (s *AgendaService) Remove(eventID, itemID, actingUserID int64, actingUserRole string) error {
	item, err := s.agendaRepo.GetByID(itemID)
	if err != nil {
		return err
	}

	if item == nil || item.EventID != eventID {
		return ErrAgendaItemNotFound
	}

	if err := s.checkManagePermission(eventID, actingUserID, actingUserRole); err != nil {
		return err
	}

	return s.agendaRepo.Delete(itemID)
}

func (s *AgendaService) ListByEventID(eventID int64) ([]model.AgendaItem, error) {
	return s.agendaRepo.ListByEventID(eventID)
}

func (s *AgendaService) checkManagePermission(eventID, actingUserID int64, actingUserRole string) error {
	if actingUserRole == "admin" {
		return nil
	}

	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	if event == nil {
		return ErrEventNotFound
	}

	if event.CreatorID != actingUserID {
		return ErrPermissionDenied
	}

	return nil
}
