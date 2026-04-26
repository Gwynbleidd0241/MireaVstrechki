package service

import (
	"errors"
	"strings"
	"time"

	"meeting-service/internal/model"
	"meeting-service/internal/repository/postgres"
)

var (
	ErrEventTitleRequired = errors.New("event title required")
	ErrInvalidEventTime   = errors.New("invalid event time")
	ErrPermissionDenied   = errors.New("permission denied")
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

func (s *EventService) Create(req CreateEventRequest) (*model.Event, error) {
	if req.CreatorRole != "admin" && req.CreatorRole != "organizer" {
		return nil, ErrPermissionDenied
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrEventTitleRequired
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

func (s *EventService) List() ([]model.Event, error) {
	return s.eventRepo.List()
}
