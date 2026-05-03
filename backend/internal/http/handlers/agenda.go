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

type AgendaHandler struct {
	agendaService *service.AgendaService
	logger        *zap.Logger
}

func NewAgendaHandler(agendaService *service.AgendaService, logger *zap.Logger) *AgendaHandler {
	return &AgendaHandler{
		agendaService: agendaService,
		logger:        logger,
	}
}

type addAgendaItemRequest struct {
	Title           string `json:"title"            example:"Цели спринта"`
	Description     string `json:"description"`
	DurationMinutes *int   `json:"duration_minutes" example:"15"`
}

type updateAgendaItemRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	DurationMinutes *int   `json:"duration_minutes"`
	IsDone          bool   `json:"is_done"`
}

type agendaItemResponse struct {
	ID              int64  `json:"id"`
	EventID         int64  `json:"event_id"`
	Position        int    `json:"position"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	DurationMinutes *int   `json:"duration_minutes"`
	IsDone          bool   `json:"is_done"`
	CreatedAt       string `json:"created_at"`
}

// Add
// @Summary  Добавить пункт повестки
// @Tags     agenda
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id    path     int                  true "Event ID"
// @Param    input body     addAgendaItemRequest true "Пункт повестки"
// @Success  201   {object} agendaItemResponse
// @Failure  400   {string} string "validation error"
// @Failure  403   {string} string "permission denied"
// @Router   /events/{id}/agenda [post]
func (h *AgendaHandler) Add(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, err := parseEventResource(r.URL.Path, "agenda")
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	var req addAgendaItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	item, err := h.agendaService.Add(service.AddAgendaItemRequest{
		EventID:         eventID,
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		ActingUserID:    claims.UserID,
		ActingUserRole:  claims.Role,
	})
	if err != nil {
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusCreated, toAgendaItemResponse(*item))
}

// Update
// @Summary  Обновить пункт повестки
// @Tags     agenda
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id       path     int                     true "Event ID"
// @Param    agendaId path     int                     true "Agenda Item ID"
// @Param    input    body     updateAgendaItemRequest true "Новые поля"
// @Success  200      {object} agendaItemResponse
// @Failure  403      {string} string "permission denied"
// @Failure  404      {string} string "agenda item not found"
// @Router   /events/{id}/agenda/{agendaId} [patch]
func (h *AgendaHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, itemID, err := parseEventSubResource(r.URL.Path, "agenda")
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	var req updateAgendaItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	item, err := h.agendaService.Update(service.UpdateAgendaItemRequest{
		EventID:         eventID,
		ItemID:          itemID,
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		IsDone:          req.IsDone,
		ActingUserID:    claims.UserID,
		ActingUserRole:  claims.Role,
	})
	if err != nil {
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusOK, toAgendaItemResponse(*item))
}

// Delete
// @Summary  Удалить пункт повестки
// @Tags     agenda
// @Security BearerAuth
// @Param    id       path     int true "Event ID"
// @Param    agendaId path     int true "Agenda Item ID"
// @Success  204      {string} string "no content"
// @Failure  403      {string} string "permission denied"
// @Failure  404      {string} string "agenda item not found"
// @Router   /events/{id}/agenda/{agendaId} [delete]
func (h *AgendaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, itemID, err := parseEventSubResource(r.URL.Path, "agenda")
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	if err := h.agendaService.Remove(eventID, itemID, claims.UserID, claims.Role); err != nil {
		writeError(w, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List
// @Summary  Список пунктов повестки
// @Tags     agenda
// @Produce  json
// @Security BearerAuth
// @Param    id  path  int true "Event ID"
// @Success  200 {array} agendaItemResponse
// @Router   /events/{id}/agenda [get]
func (h *AgendaHandler) List(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseEventResource(r.URL.Path, "agenda")
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	items, err := h.agendaService.ListByEventID(eventID)
	if err != nil {
		h.logger.Error("failed to list agenda items", zap.Error(err))
		http.Error(w, "failed to list agenda items", http.StatusInternalServerError)
		return
	}

	resp := make([]agendaItemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, toAgendaItemResponse(item))
	}

	writeJSON(w, http.StatusOK, resp)
}

func toAgendaItemResponse(item model.AgendaItem) agendaItemResponse {
	return agendaItemResponse{
		ID:              item.ID,
		EventID:         item.EventID,
		Position:        item.Position,
		Title:           item.Title,
		Description:     item.Description,
		DurationMinutes: item.DurationMinutes,
		IsDone:          item.IsDone,
		CreatedAt:       item.CreatedAt.Format(time.RFC3339),
	}
}
