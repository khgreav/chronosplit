package repos

import (
	"github.com/khgreav/chronosplit/models"
	"database/sql"
	"fmt"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) List() ([]models.Project, error) {
	rows, err := r.db.Query(
		`
		SELECT id, name
		FROM projects
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, fmt.Errorf("Failed to scan project row: %w", err)
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error occurred while iterating rows: %w", err)
	}

	return projects, nil
}

func (r *ProjectRepo) Get(id int64) (*models.Project, error) {
	p := &models.Project{}
	err := r.db.QueryRow(
		`
		SELECT id, name
		FROM projects
		WHERE id = ?
		`,
		id,
	).Scan(
		&p.ID,
		&p.Name,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProjectRepo) Create(name string) (*models.Project, error) {
	result, err := r.db.Exec(
		"INSERT INTO projects (name) VALUES (?)",
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to insert create new project: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve project last insert ID: %w", err)
	}
	return r.Get(id)
}

func (r *ProjectRepo) Delete(id int64) error {
	result, err := r.db.Exec(
		`
		DELETE FROM projects
		WHERE id = ?
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("Failed to delete project ID %d: %w", id, err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("No such project exists.")
	}
	return nil
}
