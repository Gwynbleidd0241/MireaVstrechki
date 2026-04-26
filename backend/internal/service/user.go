package service

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"meeting-service/internal/auth"
	"meeting-service/internal/model"
	"meeting-service/internal/repository/postgres"
)

var (
	ErrInvalidRole        = errors.New("invalid role")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailRequired      = errors.New("email required")
	ErrPasswordRequired   = errors.New("password required")
)

type UserService struct {
	userRepo  *postgres.UserRepository
	jwtSecret string
}

func NewUserService(userRepo *postgres.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

type RegisterUserRequest struct {
	Email    string
	Password string
	Role     string
}

type LoginUserRequest struct {
	Email    string
	Password string
}

type LoginResult struct {
	User  *model.User
	Token string
}

func (s *UserService) Register(req RegisterUserRequest) (*model.User, error) {
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	if email == "" {
		return nil, ErrEmailRequired
	}

	if password == "" {
		return nil, ErrPasswordRequired
	}

	role := req.Role
	if role == "" {
		role = "employee"
	}

	if role != "admin" && role != "organizer" && role != "employee" {
		return nil, ErrInvalidRole
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	return s.userRepo.Create(email, string(passwordHash), role)
}

func (s *UserService) Login(req LoginUserRequest) (*LoginResult, error) {
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:  user,
		Token: token,
	}, nil
}
