package model

import "time"

type Event struct {
	ID          int64
	Title       string
	Description string
	Status      string
	Location    string
	MeetingURL  string
	StartTime   time.Time
	EndTime     time.Time
	CreatorID   int64
	CreatedAt   time.Time
}
