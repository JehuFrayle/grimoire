package auth

import (
	"testing"

	"github.com/jehufrayle/grimoire/internal/users"
)

func TestGenerarYValidarToken(t *testing.T) {
	userID := "user123"
	role := users.RoleAdmin

	token, err := GenerateToken(userID, role)
	if err != nil {
		t.Fatalf("Error generando token: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Error validando token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Esperaba userID %s, pero recibí %s", userID, claims.UserID)
	}

	if claims.Role != role {
		t.Errorf("Esperaba role %s, pero recibí %s", role, claims.Role)
	}
}
