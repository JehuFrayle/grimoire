package users

import (
	"context"
)

// Interface
type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error) // Get all users
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User, password string) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
