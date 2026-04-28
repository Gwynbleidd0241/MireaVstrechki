package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"meeting-service/internal/auth"
	"meeting-service/internal/http/middleware"
	"meeting-service/internal/service"
)

type UserHandler struct {
	userService  *service.UserService
	eventService *service.EventService
	taskService  *service.TaskService
	logger       *zap.Logger
}

func NewUserHandler(
	userService *service.UserService,
	eventService *service.EventService,
	taskService *service.TaskService,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		userService:  userService,
		eventService: eventService,
		taskService:  taskService,
		logger:       logger,
	}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type registerResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

type meResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(service.RegisterUserRequest{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		h.handleUserError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, registerResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	result, err := h.userService.Login(service.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.handleUserError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{
		ID:    result.User.ID,
		Email: result.User.Email,
		Role:  result.User.Role,
		Token: result.Token,
	})
}

func (h *UserHandler) handleUserError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrInvalidEmail:
		http.Error(w, "invalid email", http.StatusBadRequest)
	case service.ErrPasswordTooShort:
		http.Error(w, "password too short", http.StatusBadRequest)
	case service.ErrPasswordTooLong:
		http.Error(w, "password too long", http.StatusBadRequest)
	case service.ErrInvalidRole:
		http.Error(w, "invalid role", http.StatusBadRequest)
	case service.ErrInvalidCredentials:
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
	default:
		h.logger.Error("user handler error", zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, meResponse{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	})
}

func (h *UserHandler) MyEvents(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	events, err := h.eventService.ListForUser(claims.UserID)
	if err != nil {
		h.logger.Error("failed to list user events", zap.Error(err))
		http.Error(w, "failed to list events", http.StatusInternalServerError)
		return
	}

	resp := make([]eventResponse, 0, len(events))
	for _, event := range events {
		resp = append(resp, toEventResponse(event))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.List()
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		http.Error(w, "failed to list users", http.StatusInternalServerError)
		return
	}

	resp := make([]userResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, userResponse{
			ID:    u.ID,
			Email: u.Email,
			Role:  u.Role,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) MyTasks(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := h.taskService.ListForAssignee(claims.UserID)
	if err != nil {
		h.logger.Error("failed to list user tasks", zap.Error(err))
		http.Error(w, "failed to list tasks", http.StatusInternalServerError)
		return
	}

	resp := make([]taskResponse, 0, len(tasks))
	for _, task := range tasks {
		resp = append(resp, toTaskResponse(task))
	}

	writeJSON(w, http.StatusOK, resp)
}
