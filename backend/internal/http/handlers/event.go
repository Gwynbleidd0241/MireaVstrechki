package handlers

import (
	"encoding/json"
	"net/http"
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
	case service.ErrInvalidEventTime:
		http.Error(w, "invalid event time", http.StatusBadRequest)
	case service.ErrPermissionDenied:
		http.Error(w, "permission denied", http.StatusForbidden)
	default:
		h.logger.Error("event handler error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
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
