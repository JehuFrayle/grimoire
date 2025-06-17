// Example for notes/handler.go
package notes

import (
	"context"
	"encoding/json"
	"net/http"
)

var store Store // This should be set from main or server package

func SetStore(s Store) {
	store = s
}

func NotesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODO: Parse query params if needed
		notes, err := store.GetAll(context.Background())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get notes"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	case http.MethodPost:
		// TODO: Parse request body for new note
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Create note not implemented"))
	case http.MethodPut:
		// TODO: Parse request body for update
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Update note not implemented"))
	case http.MethodDelete:
		// TODO: Parse request for note ID to delete
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Delete note not implemented"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}
