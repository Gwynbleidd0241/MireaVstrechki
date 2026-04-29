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
	Email    string `json:"email"    example:"petr@example.com"`
	Password string `json:"password" example:"super-secret"`
	Role     string `json:"role"     example:"employee" enums:"admin,organizer,employee"`
}

type registerResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email"    example:"petr@example.com"`
	Password string `json:"password" example:"super-secret"`
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

// Register godoc
// @Summary  Регистрация нового пользователя
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    input body     registerRequest true "Учётные данные"
// @Success  201   {object} registerResponse
// @Failure  400   {string} string "validation error"
// @Router   /register [post]
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
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusCreated, registerResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}

// Login godoc
// @Summary  Вход в систему, возвращает JWT
// @Tags     auth
// @Accept   json
// @Produce  json
// @Param    input body     loginRequest true "Email и пароль"
// @Success  200   {object} loginResponse
// @Failure  401   {string} string "invalid email or password"
// @Router   /login [post]
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
		writeError(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{
		ID:    result.User.ID,
		Email: result.User.Email,
		Role:  result.User.Role,
		Token: result.Token,
	})
}

// Me godoc
// @Summary  Профиль текущего пользователя
// @Tags     auth
// @Produce  json
// @Security BearerAuth
// @Success  200 {object} meResponse
// @Failure  401 {string} string "unauthorized"
// @Router   /me [get]
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

// MyEvents godoc
// @Summary  Встречи, в которых я создатель или участник
// @Tags     auth
// @Produce  json
// @Security BearerAuth
// @Success  200 {array}  eventResponse
// @Failure  401 {string} string "unauthorized"
// @Router   /me/events [get]
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

// MyTasks godoc
// @Summary  Задачи, назначенные на меня
// @Tags     auth
// @Produce  json
// @Security BearerAuth
// @Success  200 {array}  taskResponse
// @Failure  401 {string} string "unauthorized"
// @Router   /me/tasks [get]
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

// List godoc
// @Summary  Список всех пользователей
// @Tags     users
// @Produce  json
// @Security BearerAuth
// @Success  200 {array}  userResponse
// @Failure  401 {string} string "unauthorized"
// @Router   /users [get]
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
