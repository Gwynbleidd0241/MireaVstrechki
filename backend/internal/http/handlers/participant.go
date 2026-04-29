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

type ParticipantHandler struct {
	participantService *service.ParticipantService
	logger             *zap.Logger
}

func NewParticipantHandler(participantService *service.ParticipantService, logger *zap.Logger) *ParticipantHandler {
	return &ParticipantHandler{
		participantService: participantService,
		logger:             logger,
	}
}

type addParticipantRequest struct {
	UserID int64  `json:"user_id" example:"5"`
	Role   string `json:"role"    example:"participant" enums:"participant,responsible"`
}

type updateParticipantRequest struct {
	Role string `json:"role" enums:"participant,responsible"`
}

type participantResponse struct {
	ID        int64  `json:"id"`
	EventID   int64  `json:"event_id"`
	UserID    int64  `json:"user_id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// Add godoc
// @Summary  Добавить участника во встречу
// @Tags     participants
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id    path     int                   true "Event ID"
// @Param    input body     addParticipantRequest true "Пользователь и роль"
// @Success  201   {object} participantResponse
// @Failure  400   {string} string "validation error"
// @Failure  403   {string} string "permission denied"
// @Router   /events/{id}/participants [post]
func (h *ParticipantHandler) Add(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, err := parseEventResource(r.URL.Path, "participants")
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	var req addParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	participant, err := h.participantService.Add(service.AddParticipantRequest{
		EventID:        eventID,
		UserID:         req.UserID,
		Role:           req.Role,
		ActingUserID:   claims.UserID,
		ActingUserRole: claims.Role,
	})
	if err != nil {
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusCreated, toParticipantResponse(*participant))
}

// Update godoc
// @Summary  Изменить роль участника
// @Tags     participants
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id            path     int                      true "Event ID"
// @Param    participantId path     int                      true "Participant ID"
// @Param    input         body     updateParticipantRequest true "Новая роль"
// @Success  200           {object} participantResponse
// @Failure  403           {string} string "permission denied"
// @Failure  404           {string} string "participant not found"
// @Router   /events/{id}/participants/{participantId} [patch]
func (h *ParticipantHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, participantID, err := parseEventSubResource(r.URL.Path, "participants")
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	var req updateParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	participant, err := h.participantService.Update(service.UpdateParticipantRequest{
		EventID:        eventID,
		ParticipantID:  participantID,
		Role:           req.Role,
		ActingUserID:   claims.UserID,
		ActingUserRole: claims.Role,
	})
	if err != nil {
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusOK, toParticipantResponse(*participant))
}

// Remove godoc
// @Summary  Удалить участника из встречи
// @Tags     participants
// @Security BearerAuth
// @Param    id            path     int true "Event ID"
// @Param    participantId path     int true "Participant ID"
// @Success  204           {string} string "no content"
// @Failure  403           {string} string "permission denied"
// @Failure  404           {string} string "participant not found"
// @Router   /events/{id}/participants/{participantId} [delete]
func (h *ParticipantHandler) Remove(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, participantID, err := parseEventSubResource(r.URL.Path, "participants")
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	if err := h.participantService.Remove(eventID, participantID, claims.UserID, claims.Role); err != nil {
		writeError(w, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List godoc
// @Summary  Список участников встречи
// @Tags     participants
// @Produce  json
// @Security BearerAuth
// @Param    id  path  int true "Event ID"
// @Success  200 {array} participantResponse
// @Router   /events/{id}/participants [get]
func (h *ParticipantHandler) List(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseEventResource(r.URL.Path, "participants")
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	participants, err := h.participantService.ListByEventID(eventID)
	if err != nil {
		h.logger.Error("failed to list participants", zap.Error(err))
		http.Error(w, "failed to list participants", http.StatusInternalServerError)
		return
	}

	resp := make([]participantResponse, 0, len(participants))
	for _, participant := range participants {
		resp = append(resp, toParticipantResponse(participant))
	}

	writeJSON(w, http.StatusOK, resp)
}

func toParticipantResponse(participant model.Participant) participantResponse {
	return participantResponse{
		ID:        participant.ID,
		EventID:   participant.EventID,
		UserID:    participant.UserID,
		Role:      participant.Role,
		CreatedAt: participant.CreatedAt.Format(time.RFC3339),
	}
}
