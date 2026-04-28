package model

import "time"

type AgendaItem struct {
	ID              int64
	EventID         int64
	Position        int
	Title           string
	Description     string
	DurationMinutes *int
	IsDone          bool
	CreatedAt       time.Time
}
