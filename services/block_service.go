// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package services

import (
	"time"

	"github.com/khgreav/chronosplit/models"
	"github.com/khgreav/chronosplit/repos"
)

type BlockService struct {
	repo *repos.BlockRepo
}

func NewBlockService(repo *repos.BlockRepo) *BlockService {
	return &BlockService{repo: repo}
}

func (s *BlockService) ActiveBlockExists() (bool, error) {
	return s.repo.ActiveBlockExists()
}

func (s *BlockService) GetActiveBlock() (*models.Block, error) {
	return s.repo.GetActiveBlock()
}

func (s *BlockService) GetActiveBlockID() (*int64, error) {
	return s.repo.GetActiveBlockID()
}

func (s *BlockService) StartBlock() (*int64, error) {
	return s.repo.StartBlock()
}

func (s *BlockService) StopBlock(id int64, timestamp time.Time) error {
	return s.repo.StopBlock(id, timestamp)
}
