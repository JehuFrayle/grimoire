// Example for notes/handler.go
package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jehufrayle/grimoire/utils"
)

type Handler struct {
	repo UserRepository
}

func NewHandler(repo UserRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	// Handle single user requests here
	id := utils.ExtractID(r.URL.Path, "/api/users/")

	if id != "" {
		switch r.Method {
		case http.MethodGet:
			h.GetUser(w, r) // Get a specific user
		case http.MethodPut:
			h.UpdateUser(w, r) // Update user information
		case http.MethodDelete:
			h.DeleteUser(w, r) // Delete a user
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetAllUsers(w, r) // Get all users
	case http.MethodPost:
		h.CreateUser(w, r) // Create a new user
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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
	id := utils.ExtractID(r.URL.Path, "/api/users/")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	user, err := repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Convert user to JSON and write to response
	utils.JSONResponse(w, user, http.StatusOK)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Create a new user
	repo := h.repo
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := repo.Create(r.Context(), &user); err != nil {
		fmt.Println("Error creating user:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, user, http.StatusCreated) // Respond with the created user
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
	// Delete a user by ID
	repo := h.repo
	id := utils.ExtractID(r.URL.Path, "/api/users/")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	if err := repo.Delete(r.Context(), id); err != nil {
		fmt.Println("Error deleting user:", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // Respond with no content on successful deletion
}
