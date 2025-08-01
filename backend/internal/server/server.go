package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jehufrayle/grimoire/internal/auth"
	"github.com/jehufrayle/grimoire/internal/database"
	"github.com/jehufrayle/grimoire/internal/notes"
	"github.com/jehufrayle/grimoire/internal/users"
	"github.com/jehufrayle/grimoire/middleware"
	"github.com/rs/cors"
)

func StartServer(ctx context.Context, addr string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	mux.HandleFunc("GET /hello", helloHandler)

	// User-related endpoints
	userRepo := users.NewPgUserRepository(database.DB) // change this to change between memory and database
	userHandler := users.NewHandler(userRepo)
	mux.HandleFunc("GET /api/users", userHandler.GetAllUsers)
	mux.HandleFunc("POST /api/users", userHandler.CreateUser)
	mux.HandleFunc("GET /api/users/{id}", userHandler.GetUser)
	mux.HandleFunc("PATCH /api/users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", userHandler.DeleteUser)

	// Login related endpoints
	authHandler := auth.NewHandler(userRepo)
	mux.HandleFunc("POST /api/auth/login", authHandler.LoginHandler)
	mux.HandleFunc("POST /api/auth/signup", authHandler.SignupHandler)
	mux.HandleFunc("GET /api/token", tokenValidatorHandler)

	// Notes related endpoints
	noteRepo := notes.NewPgNoteRepository(database.DB) // change this to change between memory and()
	noteHandler := notes.NewHandler(noteRepo)
	mux.HandleFunc("GET /api/notes", noteHandler.GetUserNotes)
	mux.HandleFunc("GET /api/admin/notes", noteHandler.GetAllNotes)
	mux.HandleFunc("GET /api/notes/{id}", noteHandler.GetUserNoteByID)
	mux.HandleFunc("GET /api/admin/notes/{id}", noteHandler.GetNoteByID)
	mux.HandleFunc("POST /api/notes", noteHandler.CreateNote)
	mux.HandleFunc("PATCH /api/notes/{id}", noteHandler.UpdateNote)
	mux.HandleFunc("DELETE /api/notes/{id}", noteHandler.DeleteNote)

	// Create the HTTP server
	middlewares := middleware.CreateStack(middleware.Logging, middleware.Authentication, middleware.Authorization)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(middlewares(mux))

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Channel to stop the server when necessary
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
		log.Fatalf("graceful shutdown failed: %v", err)
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

func tokenValidatorHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}
	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}
	token := authHeader[len(prefix):]
	claims, err := auth.ValidateToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(claims); err != nil {
		http.Error(w, "Failed to encode claims", http.StatusInternalServerError)
		return
	}
}
