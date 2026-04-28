package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

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
	})
	if err != nil {
		h.handleAgendaError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toAgendaItemResponse(*item))
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
	default:
		h.logger.Error("agenda handler error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func eventIDFromAgendaPath(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 3 || parts[0] != "events" || parts[2] != "agenda" {
		return 0, strconv.ErrSyntax
	}

	return strconv.ParseInt(parts[1], 10, 64)
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
