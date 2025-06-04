// Example for notes/handler.go
package notes

import "net/http"

func NotesHandler(w http.ResponseWriter, r *http.Request) {
	reqCode := r.Method
	switch reqCode {
	case http.MethodGet:
		// Handle GET request for notes
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("List of notes"))
	case http.MethodPost:
		// Handle POST request to create a new note
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Note created"))
	case http.MethodPut:
		// Handle PUT request to update an existing note
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Note updated"))
	case http.MethodDelete:
		// Handle DELETE request to delete a note
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Note deleted"))
	default:
		// Handle unsupported methods
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}
