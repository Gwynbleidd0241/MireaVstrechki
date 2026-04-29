package service

import "testing"

func TestIsValidTaskStatus(t *testing.T) {
	valid := []string{"todo", "in_progress", "done"}
	for _, s := range valid {
		if !isValidTaskStatus(s) {
			t.Errorf("isValidTaskStatus(%q) = false, want true", s)
		}
	}

	invalid := []string{"", "TODO", "Done", "doing", "completed", "in-progress", "x", "DROP TABLE tasks;--"}
	for _, s := range invalid {
		if isValidTaskStatus(s) {
			t.Errorf("isValidTaskStatus(%q) = true, want false", s)
		}
	}
}

func TestIsValidParticipantRole(t *testing.T) {
	valid := []string{"participant", "responsible"}
	for _, r := range valid {
		if !isValidParticipantRole(r) {
			t.Errorf("isValidParticipantRole(%q) = false, want true", r)
		}
	}

	invalid := []string{"", "admin", "organizer", "PARTICIPANT", "guest", "owner"}
	for _, r := range invalid {
		if isValidParticipantRole(r) {
			t.Errorf("isValidParticipantRole(%q) = true, want false", r)
		}
	}
}

func TestCanEditTask(t *testing.T) {
	const (
		creatorID  int64 = 10
		assigneeID int64 = 20
		strangerID int64 = 99
	)

	assignee := assigneeID

	tests := []struct {
		name       string
		assignee   *int64
		userID     int64
		userRole   string
		want       bool
	}{
		{"admin always can", &assignee, strangerID, "admin", true},
		{"event creator can", &assignee, creatorID, "organizer", true},
		{"assignee can edit own task", &assignee, assigneeID, "employee", true},
		{"stranger cannot", &assignee, strangerID, "employee", false},
		{"unassigned task: stranger cannot", nil, strangerID, "employee", false},
		{"unassigned task: admin still can", nil, strangerID, "admin", true},
		{"unassigned task: creator still can", nil, creatorID, "organizer", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canEditTask(creatorID, tt.assignee, tt.userID, tt.userRole)
			if got != tt.want {
				t.Errorf("canEditTask = %v, want %v", got, tt.want)
			}
		})
	}
}
