package auth

import (
	"time"

	"github.com/google/uuid"
)

// Session represents an authenticated session for a user.
type Session struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// AuthCredentials represents the payload used for login.
type AuthCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is the response sent back after successful login.
type AuthResponse struct {
	SessionToken string    `json:"token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       uuid.UUID `json:"user_id"`
}
