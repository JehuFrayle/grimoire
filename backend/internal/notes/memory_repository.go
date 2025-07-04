package notes

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type InMemoryNoteRepository struct {
	mu    sync.RWMutex
	notes map[string]*Note
}

func NewInMemoryNoteRepository() *InMemoryNoteRepository {
	return &InMemoryNoteRepository{
		notes: make(map[string]*Note),
	}
}

func (r *InMemoryNoteRepository) Create(ctx context.Context, note *Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	note.ID = uuid.New()
	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now
	note.DeletedAt = nil
	for i, tag := range note.Tags {
		note.Tags[i] = Tag{
			ID:   uuid.New(),
			Name: strings.ToLower(tag.Name), // Normalize tag names to lowercase
		}
	}

	r.notes[note.ID.String()] = note
	return nil
}

func (r *InMemoryNoteRepository) GetAll(ctx context.Context) ([]Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Note
	for _, n := range r.notes {
		if n.DeletedAt == nil {
			result = append(result, *n)
		}
	}
	return result, nil
}

func (r *InMemoryNoteRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []Note

	for _, n := range r.notes {
		if n.DeletedAt == nil && n.UserID == userID {
			result = append(result, *n)
		}
	}
	return result, nil
}

func (r *InMemoryNoteRepository) GetByID(ctx context.Context, id string) (*Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	note, ok := r.notes[id]
	if !ok || note.DeletedAt != nil {
		return nil, errors.New("note not found")
	}
	return note, nil
}

func (r *InMemoryNoteRepository) GetByTags(ctx context.Context, tags []string) ([]Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Note
	tagSet := make(map[string]struct{})
	for _, t := range tags {
		tagSet[strings.ToLower(t)] = struct{}{}
	}

	for _, n := range r.notes {
		if n.DeletedAt != nil {
			continue
		}
		for _, tag := range n.Tags {
			if _, ok := tagSet[strings.ToLower(tag.Name)]; ok {
				result = append(result, *n)
				break
			}
		}
	}
	return result, nil
}

func (r *InMemoryNoteRepository) Update(ctx context.Context, note *Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.notes[note.ID.String()]
	if !ok || existing.DeletedAt != nil {
		return errors.New("note not found")
	}

	note.UpdatedAt = time.Now()
	r.notes[note.ID.String()] = note
	return nil
}

func (r *InMemoryNoteRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	note, ok := r.notes[id]
	if !ok || note.DeletedAt != nil {
		return errors.New("note not found or already deleted")
	}

	now := time.Now()
	note.DeletedAt = &now
	note.UpdatedAt = now
	return nil
}
