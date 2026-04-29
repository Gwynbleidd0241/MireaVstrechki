package handlers

import "testing"

func FuzzEventIDFromPath(f *testing.F) {
	f.Add("/events/1/tasks")
	f.Add("/events/abc/tasks")
	f.Add("/events//tasks")
	f.Add("")
	f.Add("/")
	f.Add("///")
	f.Add("/events")
	f.Add("/events/-1/tasks")
	f.Add("/events/9999999999999999999999999/tasks")
	f.Add("/admin/events/1/tasks")
	f.Add("/events/1/tasks/extra")

	f.Fuzz(func(t *testing.T, path string) {
		id, err := parseEventResource(path, "tasks")

		if err != nil {
			return
		}

		if id <= 0 {
			t.Errorf("parseEventResource(%q, tasks) returned id=%d without error", path, id)
		}
	})
}

func FuzzEventAndTaskIDFromPath(f *testing.F) {
	f.Add("/events/1/tasks/2")
	f.Add("/events/abc/tasks/5")
	f.Add("/events/1/tasks/abc")
	f.Add("")
	f.Add("/events/1/tasks/")
	f.Add("/events/1/tasks/2/extra")
	f.Add("/events//tasks//")

	f.Fuzz(func(t *testing.T, path string) {
		eventID, taskID, err := parseEventSubResource(path, "tasks")

		if err != nil {
			return
		}

		if eventID <= 0 || taskID <= 0 {
			t.Errorf("parseEventSubResource(%q, tasks) = (%d, %d) without error",
				path, eventID, taskID)
		}
	})
}
