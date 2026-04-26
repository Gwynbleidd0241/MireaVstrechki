package service

import (
	"errors"
	"strings"
	"time"

	"meeting-service/internal/model"
	"meeting-service/internal/repository/postgres"
)

var (
	ErrTaskTitleRequired = errors.New("task title required")
	ErrInvalidTaskStatus = errors.New("invalid task status")
)

type TaskService struct {
	taskRepo *postgres.TaskRepository
}

func NewTaskService(taskRepo *postgres.TaskRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

type CreateTaskRequest struct {
	EventID     int64
	Title       string
	Description string
	Status      string
	AssigneeID  *int64
	DueDate     *time.Time
}

func (s *TaskService) Create(req CreateTaskRequest) (*model.Task, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrTaskTitleRequired
	}

	status := req.Status
	if status == "" {
		status = "todo"
	}

	if status != "todo" && status != "in_progress" && status != "done" {
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

func (s *TaskService) ListByEventID(eventID int64) ([]model.Task, error) {
	return s.taskRepo.ListByEventID(eventID)
}
