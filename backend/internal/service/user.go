package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"meeting-service/internal/model"
	"meeting-service/internal/repository/postgres"
)

var (
	ErrInvalidRole = errors.New("invalid role")
)

type UserService struct {
	userRepo *postgres.UserRepository
}

func NewUserService(userRepo *postgres.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type RegisterUserRequest struct {
	Email    string
	Password string
	Role     string
}

func (s *UserService) Register(req RegisterUserRequest) (*model.User, error) {
	role := req.Role
	if role == "" {
		role = "employee"
	}

	if role != "admin" && role != "organizer" && role != "employee" {
		return nil, ErrInvalidRole
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	return s.userRepo.Create(req.Email, string(passwordHash), role)
}
