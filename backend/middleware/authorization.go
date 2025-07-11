package middleware

import (
	"net/http"
	"strings"

	"github.com/jehufrayle/grimoire/internal/shared"
)

// Authorization is a middleware that checks for role-based access.
// It protects routes prefixed with "/api/admin" by requiring an admin role.
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the route is an admin route
		if strings.HasPrefix(r.URL.Path, "/api/admin") {
			// Get user role from context (set by Authentication middleware)
			role, ok := r.Context().Value(shared.UserRoleKey).(shared.Role)
			if !ok {
				// This should not happen if Authentication middleware is properly configured
				http.Error(w, "User role not found in context", http.StatusInternalServerError)
				return
			}

			// Check if the user has the admin role
			if role != shared.RoleAdmin {
				http.Error(w, "Forbidden: You do not have access to this resource", http.StatusForbidden)
				return
			}
		}

		// For non-admin routes or authorized admin routes, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
