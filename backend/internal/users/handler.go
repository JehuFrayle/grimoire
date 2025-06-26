// Example for notes/handler.go
package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jehufrayle/grimoire/utils"
)

type Handler struct {
	repo UserRepository
}

func NewHandler(repo UserRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) { // Initialize an empty slice of User
	repo := h.repo // Assuming database.DB is your DB connection

	users, err := repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Convert users to JSON and write to response}
	utils.JSONResponse(w, users, http.StatusOK)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get a specific user by ID
	repo := h.repo
	id := r.PathValue("id")
	log.Print(id)

	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	user, err := repo.GetByID(r.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if user == nil {
		log.Print(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Convert user to JSON and write to response
	utils.JSONResponse(w, user, http.StatusOK)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	repo := h.repo

	// Define a struct to capture both user fields and password from the request body
	var req struct {
		User
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Pass the password separately to the repo.Create method
	if err := repo.Create(r.Context(), &req.User, req.Password); err != nil {
		fmt.Println("Error creating user:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, req.User, http.StatusCreated) // Respond with the created user
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Update user information
	repo := h.repo
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if user.ID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	if err := repo.Update(r.Context(), &user); err != nil {
		fmt.Println("Error updating user:", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, user, http.StatusOK) // Respond with the updated user
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	repo := h.repo
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	err := repo.Delete(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Soft delete successful
	w.WriteHeader(http.StatusNoContent)
}
