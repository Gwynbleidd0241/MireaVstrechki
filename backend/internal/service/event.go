package service

import (
	"errors"
	"strings"
	"time"

	"meeting-service/internal/model"
	"meeting-service/internal/repository/postgres"
	"meeting-service/internal/validation"
)

var (
	ErrEventTitleRequired      = errors.New("event title required")
	ErrEventTitleTooLong       = errors.New("event title too long")
	ErrEventDescriptionTooLong = errors.New("event description too long")
	ErrInvalidEventTime        = errors.New("invalid event time")
	ErrEventNotFound           = errors.New("event not found")
	ErrPermissionDenied        = errors.New("permission denied")
)

type EventService struct {
	eventRepo *postgres.EventRepository
}

func NewEventService(eventRepo *postgres.EventRepository) *EventService {
	return &EventService{eventRepo: eventRepo}
}

type CreateEventRequest struct {
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CreatorID   int64
	CreatorRole string
}

type UpdateEventRequest struct {
	EventID     int64
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	UserID      int64
	UserRole    string
}

func (s *EventService) Create(req CreateEventRequest) (*model.Event, error) {
	if req.CreatorRole != "admin" && req.CreatorRole != "organizer" {
		return nil, ErrPermissionDenied
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrEventTitleRequired
	}

	if len(title) > validation.MaxTitleLength {
		return nil, ErrEventTitleTooLong
	}

	if len(req.Description) > validation.MaxDescriptionLength {
		return nil, ErrEventDescriptionTooLong
	}

	if req.StartTime.IsZero() || req.EndTime.IsZero() || !req.EndTime.After(req.StartTime) {
		return nil, ErrInvalidEventTime
	}

	return s.eventRepo.Create(model.Event{
		Title:       title,
		Description: strings.TrimSpace(req.Description),
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		CreatorID:   req.CreatorID,
	})
}

func (s *EventService) Get(eventID int64) (*model.Event, error) {
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, ErrEventNotFound
	}

	return event, nil
}

func (s *EventService) Update(req UpdateEventRequest) (*model.Event, error) {
	event, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, ErrEventNotFound
	}

	if req.UserRole != "admin" && event.CreatorID != req.UserID {
		return nil, ErrPermissionDenied
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrEventTitleRequired
	}

	if len(title) > validation.MaxTitleLength {
		return nil, ErrEventTitleTooLong
	}

	if len(req.Description) > validation.MaxDescriptionLength {
		return nil, ErrEventDescriptionTooLong
	}

	if req.StartTime.IsZero() || req.EndTime.IsZero() || !req.EndTime.After(req.StartTime) {
		return nil, ErrInvalidEventTime
	}

	return s.eventRepo.Update(model.Event{
		ID:          req.EventID,
		Title:       title,
		Description: strings.TrimSpace(req.Description),
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	})
}

func (s *EventService) Delete(eventID, userID int64, userRole string) error {
	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	if event == nil {
		return ErrEventNotFound
	}

	if userRole != "admin" && event.CreatorID != userID {
		return ErrPermissionDenied
	}

	return s.eventRepo.Delete(eventID)
}

func (s *EventService) List() ([]model.Event, error) {
	return s.eventRepo.List()
}

func (s *EventService) ListForUser(userID int64) ([]model.Event, error) {
	return s.eventRepo.ListForUser(userID)
}
