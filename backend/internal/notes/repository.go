package notes

import (
	"context"

	"github.com/google/uuid"
)

// Interface
type NoteRepository interface {
	GetAll(ctx context.Context) ([]Note, error) // Get all users
	GetByID(ctx context.Context, id string) (*Note, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Note, error)
	GetByTags(ctx context.Context, tags []string) ([]Note, error)
	Create(ctx context.Context, note *Note) error
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, id string) error
}
