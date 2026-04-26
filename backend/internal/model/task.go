package model

import "time"

type Task struct {
	ID          int64
	EventID     int64
	Title       string
	Description string
	Status      string
	AssigneeID  *int64
	DueDate     *time.Time
	CreatedAt   time.Time
}
