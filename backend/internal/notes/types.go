package notes

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []string  `json:"tags"`
}

type Tag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
