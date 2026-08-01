package repos

import (
	"database/sql"
	"fmt"
	"time"
)

type BlockRepo struct {
	db *sql.DB
}

func NewBlockRepo(db *sql.DB) *BlockRepo {
	return &BlockRepo{db: db}
}

func (r *BlockRepo) GetActiveBlockId() (*int64, error) {
	var id int64
	err := r.db.QueryRow(`
		SELECT id
		FROM blocks
		WHERE ended_at IS NULL
		LIMIT 1
	`).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("There are no active work blocks.")
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve active block: %w", err)
	}
	return &id, nil
}

func (r *BlockRepo) StartBlock() (*int64, error) {
	result, err := r.db.Exec("INSERT INTO blocks (created_at) VALUES (?)", time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("Failed to create start work block: %w", err)
	}
	id, err := result.LastInsertId()
	return &id, nil
}

func (r *BlockRepo) StopBlock(id int64, timestamp time.Time) error {
	result, err := r.db.Exec(
		`
			UPDATE blocks
			SET ended_at = ?
			WHERE id = ?
		`,
		timestamp.UTC().Format(""),
		id,
	)
	if err != nil {
		return fmt.Errorf("Failed to stop work block: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("No active blocks detected.")
	}
	return nil
}
