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
	Title       string  `json:"title"        example:"Подготовить макеты"`
	Description string  `json:"description"`
	Status      string  `json:"status"       example:"todo" enums:"todo,in_progress,done"`
	AssigneeID  *int64  `json:"assignee_id"`
	DueDate     *string `json:"due_date"     example:"2026-05-10T18:00:00Z"`
}

type updateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"      enums:"todo,in_progress,done"`
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

// Create godoc
// @Summary  Создать задачу во встрече
// @Tags     tasks
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id    path     int               true "Event ID"
// @Param    input body     createTaskRequest true "Параметры задачи"
// @Success  201   {object} taskResponse
// @Failure  400   {string} string "validation error"
// @Failure  403   {string} string "permission denied"
// @Router   /events/{id}/tasks [post]
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, err := parseEventResource(r.URL.Path, "tasks")
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
		EventID:        eventID,
		Title:          req.Title,
		Description:    req.Description,
		Status:         req.Status,
		AssigneeID:     req.AssigneeID,
		DueDate:        dueDate,
		ActingUserID:   claims.UserID,
		ActingUserRole: claims.Role,
	})
	if err != nil {
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusCreated, toTaskResponse(*task))
}

// Get godoc
// @Summary  Получить задачу по id
// @Tags     tasks
// @Produce  json
// @Security BearerAuth
// @Param    id     path     int true "Event ID"
// @Param    taskId path     int true "Task ID"
// @Success  200    {object} taskResponse
// @Failure  404    {string} string "task not found"
// @Router   /events/{id}/tasks/{taskId} [get]
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	eventID, taskID, err := parseEventSubResource(r.URL.Path, "tasks")
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.Get(eventID, taskID)
	if err != nil {
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(*task))
}

// Update godoc
// @Summary  Обновить задачу
// @Tags     tasks
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    id     path     int               true "Event ID"
// @Param    taskId path     int               true "Task ID"
// @Param    input  body     updateTaskRequest true "Новые поля"
// @Success  200    {object} taskResponse
// @Failure  400    {string} string "validation error"
// @Failure  403    {string} string "permission denied"
// @Failure  404    {string} string "task not found"
// @Router   /events/{id}/tasks/{taskId} [patch]
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, taskID, err := parseEventSubResource(r.URL.Path, "tasks")
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
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(*task))
}

// Delete godoc
// @Summary  Удалить задачу
// @Tags     tasks
// @Security BearerAuth
// @Param    id     path     int true "Event ID"
// @Param    taskId path     int true "Task ID"
// @Success  204    {string} string "no content"
// @Failure  403    {string} string "permission denied"
// @Failure  404    {string} string "task not found"
// @Router   /events/{id}/tasks/{taskId} [delete]
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	eventID, taskID, err := parseEventSubResource(r.URL.Path, "tasks")
	if err != nil {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	if err := h.taskService.Delete(eventID, taskID, claims.UserID, claims.Role); err != nil {
		writeError(w, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List godoc
// @Summary  Все задачи встречи
// @Tags     tasks
// @Produce  json
// @Security BearerAuth
// @Param    id  path  int true "Event ID"
// @Success  200 {array} taskResponse
// @Router   /events/{id}/tasks [get]
func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	eventID, err := parseEventResource(r.URL.Path, "tasks")
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
