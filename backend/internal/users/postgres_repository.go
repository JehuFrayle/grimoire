package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Implementation
type PgUserRepository struct {
	DB *pgxpool.Pool
}

func NewPgUserRepository(db *pgxpool.Pool) *PgUserRepository {
	return &PgUserRepository{DB: db}
}

func (r *PgUserRepository) GetAll(ctx context.Context) ([]User, error) {
	query := "SELECT id, username, email, created_at, updated_at, role, active FROM users ORDER BY created_at DESC"
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role, &user.Active); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

func (r *PgUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	row := r.DB.QueryRow(ctx, `SELECT * FROM users WHERE id = $1`, id)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	return &user, err
}

func (r *PgUserRepository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, role, active) 
              VALUES ($1, $2, $3, $4)
              RETURNING id, created_at, updated_at`

	err := r.DB.QueryRow(ctx, query, user.Username, user.Email, user.Role, user.Active).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PostgreSQL error: %s (Code: %s)\n", pgErr.Message, pgErr.Code)
			// Handle specific error codes
			switch pgErr.Code {
			case "23505": // unique_violation
				return fmt.Errorf("user already exists: %w", err)
			case "23503": // foreign_key_violation
				return fmt.Errorf("referenced record does not exist: %w", err)
			default:
				return fmt.Errorf("database error: %w", err)
			}
		}
		return fmt.Errorf("unexpected error: %w", err)
	}
	return nil
}

func (r *PgUserRepository) Update(ctx context.Context, user *User) error {
	return fmt.Errorf("not implemented yet") // Placeholder for user update logic
}

func (r *PgUserRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
