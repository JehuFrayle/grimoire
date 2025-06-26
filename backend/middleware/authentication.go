package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jehufrayle/grimoire/internal/auth"
)

type contextKey string

const (
	// clave para almacenar datos del usuario en el contexto
	UserIDKey   contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)

// This map sets the public endpoints
var publicPaths = map[string]bool{
	"/api/auth/login":   true,
	"/api/auth/signup":  true,
	"/api/notes/public": true,
	"/hello":            true,
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth if the route is public
		if publicPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// Leer el header Authorization: Bearer <token>
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validar el token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Agregar claims al contexto para que estén disponibles en el handler
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
