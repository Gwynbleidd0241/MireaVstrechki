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

	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			http.Error(w, "invalid due_date", http.StatusBadRequest)
			return
		}
		dueDate = &parsed
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
	case service.ErrInvalidTaskStatus:
		http.Error(w, "invalid task status", http.StatusBadRequest)
	default:
		h.logger.Error("task handler error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func eventIDFromPath(path string) (int64, error) {
	// ожидаем путь вида /events/{id}/tasks
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 3 || parts[0] != "events" || parts[2] != "tasks" {
		return 0, strconv.ErrSyntax
	}

	return strconv.ParseInt(parts[1], 10, 64)
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
