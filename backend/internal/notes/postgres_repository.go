package notes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const selectNoteWithTagsQuery = `
SELECT
    n.id,
    n.title,
    n.content,
    n.created_at,
    n.updated_at,
    n.user_id,
    n.is_public,
    COALESCE(jsonb_agg(jsonb_build_object('id', t.id, 'name', t.name)) FILTER (WHERE t.id IS NOT NULL), '[]') AS tags
FROM
    active_notes n
LEFT JOIN
    note_tags nt ON n.id = nt.note_id
LEFT JOIN
    tags t ON nt.tag_id = t.id
`

const groupByClause = " GROUP BY n.id, n.title, n.content, n.created_at, n.updated_at, n.user_id, n.is_public"

// PgNoteRepository implements the Repository interface for PostgreSQL.
type PgNoteRepository struct {
	DB *pgxpool.Pool
}

// NewPgNoteRepository creates a new instance of PgNoteRepository.
func NewPgNoteRepository(db *pgxpool.Pool) *PgNoteRepository {
	return &PgNoteRepository{DB: db}
}

// Create handles the creation of a new note and its associated tags.
func (r *PgNoteRepository) Create(ctx context.Context, note *Note) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Rollback is a no-op if tx has been committed.

	// Insert the note
	noteQuery := `
        INSERT INTO notes (user_id, title, content, is_public)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`
	err = tx.QueryRow(ctx, noteQuery, note.UserID, note.Title, note.Content, note.IsPublic).
		Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { // foreign key violation
			return fmt.Errorf("user not found: %w", err)
		}
		return fmt.Errorf("failed to insert note: %w", err)
	}

	// Handle tags
	if len(note.Tags) > 0 {
		for i, tag := range note.Tags {
			// Upsert tag and get its ID
			var tagID uuid.UUID
			tagName := strings.ToLower(tag.Name)
			tagQuery := `
                INSERT INTO tags (name) VALUES ($1)
                ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
                RETURNING id`
			err := tx.QueryRow(ctx, tagQuery, tagName).Scan(&tagID)
			if err != nil {
				return fmt.Errorf("failed to upsert tag '%s': %w", tagName, err)
			}
			note.Tags[i].ID = tagID
			note.Tags[i].Name = tagName // ensure it's lowercase

			// Link tag to note
			noteTagQuery := `
				INSERT INTO note_tags (note_id, tag_id) VALUES ($1, $2)
				ON CONFLICT DO NOTHING`
			_, err = tx.Exec(ctx, noteTagQuery, note.ID, tagID)
			if err != nil {
				return fmt.Errorf("failed to link tag to note: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}

// GetAll retrieves all non-deleted notes, along with their tags.
func (r *PgNoteRepository) GetAll(ctx context.Context) ([]Note, error) {
	query := selectNoteWithTagsQuery + groupByClause + " ORDER BY n.created_at DESC"
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var tagsJSON []byte
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.UserID, &note.IsPublic, &tagsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan note row: %w", err)
		}
		if err := json.Unmarshal(tagsJSON, &note.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags for note %s: %w", note.ID, err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return notes, nil
}

// GetByUserID retrieves all notes for a specific user.
func (r *PgNoteRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]Note, error) {
	query := selectNoteWithTagsQuery + " WHERE n.user_id = $1" + groupByClause + " ORDER BY n.created_at DESC"
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes by user ID: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var tagsJSON []byte
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.UserID, &note.IsPublic, &tagsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan note row: %w", err)
		}
		if err := json.Unmarshal(tagsJSON, &note.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags for note %s: %w", note.ID, err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return notes, nil
}

// GetByID retrieves a single note by its ID.
func (r *PgNoteRepository) GetByID(ctx context.Context, id string) (*Note, error) {
	noteID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid note ID format: %w", err)
	}

	query := selectNoteWithTagsQuery + " WHERE n.id = $1" + groupByClause
	row := r.DB.QueryRow(ctx, query, noteID)

	var note Note
	var tagsJSON []byte
	err = row.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.UserID, &note.IsPublic, &tagsJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("failed to scan note: %w", err)
	}

	if err := json.Unmarshal(tagsJSON, &note.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags for note %s: %w", note.ID, err)
	}

	return &note, nil
}

// GetByTags retrieves all notes that have at least one of the specified tags.
func (r *PgNoteRepository) GetByTags(ctx context.Context, tags []string) ([]Note, error) {
	lowerTags := make([]string, len(tags))
	for i, t := range tags {
		lowerTags[i] = strings.ToLower(t)
	}

	query := selectNoteWithTagsQuery + `
        WHERE EXISTS (
            SELECT 1
            FROM note_tags nt_sub
            JOIN tags t_sub ON nt_sub.tag_id = t_sub.id
            WHERE nt_sub.note_id = n.id AND t_sub.name = ANY($1)
        )` + groupByClause + `
        ORDER BY n.created_at DESC`

	rows, err := r.DB.Query(ctx, query, lowerTags)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes by tags: %w", err)
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		var tagsJSON []byte
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.UserID, &note.IsPublic, &tagsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan note row: %w", err)
		}
		if err := json.Unmarshal(tagsJSON, &note.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags for note %s: %w", note.ID, err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return notes, nil
}

// Update handles the modification of a note's details and its tags.
func (r *PgNoteRepository) Update(ctx context.Context, note *Note) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	updateQuery := `
        UPDATE notes
        SET title = $1, content = $2, is_public = $3, updated_at = now()
        WHERE id = $4 AND deleted_at IS NULL
        RETURNING updated_at`
	err = tx.QueryRow(ctx, updateQuery, note.Title, note.Content, note.IsPublic, note.ID).Scan(&note.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("note not found or already deleted")
		}
		return fmt.Errorf("failed to update note: %w", err)
	}

	// Clear existing tags for the note
	_, err = tx.Exec(ctx, "DELETE FROM note_tags WHERE note_id = $1", note.ID)
	if err != nil {
		return fmt.Errorf("failed to clear existing tags: %w", err)
	}

	// Add new tags
	if len(note.Tags) > 0 {
		for i, tag := range note.Tags {
			var tagID uuid.UUID
			tagName := strings.ToLower(tag.Name)
			tagQuery := `
                INSERT INTO tags (name) VALUES ($1)
                ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
                RETURNING id`
			err := tx.QueryRow(ctx, tagQuery, tagName).Scan(&tagID)
			if err != nil {
				return fmt.Errorf("failed to upsert tag '%s': %w", tagName, err)
			}
			note.Tags[i].ID = tagID
			note.Tags[i].Name = tagName

			noteTagQuery := `
				INSERT INTO note_tags (note_id, tag_id) VALUES ($1, $2)
				ON CONFLICT DO NOTHING`
			_, err = tx.Exec(ctx, noteTagQuery, note.ID, tagID)
			if err != nil {
				return fmt.Errorf("failed to link tag to note: %w", err)
			}
		}
	}

	return tx.Commit(ctx)
}

// Delete performs a soft delete on a note.
func (r *PgNoteRepository) Delete(ctx context.Context, id string) error {
	noteID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid note ID format: %w", err)
	}

	query := `
        UPDATE notes
        SET deleted_at = now(), updated_at = now()
        WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.DB.Exec(ctx, query, noteID)
	if err != nil {
		return fmt.Errorf("failed to soft delete note: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("note not found or already deleted")
	}

	return nil
}
