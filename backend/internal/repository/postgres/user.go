package postgres

import (
	"database/sql"

	"meeting-service/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, password, role string) (*model.User, error) {
	var id int64

	err := r.db.QueryRow(
		`INSERT INTO users (email, password_hash, role)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		email,
		password,
		role,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:           id,
		Email:        email,
		PasswordHash: password,
		Role:         role,
	}, nil
}

func (r *UserRepository) List() ([]model.User, error) {
	rows, err := r.db.Query(
		`SELECT id, email, password_hash, role
		 FROM users
		 ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User

		if err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(
		`SELECT id, email, password_hash, role
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
