package server

import (
	"log"
	"net/http"

	"github.com/jehufrayle/grimoire/internal/notes"
	"github.com/jehufrayle/grimoire/internal/users"
)

func StartServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/api/notes", notes.NotesHandler)
	mux.HandleFunc("/api/users", users.UsersHandler)

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to Grimoire API"))
}
