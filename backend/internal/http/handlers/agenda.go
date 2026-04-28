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
	Title           string `json:"title"`
	Description     string `json:"description"`
	DurationMinutes *int   `json:"duration_minutes"`
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

func (h *AgendaHandler) Add(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, err := eventIDFromAgendaPath(r.URL.Path)
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
		h.handleAgendaError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toAgendaItemResponse(*item))
}

func (h *AgendaHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, itemID, err := eventAndAgendaIDFromPath(r.URL.Path)
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
		h.handleAgendaError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toAgendaItemResponse(*item))
}

func (h *AgendaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, itemID, err := eventAndAgendaIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	if err := h.agendaService.Remove(eventID, itemID, claims.UserID, claims.Role); err != nil {
		h.handleAgendaError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AgendaHandler) List(w http.ResponseWriter, r *http.Request) {
	eventID, err := eventIDFromAgendaPath(r.URL.Path)
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

func (h *AgendaHandler) handleAgendaError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrAgendaItemTitleRequired:
		http.Error(w, "agenda item title required", http.StatusBadRequest)
	case service.ErrAgendaItemTitleTooLong:
		http.Error(w, "agenda item title too long", http.StatusBadRequest)
	case service.ErrAgendaItemDescriptionTooLong:
		http.Error(w, "agenda item description too long", http.StatusBadRequest)
	case service.ErrInvalidAgendaDuration:
		http.Error(w, "invalid agenda item duration", http.StatusBadRequest)
	case service.ErrAgendaItemNotFound:
		http.Error(w, "agenda item not found", http.StatusNotFound)
	case service.ErrEventNotFound:
		http.Error(w, "event not found", http.StatusNotFound)
	case service.ErrPermissionDenied:
		http.Error(w, "permission denied", http.StatusForbidden)
	default:
		h.logger.Error("agenda handler error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func eventIDFromAgendaPath(path string) (int64, error) {
	// /events/{id}/agenda
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 3 || parts[0] != "events" || parts[2] != "agenda" {
		return 0, strconv.ErrSyntax
	}

	return strconv.ParseInt(parts[1], 10, 64)
}

func eventAndAgendaIDFromPath(path string) (int64, int64, error) {
	// /events/{id}/agenda/{aid}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 4 || parts[0] != "events" || parts[2] != "agenda" {
		return 0, 0, strconv.ErrSyntax
	}

	eventID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	itemID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return eventID, itemID, nil
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
