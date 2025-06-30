package notes

type PostgresNoteRepository struct{}

func (r *PostgresNoteRepository) GetAll() ([]*Note, error) {
	return nil, nil
}
func (r *PostgresNoteRepository) Create(note *Note) error {
	return nil
}

func (r *PostgresNoteRepository) GetByID(id int) (*Note, error) {
	return nil, nil
}

func (r *PostgresNoteRepository) Update(note *Note) error {
	return nil
}

func (r *PostgresNoteRepository) Delete(id int) error {
	return nil
}

func (r *PostgresNoteRepository) List() ([]*Note, error) {
	return nil, nil
}
