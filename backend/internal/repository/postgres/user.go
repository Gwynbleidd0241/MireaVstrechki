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
