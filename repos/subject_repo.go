package repos

import (
	"github.com/khgreav/chronosplit/models"
	"database/sql"
	"fmt"
)

type SubjectRepo struct {
	db *sql.DB
}

func NewSubjectRepo(db *sql.DB) *SubjectRepo {
	return &SubjectRepo{db: db}
}

func (r *SubjectRepo) List() ([]models.Subject, error) {
	rows, err := r.db.Query(
		`
		SELECT id, name
		FROM subjects
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to query subjects: %w", err)
	}
	defer rows.Close()

	var subjects []models.Subject

	for rows.Next() {
		var s models.Subject
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, fmt.Errorf("Failed to scan subject row: %w", err)
		}
		subjects = append(subjects, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error occurred while iterating rows: %w", err)
	}

	return subjects, nil
}

func (r *SubjectRepo) Get(id int64) (*models.Subject, error) {
	s := &models.Subject{}
	err := r.db.QueryRow(
		`
		SELECT id, name
		FROM subjects
		WHERE id = ?
		`,
		id,
	).Scan(
		&s.ID,
		&s.Name,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SubjectRepo) Create(name string) (*models.Subject, error) {
	result, err := r.db.Exec(
		"INSERT INTO subjects (name) VALUES (?)",
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to insert create new subject: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve subject last insert ID: %w", err)
	}
	return r.Get(id)
}

func (r *SubjectRepo) Delete(id int64) error {
	result, err := r.db.Exec(
		`
		DELETE FROM subjects
		WHERE id = ?
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("Failed to delete subject ID %d: %w", id, err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("No such subject exists.")
	}
	return nil
}
