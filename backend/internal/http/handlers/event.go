package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"meeting-service/internal/auth"
	"meeting-service/internal/http/middleware"
	"meeting-service/internal/model"
	"meeting-service/internal/service"
)

type EventHandler struct {
	eventService *service.EventService
	logger       *zap.Logger
}

func NewEventHandler(eventService *service.EventService, logger *zap.Logger) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		logger:       logger,
	}
}

type createEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type updateEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type eventResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	CreatorID   int64  `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req createEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "invalid start_time", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "invalid end_time", http.StatusBadRequest)
		return
	}

	event, err := h.eventService.Create(service.CreateEventRequest{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		CreatorID:   claims.UserID,
		CreatorRole: claims.Role,
	})
	if err != nil {
		h.handleEventError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toEventResponse(*event))
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	eventID, err := eventIDFromEventPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	event, err := h.eventService.Get(eventID)
	if err != nil {
		h.handleEventError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toEventResponse(*event))
}

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, err := eventIDFromEventPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	var req updateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "invalid start_time", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "invalid end_time", http.StatusBadRequest)
		return
	}

	event, err := h.eventService.Update(service.UpdateEventRequest{
		EventID:     eventID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		UserID:      claims.UserID,
		UserRole:    claims.Role,
	})
	if err != nil {
		h.handleEventError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toEventResponse(*event))
}

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, err := eventIDFromEventPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	if err := h.eventService.Delete(eventID, claims.UserID, claims.Role); err != nil {
		h.handleEventError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventService.List()
	if err != nil {
		h.logger.Error("failed to list events", zap.Error(err))
		http.Error(w, "failed to list events", http.StatusInternalServerError)
		return
	}

	resp := make([]eventResponse, 0, len(events))
	for _, event := range events {
		resp = append(resp, toEventResponse(event))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *EventHandler) handleEventError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrEventTitleRequired:
		http.Error(w, "event title required", http.StatusBadRequest)
	case service.ErrEventTitleTooLong:
		http.Error(w, "event title too long", http.StatusBadRequest)
	case service.ErrEventDescriptionTooLong:
		http.Error(w, "event description too long", http.StatusBadRequest)
	case service.ErrInvalidEventTime:
		http.Error(w, "invalid event time", http.StatusBadRequest)
	case service.ErrEventNotFound:
		http.Error(w, "event not found", http.StatusNotFound)
	case service.ErrPermissionDenied:
		http.Error(w, "permission denied", http.StatusForbidden)
	default:
		h.logger.Error("event handler error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func eventIDFromEventPath(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 || parts[0] != "events" {
		return 0, strconv.ErrSyntax
	}

	return strconv.ParseInt(parts[1], 10, 64)
}

func toEventResponse(event model.Event) eventResponse {
	return eventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartTime:   event.StartTime.Format(time.RFC3339),
		EndTime:     event.EndTime.Format(time.RFC3339),
		CreatorID:   event.CreatorID,
		CreatedAt:   event.CreatedAt.Format(time.RFC3339),
	}
}
