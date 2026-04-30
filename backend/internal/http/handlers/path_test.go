package handlers

import "testing"

func TestParseEventID(t *testing.T) {
	tests := []struct {
		path    string
		want    int64
		wantErr bool
	}{
		{"/events/1", 1, false},
		{"/events/42", 42, false},
		{"/events/0", 0, true},
		{"/events/-1", 0, true},
		{"/events/abc", 0, true},
		{"/events", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		got, err := parseEventID(tt.path)

		if (err != nil) != tt.wantErr {
			t.Errorf("parseEventID(%q) err=%v, wantErr=%v", tt.path, err, tt.wantErr)
		}

		if got != tt.want {
			t.Errorf("parseEventID(%q) = %d, want %d", tt.path, got, tt.want)
		}
	}
}

func TestParseEventResource(t *testing.T) {
	tests := []struct {
		path     string
		resource string
		want     int64
		wantErr  bool
	}{
		{"/events/1/tasks", "tasks", 1, false},
		{"/events/5/agenda", "agenda", 5, false},
		{"/events/1/tasks", "agenda", 0, true},
		{"/events/0/tasks", "tasks", 0, true},
		{"/events/-1/tasks", "tasks", 0, true},
		{"/events/abc/tasks", "tasks", 0, true},
	}

	for _, tt := range tests {
		got, err := parseEventResource(tt.path, tt.resource)

		if (err != nil) != tt.wantErr {
			t.Errorf("parseEventResource(%q, %q) err=%v, wantErr=%v",
				tt.path, tt.resource, err, tt.wantErr)
		}

		if got != tt.want {
			t.Errorf("parseEventResource(%q, %q) = %d, want %d",
				tt.path, tt.resource, got, tt.want)
		}
	}
}
