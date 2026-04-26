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
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

type participantResponse struct {
	ID        int64  `json:"id"`
	EventID   int64  `json:"event_id"`
	UserID    int64  `json:"user_id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

func (h *ParticipantHandler) Add(w http.ResponseWriter, r *http.Request) {
	eventID, err := eventIDFromParticipantsPath(r.URL.Path)
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
		EventID: eventID,
		UserID:  req.UserID,
		Role:    req.Role,
	})
	if err != nil {
		h.handleParticipantError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toParticipantResponse(*participant))
}

func (h *ParticipantHandler) List(w http.ResponseWriter, r *http.Request) {
	eventID, err := eventIDFromParticipantsPath(r.URL.Path)
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

func (h *ParticipantHandler) handleParticipantError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrParticipantUserRequired:
		http.Error(w, "participant user required", http.StatusBadRequest)
	case service.ErrInvalidParticipantRole:
		http.Error(w, "invalid participant role", http.StatusBadRequest)
	default:
		h.logger.Error("participant handler error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func eventIDFromParticipantsPath(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 3 || parts[0] != "events" || parts[2] != "participants" {
		return 0, strconv.ErrSyntax
	}

	return strconv.ParseInt(parts[1], 10, 64)
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
