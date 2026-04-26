package model

import "time"

type Participant struct {
	ID        int64
	EventID   int64
	UserID    int64
	Role      string
	CreatedAt time.Time
}
