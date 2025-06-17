package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jehufrayle/grimoire/internal/database"
	"github.com/jehufrayle/grimoire/internal/notes"
	"github.com/jehufrayle/grimoire/internal/users"
)

func StartServer(ctx context.Context, addr string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	mux.HandleFunc("/", http.NotFound)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/api/notes", notes.NotesHandler)

	// User-related endpoint
	userRepo := users.NewPgUserRepository(database.DB) // change this to change between memory and database
	userHandler := users.NewHandler(userRepo)
	mux.HandleFunc("/api/users/", userHandler.UsersHandler)

	// Create the HTTP server
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	serverStopped := make(chan struct{})

	go func() {
		log.Printf("🧙 Server running on port %s", addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
		close(serverStopped)
	}()

	// Espera a que se cancele el contexto (Ctrl+C, SIGTERM)
	<-ctx.Done()

	log.Println("🛑 Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %w", err)
	}

	<-serverStopped
	log.Println("✅ Server stopped cleanly")
}

// helloHandler is a simple handler that responds with a welcome message
// just for demonstration purposes.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to Grimoire API"))
}
