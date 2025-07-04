// Example for notes/handler.go
package notes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jehufrayle/grimoire/internal/users"
	"github.com/jehufrayle/grimoire/middleware"
	"github.com/jehufrayle/grimoire/utils"
)

type Handler struct {
	repo NoteRepository
}

func NewHandler(repo NoteRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	// Get all notes
	repo := h.repo
	role := r.Context().Value(middleware.UserRoleKey).(users.Role)
	log.Print(role)

	if role != users.RoleAdmin {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	notes, err := repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve notes", http.StatusInternalServerError)
		return
	}
	if notes == nil {
		notes = []Note{}
	}

	// Convert notes to JSON and write to response
	utils.JSONResponse(w, notes, http.StatusOK)
}

func (h *Handler) GetUserNotes(w http.ResponseWriter, r *http.Request) {
	// Get notes for a specific user
	repo := h.repo
	userIDstr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	notes, err := repo.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve notes for user", http.StatusInternalServerError)
		return
	}
	if notes == nil {
		notes = []Note{}
	}

	utils.JSONResponse(w, notes, http.StatusOK)
}

func (h *Handler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	// Get a specific note by ID
	repo := h.repo
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "Note ID is required", http.StatusBadRequest)
		return
	}
	note, err := repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	if note == nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	// Convert note to JSON and write to response
	utils.JSONResponse(w, note, http.StatusOK)
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	// Create a new note
	repo := h.repo
	type req struct {
		Title    string   `json:"title"`
		Content  string   `json:"content"`
		Tags     []string `json:"tags"`
		IsPublic bool     `json:"is_public"`
	}
	var note Note
	var requestBody req

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDstr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	tags := make([]Tag, len(requestBody.Tags))
	for i, tag := range requestBody.Tags {
		tags[i] = Tag{
			Name: tag,
		}
	}
	note = Note{
		Title:    requestBody.Title,
		Content:  requestBody.Content,
		UserID:   userID, // Assuming user_id is passed in the path
		IsPublic: requestBody.IsPublic,
		Tags:     tags,
	}
	if err := repo.Create(r.Context(), &note); err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	// Return the created note with status 201 Created
	utils.JSONResponse(w, note, http.StatusCreated)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	// Update an existing note
	repo := h.repo
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "Note ID is required", http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid note ID format", http.StatusBadRequest)
		return
	}

	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := repo.Update(r.Context(), &note); err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	note.ID = uid

	utils.JSONResponse(w, note, http.StatusOK)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	// Delete a note by ID
	repo := h.repo
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "Note ID is required", http.StatusBadRequest)
		return
	}

	if err := repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (h *Handler) GetNotesByTags(w http.ResponseWriter, r *http.Request) {
	// Get notes by tags
	repo := h.repo
	tagsQuery := r.URL.Query()
	if tagsQuery == nil {
		http.Error(w, "Tags query parameter is required", http.StatusBadRequest)
		return
	}
	tagsList := tagsQuery["tags"]
	if len(tagsList) == 0 {
		http.Error(w, "At least one tag is required", http.StatusBadRequest)
		return
	}
	var tags []string
	for _, tag := range tagsList {
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	if len(tagsQuery) == 0 {
		http.Error(w, "At least one tag is required", http.StatusBadRequest)
		return
	}

	notes, err := repo.GetByTags(r.Context(), tags)
	if err != nil {
		http.Error(w, "Failed to retrieve notes by tags", http.StatusInternalServerError)
		return
	}
	if notes == nil {
		notes = []Note{}
	}

	utils.JSONResponse(w, notes, http.StatusOK)
}
