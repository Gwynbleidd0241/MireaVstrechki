//go:build integration

package postgres_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"meeting-service/internal/model"
	"meeting-service/internal/repository/postgres"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("INTEGRATION_DSN")
	if dsn == "" {
		t.Skip("INTEGRATION_DSN is not set; skipping integration test")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v (DSN may point to an unreachable host)", err)
	}

	t.Cleanup(func() { _ = db.Close() })
	return db
}

func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s_%d@integration.test", prefix, time.Now().UnixNano())
}

func TestUserRepository_CreateAndGetByEmail(t *testing.T) {
	db := setupDB(t)
	repo := postgres.NewUserRepository(db)

	email := uniqueEmail("user")
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM users WHERE email = $1`, email) })

	created, err := repo.Create(email, "hash-placeholder", "employee")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if created.ID == 0 {
		t.Error("expected non-zero ID after insert")
	}

	if created.Email != email {
		t.Errorf("Email = %q, want %q", created.Email, email)
	}

	fetched, err := repo.GetByEmail(email)
	if err != nil {
		t.Fatalf("GetByEmail: %v", err)
	}

	if fetched.ID != created.ID {
		t.Errorf("fetched.ID = %d, want %d", fetched.ID, created.ID)
	}

	if fetched.Role != "employee" {
		t.Errorf("Role = %q, want employee", fetched.Role)
	}
}

func TestUserRepository_DuplicateEmailRejected(t *testing.T) {
	db := setupDB(t)
	repo := postgres.NewUserRepository(db)

	email := uniqueEmail("dup")
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM users WHERE email = $1`, email) })

	if _, err := repo.Create(email, "hash", "employee"); err != nil {
		t.Fatalf("first Create: %v", err)
	}

	if _, err := repo.Create(email, "hash", "employee"); err == nil {
		t.Error("expected error on duplicate email, got nil")
	}
}

func TestEventRepository_FullCRUDLifecycle(t *testing.T) {
	db := setupDB(t)

	users := postgres.NewUserRepository(db)
	events := postgres.NewEventRepository(db)

	user, err := users.Create(uniqueEmail("event-creator"), "hash", "organizer")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM users WHERE id = $1`, user.ID) })

	start := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	end := start.Add(time.Hour)

	created, err := events.Create(model.Event{
		Title:       "integration meeting",
		Description: "test desc",
		StartTime:   start,
		EndTime:     end,
		CreatorID:   user.ID,
	})
	if err != nil {
		t.Fatalf("Create event: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM events WHERE id = $1`, created.ID) })

	if created.ID == 0 {
		t.Fatal("expected non-zero event ID")
	}

	fetched, err := events.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if fetched == nil {
		t.Fatal("event not found right after Create")
	}
	if fetched.Title != "integration meeting" {
		t.Errorf("Title = %q, want %q", fetched.Title, "integration meeting")
	}

	fetched.Title = "updated title"
	updated, err := events.Update(*fetched)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Title != "updated title" {
		t.Errorf("after Update Title = %q, want %q", updated.Title, "updated title")
	}

	if err := events.Delete(created.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	after, err := events.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID after Delete: %v", err)
	}
	if after != nil {
		t.Error("expected nil after Delete, got event")
	}
}

func TestEventRepository_ListForUserIncludesCreatorAndParticipant(t *testing.T) {
	db := setupDB(t)

	users := postgres.NewUserRepository(db)
	events := postgres.NewEventRepository(db)
	participants := postgres.NewParticipantRepository(db)

	owner, err := users.Create(uniqueEmail("owner"), "hash", "organizer")
	if err != nil {
		t.Fatalf("create owner: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM users WHERE id = $1`, owner.ID) })

	guest, err := users.Create(uniqueEmail("guest"), "hash", "employee")
	if err != nil {
		t.Fatalf("create guest: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM users WHERE id = $1`, guest.ID) })

	start := time.Now().Add(2 * time.Hour).UTC().Truncate(time.Second)
	event, err := events.Create(model.Event{
		Title:     "list test",
		StartTime: start,
		EndTime:   start.Add(time.Hour),
		CreatorID: owner.ID,
	})
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM events WHERE id = $1`, event.ID) })

	if _, err := participants.Add(model.Participant{
		EventID: event.ID,
		UserID:  guest.ID,
		Role:    "participant",
	}); err != nil {
		t.Fatalf("add participant: %v", err)
	}

	ownerEvents, err := events.ListForUser(owner.ID)
	if err != nil {
		t.Fatalf("ListForUser(owner): %v", err)
	}
	if !containsEvent(ownerEvents, event.ID) {
		t.Error("event not visible to its creator")
	}

	guestEvents, err := events.ListForUser(guest.ID)
	if err != nil {
		t.Fatalf("ListForUser(guest): %v", err)
	}
	if !containsEvent(guestEvents, event.ID) {
		t.Error("event not visible to its participant")
	}
}

func TestTaskRepository_StatusUpdateRoundtrip(t *testing.T) {
	db := setupDB(t)

	users := postgres.NewUserRepository(db)
	events := postgres.NewEventRepository(db)
	tasks := postgres.NewTaskRepository(db)

	user, err := users.Create(uniqueEmail("task-owner"), "hash", "organizer")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM users WHERE id = $1`, user.ID) })

	now := time.Now().UTC().Truncate(time.Second)
	event, err := events.Create(model.Event{
		Title:     "task test",
		StartTime: now.Add(time.Hour),
		EndTime:   now.Add(2 * time.Hour),
		CreatorID: user.ID,
	})
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Exec(`DELETE FROM events WHERE id = $1`, event.ID) })

	task, err := tasks.Create(model.Task{
		EventID: event.ID,
		Title:   "do the thing",
		Status:  "todo",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	task.Status = "done"
	updated, err := tasks.Update(*task)
	if err != nil {
		t.Fatalf("update task: %v", err)
	}

	if updated.Status != "done" {
		t.Errorf("Status = %q, want done", updated.Status)
	}
}

func containsEvent(list []model.Event, id int64) bool {
	for _, e := range list {
		if e.ID == id {
			return true
		}
	}
	return false
}
