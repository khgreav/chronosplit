// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package repos

import (
	"database/sql"
	"fmt"

	"github.com/khgreav/chronosplit/models"
)

type CheckpointRepo struct {
	db *sql.DB
}

func NewCheckpointRepo(db *sql.DB) *CheckpointRepo {
	return &CheckpointRepo{db: db}
}

func (r *CheckpointRepo) List() ([]models.Checkpoint, error) {
	rows, err := r.db.Query(
		`
		SELECT id, block_id, project_id, subject_id, previous_checkpoint_id, start_time, end_time, description
		FROM checkpoints
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to query checkpoints: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var checkpoints []models.Checkpoint

	for rows.Next() {
		var c models.Checkpoint
		if err := rows.Scan(
			&c.ID,
			&c.BlockID,
			&c.ProjectID,
			&c.SubjectID,
			&c.PreviousCheckpointID,
			&c.StartTime,
			&c.EndTime,
			&c.Description,
		); err != nil {
			return nil, fmt.Errorf("Failed to scan checkpoint row: %w", err)
		}
		checkpoints = append(checkpoints, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error occurred while iterating rows: %w", err)
	}

	return checkpoints, nil
}

func (r *CheckpointRepo) ListByBlock(blockID int64) ([]models.Checkpoint, error) {
	rows, err := r.db.Query(
		`
		SELECT id, block_id, project_id, subject_id, previous_checkpoint_id, start_time, end_time, description
		FROM checkpoints
		WHERE block_id = ?
		`,
		blockID,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to query checkpoints for block %d: %w", blockID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var checkpoints []models.Checkpoint

	for rows.Next() {
		var c models.Checkpoint
		if err := rows.Scan(
			&c.ID,
			&c.BlockID,
			&c.ProjectID,
			&c.SubjectID,
			&c.PreviousCheckpointID,
			&c.StartTime,
			&c.EndTime,
			&c.Description,
		); err != nil {
			return nil, fmt.Errorf("Failed to scan checkpoint row: %w", err)
		}
		checkpoints = append(checkpoints, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error occurred while iterating rows: %w", err)
	}

	return checkpoints, nil
}

func (r *CheckpointRepo) ListByProject(projectID int64) ([]models.Checkpoint, error) {
	rows, err := r.db.Query(
		`
		SELECT id, block_id, project_id, subject_id, previous_checkpoint_id, start_time, end_time, description
		FROM checkpoints
		WHERE project_id = ?
		`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to query checkpoints for project %d: %w", projectID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var checkpoints []models.Checkpoint

	for rows.Next() {
		var c models.Checkpoint
		if err := rows.Scan(
			&c.ID,
			&c.BlockID,
			&c.ProjectID,
			&c.SubjectID,
			&c.PreviousCheckpointID,
			&c.StartTime,
			&c.EndTime,
			&c.Description,
		); err != nil {
			return nil, fmt.Errorf("Failed to scan checkpoint row: %w", err)
		}
		checkpoints = append(checkpoints, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error occurred while iterating rows: %w", err)
	}

	return checkpoints, nil
}

func (r *CheckpointRepo) ListBySubject(subjectID int64) ([]models.Checkpoint, error) {
	rows, err := r.db.Query(
		`
		SELECT id, block_id, project_id, subject_id, previous_checkpoint_id, start_time, end_time, description
		FROM checkpoints
		WHERE subject_id = ?
		`,
		subjectID,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to query checkpoints for subject %d: %w", subjectID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var checkpoints []models.Checkpoint

	for rows.Next() {
		var c models.Checkpoint
		if err := rows.Scan(
			&c.ID,
			&c.BlockID,
			&c.ProjectID,
			&c.SubjectID,
			&c.PreviousCheckpointID,
			&c.StartTime,
			&c.EndTime,
			&c.Description,
		); err != nil {
			return nil, fmt.Errorf("Failed to scan checkpoint row: %w", err)
		}
		checkpoints = append(checkpoints, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error occurred while iterating rows: %w", err)
	}

	return checkpoints, nil
}

func (r *CheckpointRepo) Get(id int64) (*models.Checkpoint, error) {
	c := &models.Checkpoint{}
	err := r.db.QueryRow(
		`
		SELECT id, block_id, project_id, subject_id, previous_checkpoint_id, start_time, end_time, description
		FROM checkpoints
		WHERE id = ?
		`,
		id,
	).Scan(
		&c.ID,
		&c.BlockID,
		&c.ProjectID,
		&c.SubjectID,
		&c.PreviousCheckpointID,
		&c.StartTime,
		&c.EndTime,
		&c.Description,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CheckpointRepo) Create(c *models.Checkpoint) (*models.Checkpoint, error) {
	result, err := r.db.Exec(
		`
		INSERT INTO checkpoints (block_id, project_id, subject_id, previous_checkpoint_id, start_time, end_time, description)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
		c.BlockID,
		c.ProjectID,
		c.SubjectID,
		c.PreviousCheckpointID,
		c.StartTime,
		c.EndTime,
		c.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to insert new checkpoint: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve checkpoint last insert ID: %w", err)
	}
	c.ID = id
	return r.Get(id)
}

func (r *CheckpointRepo) Update(c *models.Checkpoint) error {
	result, err := r.db.Exec(
		`
		UPDATE checkpoints
		SET block_id = ?, project_id = ?, subject_id = ?, previous_checkpoint_id = ?, start_time = ?, end_time = ?, description = ?
		WHERE id = ?
		`,
		c.BlockID,
		c.ProjectID,
		c.SubjectID,
		c.PreviousCheckpointID,
		c.StartTime,
		c.EndTime,
		c.Description,
		c.ID,
	)
	if err != nil {
		return fmt.Errorf("Failed to update checkpoint: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("No such checkpoint exists.")
	}
	return nil
}

func (r *CheckpointRepo) HasCheckpoints(blockID int64) (bool, error) {
	var exists int
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM checkpoints WHERE block_id = ?)
	`, blockID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("Failed to check if checkpoints exist for block %d: %w", blockID, err)
	}
	return exists == 1, nil
}

func (r *CheckpointRepo) GetLastCheckpoint(blockID int64) (*models.Checkpoint, error) {
	var c models.Checkpoint
	err := r.db.QueryRow(`
		SELECT id, block_id, project_id, subject_id, previous_checkpoint_id, start_time, end_time, description
		FROM checkpoints
		WHERE block_id = ? AND end_time IS NULL
		LIMIT 1
	`, blockID).Scan(
		&c.ID,
		&c.BlockID,
		&c.ProjectID,
		&c.SubjectID,
		&c.PreviousCheckpointID,
		&c.StartTime,
		&c.EndTime,
		&c.Description,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve last checkpoint for block %d: %w", blockID, err)
	}
	return &c, nil
}

func (r *CheckpointRepo) Delete(id int64) error {
	result, err := r.db.Exec(
		`
		DELETE FROM checkpoints
		WHERE id = ?
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("Failed to delete checkpoint ID %d: %w", id, err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("No such checkpoint exists.")
	}
	return nil
}
