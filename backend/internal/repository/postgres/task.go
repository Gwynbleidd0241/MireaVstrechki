package postgres

import (
	"database/sql"

	"meeting-service/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task model.Task) (*model.Task, error) {
	var created model.Task

	err := r.db.QueryRow(
		`INSERT INTO tasks (event_id, title, description, status, assignee_id, due_date)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, event_id, title, description, status, assignee_id, due_date, created_at`,
		task.EventID,
		task.Title,
		task.Description,
		task.Status,
		task.AssigneeID,
		task.DueDate,
	).Scan(
		&created.ID,
		&created.EventID,
		&created.Title,
		&created.Description,
		&created.Status,
		&created.AssigneeID,
		&created.DueDate,
		&created.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *TaskRepository) ListByEventID(eventID int64) ([]model.Task, error) {
	rows, err := r.db.Query(
		`SELECT id, event_id, title, description, status, assignee_id, due_date, created_at
		 FROM tasks
		 WHERE event_id = $1
		 ORDER BY created_at ASC`,
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var task model.Task

		if err := rows.Scan(
			&task.ID,
			&task.EventID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.AssigneeID,
			&task.DueDate,
			&task.CreatedAt,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
