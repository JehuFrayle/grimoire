package notes

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store interface for notes repository
type Store interface {
	Create(ctx context.Context, note Note) error
	GetByID(ctx context.Context, id string) (*Note, error)
	GetAll(ctx context.Context) ([]Note, error)
	Delete(ctx context.Context, id string) error
}

// DBStore implements Store using a PostgreSQL database
type DBStore struct {
	DB *pgxpool.Pool
}

func NewDBStore(db *pgxpool.Pool) Store {
	return &DBStore{DB: db}
}

func (s *DBStore) Create(ctx context.Context, note Note) error {
	// TODO: Implement SQL INSERT for notes
	// Example: INSERT INTO notes (id, ...) VALUES ($1, ...)
	return fmt.Errorf("not implemented")
}

func (s *DBStore) GetByID(ctx context.Context, id string) (*Note, error) {
	// TODO: Implement SQL SELECT by id
	// Example: SELECT * FROM notes WHERE id = $1
	return nil, fmt.Errorf("not implemented")
}

func (s *DBStore) GetAll(ctx context.Context) ([]Note, error) {
	// TODO: Implement SQL SELECT all notes
	// Example: SELECT * FROM notes
	return nil, fmt.Errorf("not implemented")
}

func (s *DBStore) Delete(ctx context.Context, id string) error {
	// TODO: Implement SQL DELETE by id
	// Example: DELETE FROM notes WHERE id = $1
	return fmt.Errorf("not implemented")
}
