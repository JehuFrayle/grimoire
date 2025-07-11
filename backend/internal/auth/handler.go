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

	// Verify password
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

func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newUser := users.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     users.RoleUser, // por defecto
		Active:   true,
		Profile: users.Profile{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
	}

	// Insertar en base de datos y dejar que esta genere ID y timestamps
	if err := h.userRepo.Create(r.Context(), &newUser, req.Password); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Luego del insert, tu repo debe devolver el `ID` generado al struct `newUser.ID`
	token, err := GenerateToken(newUser.ID, newUser.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	resp := AuthResponse{
		SessionToken: token,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
		UserID:       uuid.MustParse(newUser.ID),
	}

	utils.JSONResponse(w, resp, http.StatusCreated)
}
