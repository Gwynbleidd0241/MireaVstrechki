package service

import (
	"errors"
	"strings"
	"time"

	"meeting-service/internal/model"
	"meeting-service/internal/repository/postgres"
	"meeting-service/internal/validation"
)

var (
	ErrTaskTitleRequired      = errors.New("task title required")
	ErrTaskTitleTooLong       = errors.New("task title too long")
	ErrTaskDescriptionTooLong = errors.New("task description too long")
	ErrInvalidTaskStatus      = errors.New("invalid task status")
	ErrTaskNotFound           = errors.New("task not found")
)

type TaskService struct {
	taskRepo  *postgres.TaskRepository
	eventRepo *postgres.EventRepository
}

func NewTaskService(taskRepo *postgres.TaskRepository, eventRepo *postgres.EventRepository) *TaskService {
	return &TaskService{
		taskRepo:  taskRepo,
		eventRepo: eventRepo,
	}
}

type CreateTaskRequest struct {
	EventID        int64
	Title          string
	Description    string
	Status         string
	AssigneeID     *int64
	DueDate        *time.Time
	ActingUserID   int64
	ActingUserRole string
}

type UpdateTaskRequest struct {
	EventID     int64
	TaskID      int64
	Title       string
	Description string
	Status      string
	AssigneeID  *int64
	DueDate     *time.Time
	UserID      int64
	UserRole    string
}

func (s *TaskService) Create(req CreateTaskRequest) (*model.Task, error) {
	if err := s.checkEventManagePermission(req.EventID, req.ActingUserID, req.ActingUserRole); err != nil {
		return nil, err
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrTaskTitleRequired
	}

	if len(title) > validation.MaxTitleLength {
		return nil, ErrTaskTitleTooLong
	}

	if len(req.Description) > validation.MaxDescriptionLength {
		return nil, ErrTaskDescriptionTooLong
	}

	status := req.Status
	if status == "" {
		status = "todo"
	}

	if !isValidTaskStatus(status) {
		return nil, ErrInvalidTaskStatus
	}

	return s.taskRepo.Create(model.Task{
		EventID:     req.EventID,
		Title:       title,
		Description: strings.TrimSpace(req.Description),
		Status:      status,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
	})
}

func (s *TaskService) Get(eventID, taskID int64) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	if task == nil || task.EventID != eventID {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (s *TaskService) Update(req UpdateTaskRequest) (*model.Task, error) {
	task, err := s.taskRepo.GetByID(req.TaskID)
	if err != nil {
		return nil, err
	}

	if task == nil || task.EventID != req.EventID {
		return nil, ErrTaskNotFound
	}

	event, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, ErrEventNotFound
	}

	if !canEditTask(event.CreatorID, task.AssigneeID, req.UserID, req.UserRole) {
		return nil, ErrPermissionDenied
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrTaskTitleRequired
	}

	if len(title) > validation.MaxTitleLength {
		return nil, ErrTaskTitleTooLong
	}

	if len(req.Description) > validation.MaxDescriptionLength {
		return nil, ErrTaskDescriptionTooLong
	}

	if !isValidTaskStatus(req.Status) {
		return nil, ErrInvalidTaskStatus
	}

	return s.taskRepo.Update(model.Task{
		ID:          req.TaskID,
		EventID:     req.EventID,
		Title:       title,
		Description: strings.TrimSpace(req.Description),
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
	})
}

func (s *TaskService) Delete(eventID, taskID, userID int64, userRole string) error {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return err
	}

	if task == nil || task.EventID != eventID {
		return ErrTaskNotFound
	}

	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	if event == nil {
		return ErrEventNotFound
	}

	if userRole != "admin" && event.CreatorID != userID {
		return ErrPermissionDenied
	}

	return s.taskRepo.Delete(taskID)
}

func (s *TaskService) ListByEventID(eventID int64) ([]model.Task, error) {
	return s.taskRepo.ListByEventID(eventID)
}

func (s *TaskService) ListForAssignee(userID int64) ([]model.Task, error) {
	return s.taskRepo.ListForAssignee(userID)
}

func isValidTaskStatus(status string) bool {
	return status == "todo" || status == "in_progress" || status == "done"
}

func (s *TaskService) checkEventManagePermission(eventID, actingUserID int64, actingUserRole string) error {
	if actingUserRole == "admin" {
		return nil
	}

	event, err := s.eventRepo.GetByID(eventID)
	if err != nil {
		return err
	}

	if event == nil {
		return ErrEventNotFound
	}

	if event.CreatorID != actingUserID {
		return ErrPermissionDenied
	}

	return nil
}

func canEditTask(creatorID int64, assigneeID *int64, userID int64, userRole string) bool {
	if userRole == "admin" {
		return true
	}

	if creatorID == userID {
		return true
	}

	if assigneeID != nil && *assigneeID == userID {
		return true
	}

	return false
}
