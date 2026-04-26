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
	userService *service.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
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
	case service.ErrEmailRequired:
		http.Error(w, "email required", http.StatusBadRequest)
	case service.ErrPasswordRequired:
		http.Error(w, "password required", http.StatusBadRequest)
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
