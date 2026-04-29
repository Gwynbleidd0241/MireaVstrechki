package seed

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

const demoAdminEmail = "admin@demo.local"

func Run(db *sql.DB) error {
	var exists bool
	if err := db.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`,
		demoAdminEmail,
	).Scan(&exists); err != nil {
		return err
	}

	if exists {
		log.Println("seed: already applied, skipping")
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	pwd := string(hash)

	var adminID, organizerID, anneID, ivanID int64

	if err := db.QueryRow(
		`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, 'admin') RETURNING id`,
		demoAdminEmail, pwd,
	).Scan(&adminID); err != nil {
		return err
	}

	if err := db.QueryRow(
		`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, 'organizer') RETURNING id`,
		"organizer@demo.local", pwd,
	).Scan(&organizerID); err != nil {
		return err
	}

	if err := db.QueryRow(
		`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, 'employee') RETURNING id`,
		"anna@demo.local", pwd,
	).Scan(&anneID); err != nil {
		return err
	}

	if err := db.QueryRow(
		`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, 'employee') RETURNING id`,
		"ivan@demo.local", pwd,
	).Scan(&ivanID); err != nil {
		return err
	}

	var planningID, retroID int64

	if err := db.QueryRow(
		`INSERT INTO events (title, description, start_time, end_time, creator_id)
		 VALUES ($1, $2, NOW() + interval '1 day', NOW() + interval '1 day' + interval '1 hour', $3)
		 RETURNING id`,
		"Планирование спринта",
		"Обсуждаем задачи на ближайший спринт, расставляем приоритеты и распределяем ответственных.",
		organizerID,
	).Scan(&planningID); err != nil {
		return err
	}

	if err := db.QueryRow(
		`INSERT INTO events (title, description, start_time, end_time, creator_id)
		 VALUES ($1, $2, NOW() + interval '3 day', NOW() + interval '3 day' + interval '90 minutes', $3)
		 RETURNING id`,
		"Ретроспектива",
		"Что получилось хорошо, что можно улучшить, какие действия зафиксируем на следующий спринт.",
		organizerID,
	).Scan(&retroID); err != nil {
		return err
	}

	if _, err := db.Exec(
		`INSERT INTO event_participants (event_id, user_id, role) VALUES
		 ($1, $2, 'responsible'),
		 ($1, $3, 'participant'),
		 ($4, $2, 'participant'),
		 ($4, $3, 'responsible'),
		 ($4, $5, 'participant')`,
		planningID, anneID, ivanID,
		retroID, adminID,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		`INSERT INTO tasks (event_id, title, description, status, assignee_id, due_date) VALUES
		 ($1, 'Подготовить дизайн макетов', 'Высокоуровневые макеты главной страницы и карточки встречи', 'in_progress', $2, NOW() + interval '5 day'),
		 ($1, 'Согласовать API контракт', 'OpenAPI-схема с фронтендом', 'todo', $3, NULL),
		 ($1, 'Настроить CI', 'Базовый pipeline: lint, build, test, fuzz smoke', 'done', $3, NULL)`,
		planningID, anneID, ivanID,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		`INSERT INTO agenda_items (event_id, position, title, description, duration_minutes, is_done) VALUES
		 ($1, 1, 'Цели спринта', 'Зафиксируем три главных результата', 10, false),
		 ($1, 2, 'Распределение задач', 'Кто что берёт на себя', 20, false),
		 ($1, 3, 'Риски и блокеры', 'Что может помешать спринту', 15, false)`,
		planningID,
	); err != nil {
		return err
	}

	if _, err := db.Exec(
		`INSERT INTO agenda_items (event_id, position, title, description, duration_minutes, is_done) VALUES
		 ($1, 1, 'Что получилось хорошо', '', 20, false),
		 ($1, 2, 'Что можно улучшить', '', 25, false),
		 ($1, 3, 'Action items', 'Конкретные действия на следующий спринт', 20, false)`,
		retroID,
	); err != nil {
		return err
	}

	log.Println("seed: applied successfully (4 users, 2 events, 3 tasks, 6 agenda items)")
	return nil
}
