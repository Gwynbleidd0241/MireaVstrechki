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

type TaskHandler struct {
	taskService *service.TaskService
	logger      *zap.Logger
}

func NewTaskHandler(taskService *service.TaskService, logger *zap.Logger) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		logger:      logger,
	}
}

type createTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	AssigneeID  *int64  `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}

type updateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	AssigneeID  *int64  `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}

type taskResponse struct {
	ID          int64   `json:"id"`
	EventID     int64   `json:"event_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	AssigneeID  *int64  `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
	CreatedAt   string  `json:"created_at"`
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	eventID, err := eventIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	dueDate, err := parseOptionalDueDate(req.DueDate)
	if err != nil {
		http.Error(w, "invalid due_date", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.Create(service.CreateTaskRequest{
		EventID:     eventID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
		DueDate:     dueDate,
	})
	if err != nil {
		h.handleTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toTaskResponse(*task))
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	eventID, taskID, err := eventAndTaskIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.Get(eventID, taskID)
	if err != nil {
		h.handleTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(*task))
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, taskID, err := eventAndTaskIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	dueDate, err := parseOptionalDueDate(req.DueDate)
	if err != nil {
		http.Error(w, "invalid due_date", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.Update(service.UpdateTaskRequest{
		EventID:     eventID,
		TaskID:      taskID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
		DueDate:     dueDate,
		UserID:      claims.UserID,
		UserRole:    claims.Role,
	})
	if err != nil {
		h.handleTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(*task))
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, taskID, err := eventAndTaskIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	if err := h.taskService.Delete(eventID, taskID, claims.UserID, claims.Role); err != nil {
		h.handleTaskError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	eventID, err := eventIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid event id", http.StatusBadRequest)
		return
	}

	tasks, err := h.taskService.ListByEventID(eventID)
	if err != nil {
		h.logger.Error("failed to list tasks", zap.Error(err))
		http.Error(w, "failed to list tasks", http.StatusInternalServerError)
		return
	}

	resp := make([]taskResponse, 0, len(tasks))
	for _, task := range tasks {
		resp = append(resp, toTaskResponse(task))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) handleTaskError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrTaskTitleRequired:
		http.Error(w, "task title required", http.StatusBadRequest)
	case service.ErrTaskTitleTooLong:
		http.Error(w, "task title too long", http.StatusBadRequest)
	case service.ErrTaskDescriptionTooLong:
		http.Error(w, "task description too long", http.StatusBadRequest)
	case service.ErrInvalidTaskStatus:
		http.Error(w, "invalid task status", http.StatusBadRequest)
	case service.ErrTaskNotFound:
		http.Error(w, "task not found", http.StatusNotFound)
	case service.ErrEventNotFound:
		http.Error(w, "event not found", http.StatusNotFound)
	case service.ErrPermissionDenied:
		http.Error(w, "permission denied", http.StatusForbidden)
	default:
		h.logger.Error("task handler error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func eventIDFromPath(path string) (int64, error) {
	// /events/{id}/tasks
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 3 || parts[0] != "events" || parts[2] != "tasks" {
		return 0, strconv.ErrSyntax
	}

	return strconv.ParseInt(parts[1], 10, 64)
}

func eventAndTaskIDFromPath(path string) (int64, int64, error) {
	// /events/{id}/tasks/{taskID}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 4 || parts[0] != "events" || parts[2] != "tasks" {
		return 0, 0, strconv.ErrSyntax
	}

	eventID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	taskID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return eventID, taskID, nil
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

func toTaskResponse(task model.Task) taskResponse {
	var dueDate *string
	if task.DueDate != nil {
		formatted := task.DueDate.Format(time.RFC3339)
		dueDate = &formatted
	}

	return taskResponse{
		ID:          task.ID,
		EventID:     task.EventID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		AssigneeID:  task.AssigneeID,
		DueDate:     dueDate,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
	}
}
