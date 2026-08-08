// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package services

import (
	"github.com/khgreav/chronosplit/models"
	"github.com/khgreav/chronosplit/repos"
)

type CheckpointService struct {
	repo *repos.CheckpointRepo
}

func NewCheckpointService(repo *repos.CheckpointRepo) *CheckpointService {
	return &CheckpointService{repo: repo}
}

func (s *CheckpointService) ListCheckpoints() ([]models.Checkpoint, error) {
	return s.repo.List()
}

func (s *CheckpointService) ListCheckpointsByBlock(blockID int64) ([]models.Checkpoint, error) {
	return s.repo.ListByBlock(blockID)
}

func (s *CheckpointService) ListCheckpointsByProject(projectID int64) ([]models.Checkpoint, error) {
	return s.repo.ListByProject(projectID)
}

func (s *CheckpointService) ListCheckpointsBySubject(subjectID int64) ([]models.Checkpoint, error) {
	return s.repo.ListBySubject(subjectID)
}

func (s *CheckpointService) GetCheckpoint(id int64) (*models.Checkpoint, error) {
	return s.repo.Get(id)
}

func (s *CheckpointService) CreateCheckpoint(c *models.Checkpoint) (*models.Checkpoint, error) {
	return s.repo.Create(c)
}

func (s *CheckpointService) UpdateCheckpoint(c *models.Checkpoint) error {
	return s.repo.Update(c)
}

func (s *CheckpointService) DeleteCheckpoint(id int64) error {
	return s.repo.Delete(id)
}

func (s *CheckpointService) HasCheckpoints(blockID int64) (bool, error) {
	return s.repo.HasCheckpoints(blockID)
}

func (s *CheckpointService) GetLastCheckpoint(blockID int64) (*models.Checkpoint, error) {
	return s.repo.GetLastCheckpoint(blockID)
}
