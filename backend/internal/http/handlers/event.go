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
	Title       string `json:"title"        example:"Планирование спринта"`
	Description string `json:"description"  example:"Расставляем приоритеты"`
	Status      string `json:"status"       example:"scheduled" enums:"scheduled,cancelled,completed"`
	Location    string `json:"location"     example:"Переговорная №2"`
	MeetingURL  string `json:"meeting_url"  example:"https://meet.google.com/abc-defg-hij"`
	StartTime   string `json:"start_time"   example:"2026-05-01T10:00:00Z"`
	EndTime     string `json:"end_time"     example:"2026-05-01T11:00:00Z"`
}

type updateEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"      enums:"scheduled,cancelled,completed"`
	Location    string `json:"location"`
	MeetingURL  string `json:"meeting_url"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type eventResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Location    string `json:"location"`
	MeetingURL  string `json:"meeting_url"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	CreatorID   int64  `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
}

// Create godoc
// @Summary  Создать встречу
// @Tags     events
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    input body     createEventRequest true "Параметры встречи"
// @Success  201   {object} eventResponse
// @Failure  400   {string} string "validation error"
// @Failure  403   {string} string "permission denied"
// @Router   /events [post]
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
		Status:      req.Status,
		Location:    req.Location,
		MeetingURL:  req.MeetingURL,
		StartTime:   startTime,
		EndTime:     endTime,
		CreatorID:   claims.UserID,
		CreatorRole: claims.Role,
	})
	if err != nil {
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusCreated, toEventResponse(*event))
}

// Get godoc
// @Summary  Получить встречу по id
// @Tags     events
// @Produce  json
// @Security BearerAuth
// @Param    id  path     int true "Event ID"
// @Success  200 {object} eventResponse
// @Failure  404 {string} string "event not found"
// @Router   /events/{id} [get]
func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseEventID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	event, err := h.eventService.Get(eventID)
	if err != nil {
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusOK, toEventResponse(*event))
}

// Update godoc
// @Summary  Обновить встречу
// @Tags     events
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id    path     int                true "Event ID"
// @Param    input body     updateEventRequest true "Новые параметры"
// @Success  200   {object} eventResponse
// @Failure  400   {string} string "validation error"
// @Failure  403   {string} string "permission denied"
// @Failure  404   {string} string "event not found"
// @Router   /events/{id} [patch]
func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, err := parseEventID(r.URL.Path)
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
		Status:      req.Status,
		Location:    req.Location,
		MeetingURL:  req.MeetingURL,
		StartTime:   startTime,
		EndTime:     endTime,
		UserID:      claims.UserID,
		UserRole:    claims.Role,
	})
	if err != nil {
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusOK, toEventResponse(*event))
}

// Delete godoc
// @Summary  Удалить встречу
// @Tags     events
// @Security BearerAuth
// @Param    id  path     int true "Event ID"
// @Success  204 {string} string "no content"
// @Failure  403 {string} string "permission denied"
// @Failure  404 {string} string "event not found"
// @Router   /events/{id} [delete]
func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, err := parseEventID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	if err := h.eventService.Delete(eventID, claims.UserID, claims.Role); err != nil {
		writeError(w, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List godoc
// @Summary  Список всех встреч
// @Tags     events
// @Produce  json
// @Security BearerAuth
// @Success  200 {array} eventResponse
// @Router   /events [get]
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

func toEventResponse(event model.Event) eventResponse {
	return eventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Status:      event.Status,
		Location:    event.Location,
		MeetingURL:  event.MeetingURL,
		StartTime:   event.StartTime.Format(time.RFC3339),
		EndTime:     event.EndTime.Format(time.RFC3339),
		CreatorID:   event.CreatorID,
		CreatedAt:   event.CreatedAt.Format(time.RFC3339),
	}
}
