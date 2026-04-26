package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"meeting-service/internal/repository/postgres"
)

type RegisterHandler struct {
	userRepo *postgres.UserRepository
	logger   *zap.Logger
}

func NewRegisterHandler(userRepo *postgres.UserRepository, logger *zap.Logger) *RegisterHandler {
	return &RegisterHandler{
		userRepo: userRepo,
		logger:   logger,
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

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	if req.Role == "" {
		req.Role = "employee"
	}

	if req.Role != "admin" && req.Role != "organizer" && req.Role != "employee" {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.Create(req.Email, req.Password, req.Role)
	if err != nil {
		h.logger.Error("failed to create user", zap.Error(err))
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	resp := registerResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
	}
}
