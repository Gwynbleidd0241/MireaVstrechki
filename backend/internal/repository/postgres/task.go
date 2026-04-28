package postgres

import (
	"database/sql"
	"errors"

	"meeting-service/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetByID(id int64) (*model.Task, error) {
	var task model.Task

	err := r.db.QueryRow(
		`SELECT id, event_id, title, description, status, assignee_id, due_date, created_at
		 FROM tasks
		 WHERE id = $1`,
		id,
	).Scan(
		&task.ID,
		&task.EventID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.AssigneeID,
		&task.DueDate,
		&task.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) Update(task model.Task) (*model.Task, error) {
	var updated model.Task

	err := r.db.QueryRow(
		`UPDATE tasks
		 SET title = $1, description = $2, status = $3, assignee_id = $4, due_date = $5
		 WHERE id = $6
		 RETURNING id, event_id, title, description, status, assignee_id, due_date, created_at`,
		task.Title,
		task.Description,
		task.Status,
		task.AssigneeID,
		task.DueDate,
		task.ID,
	).Scan(
		&updated.ID,
		&updated.EventID,
		&updated.Title,
		&updated.Description,
		&updated.Status,
		&updated.AssigneeID,
		&updated.DueDate,
		&updated.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *TaskRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	return err
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

func (r *TaskRepository) ListForAssignee(userID int64) ([]model.Task, error) {
	rows, err := r.db.Query(
		`SELECT id, event_id, title, description, status, assignee_id, due_date, created_at
		 FROM tasks
		 WHERE assignee_id = $1
		 ORDER BY due_date ASC NULLS LAST, created_at ASC`,
		userID,
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
