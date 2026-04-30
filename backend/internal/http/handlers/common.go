package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"meeting-service/internal/service"
)

var statusByError = map[error]int{
	// 400 Bad Request
	service.ErrInvalidEmail:                 http.StatusBadRequest,
	service.ErrPasswordTooShort:             http.StatusBadRequest,
	service.ErrPasswordTooLong:              http.StatusBadRequest,
	service.ErrInvalidRole:                  http.StatusBadRequest,
	service.ErrEventTitleRequired:           http.StatusBadRequest,
	service.ErrEventTitleTooLong:            http.StatusBadRequest,
	service.ErrEventDescriptionTooLong:      http.StatusBadRequest,
	service.ErrInvalidEventTime:             http.StatusBadRequest,
	service.ErrTaskTitleRequired:            http.StatusBadRequest,
	service.ErrTaskTitleTooLong:             http.StatusBadRequest,
	service.ErrTaskDescriptionTooLong:       http.StatusBadRequest,
	service.ErrInvalidTaskStatus:            http.StatusBadRequest,
	service.ErrParticipantUserRequired:      http.StatusBadRequest,
	service.ErrInvalidParticipantRole:       http.StatusBadRequest,
	service.ErrAgendaItemTitleRequired:      http.StatusBadRequest,
	service.ErrAgendaItemTitleTooLong:       http.StatusBadRequest,
	service.ErrAgendaItemDescriptionTooLong: http.StatusBadRequest,
	service.ErrInvalidAgendaDuration:        http.StatusBadRequest,

	// 401 Unauthorized
	service.ErrInvalidCredentials: http.StatusUnauthorized,

	// 403 Forbidden
	service.ErrPermissionDenied: http.StatusForbidden,

	// 404 Not Found
	service.ErrEventNotFound:       http.StatusNotFound,
	service.ErrTaskNotFound:        http.StatusNotFound,
	service.ErrParticipantNotFound: http.StatusNotFound,
	service.ErrAgendaItemNotFound:  http.StatusNotFound,
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error, logger *zap.Logger) {
	if status, ok := statusByError[err]; ok {
		http.Error(w, err.Error(), status)
		return
	}

	if logger != nil {
		logger.Error("handler error", zap.Error(err))
	}

	http.Error(w, "internal error", http.StatusInternalServerError)
}

func parseEventID(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 || parts[0] != "events" {
		return 0, strconv.ErrSyntax
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, err
	}

	if id <= 0 {
		return 0, strconv.ErrSyntax
	}

	return id, nil
}

func parseEventResource(path, resource string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 3 || parts[0] != "events" || parts[2] != resource {
		return 0, strconv.ErrSyntax
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, err
	}

	if id <= 0 {
		return 0, strconv.ErrSyntax
	}

	return id, nil
}

func parseEventSubResource(path, resource string) (int64, int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 4 || parts[0] != "events" || parts[2] != resource {
		return 0, 0, strconv.ErrSyntax
	}

	eventID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	subID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	if eventID <= 0 || subID <= 0 {
		return 0, 0, strconv.ErrSyntax
	}

	return eventID, subID, nil
}

func parseOptionalDueDate(raw *string) (*time.Time, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, *raw)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
