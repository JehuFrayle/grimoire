package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jehufrayle/grimoire/internal/users"
	"github.com/jehufrayle/grimoire/utils"
)

type Handler struct {
	userRepo users.UserRepository
}

func NewHandler(userRepo users.UserRepository) *Handler {
	return &Handler{userRepo: userRepo}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(r.Context(), creds.Email)
	if err != nil {
		log.Print("User not found")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password (assuming you have a VerifyPassword method in your User type)
	if !user.VerifyPassword(creds.Password) {
		log.Print("Invalid password")
		log.Print("User hash:", user.PasswordHash)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := GenerateToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Create response
	uid, err := uuid.Parse(user.ID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}
	response := AuthResponse{
		SessionToken: token,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
		UserID:       uid,
	}

	utils.JSONResponse(w, response, http.StatusOK)
}
