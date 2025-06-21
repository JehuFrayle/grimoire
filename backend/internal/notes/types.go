package notes

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserID    uuid.UUID  `json:"user_id"`
	IsPublic  bool       `json:"is_public"`
	Tags      []Tag      `json:"tags"`
}

type Tag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
